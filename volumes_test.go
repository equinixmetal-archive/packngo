package packngo

import (
	"fmt"
	"testing"
	"time"
)

func waitVolumeActive(id string, c *Client) (*Volume, error) {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		c, _, err := c.Volumes.Get(id, nil)
		if err != nil {
			return nil, err
		}
		if c.State == "active" {
			return c, nil
		}
		<-time.After(5 * time.Second)
	}
	return nil, fmt.Errorf("volume %s is still not active after timeout", id)
}

func TestAccVolumeBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	sp := SnapshotPolicy{
		SnapshotFrequency: "1day",
		SnapshotCount:     3,
	}

	vcr := VolumeCreateRequest{
		Size:             10,
		BillingCycle:     "hourly",
		PlanID:           "storage_1",
		FacilityID:       testFacility(),
		SnapshotPolicies: []*SnapshotPolicy{&sp},
		Description:      "ahoj!",
		Locked:           true,
	}

	v, _, err := c.Volumes.Create(&vcr, projectID)
	if err != nil {
		t.Fatal(err)
	}

	v, err = waitVolumeActive(v.ID, c)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Volumes.Delete(v.ID)

	v, _, err = c.Volumes.Get(v.ID,
		&GetOptions{Includes: []string{"snapshot_policies", "facility"}})
	if err != nil {
		t.Fatal(err)
	}

	if len(v.SnapshotPolicies) != 1 {
		t.Fatal("Test volume should have one snapshot policy")
	}

	if v.SnapshotPolicies[0].SnapshotFrequency != sp.SnapshotFrequency {
		t.Fatal("Test volume has wrong snapshot frequency")
	}

	if v.SnapshotPolicies[0].SnapshotCount != sp.SnapshotCount {
		t.Fatal("Test volume has wrong snapshot count")
	}

	if v.Facility.Code != testFacility() {
		t.Fatal("Test volume has wrong facility")
	}
	_, err = c.Volumes.Unlock(v.ID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestAccVolumeUpdate(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	sp := SnapshotPolicy{
		SnapshotFrequency: "1day",
		SnapshotCount:     3,
	}

	vcr := VolumeCreateRequest{
		Size:             10,
		BillingCycle:     "hourly",
		PlanID:           "storage_1",
		FacilityID:       testFacility(),
		SnapshotPolicies: []*SnapshotPolicy{&sp},
	}

	v, _, err := c.Volumes.Create(&vcr, projectID)
	if err != nil {
		t.Fatal(err)
	}

	v, err = waitVolumeActive(v.ID, c)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Volumes.Delete(v.ID)

	vDesc := "new Desc"

	vur := VolumeUpdateRequest{Description: &vDesc}

	_, _, err = c.Volumes.Update(v.ID, &vur)
	if err != nil {
		t.Fatal(err)
	}

	v, _, err = c.Volumes.Get(v.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if v.Description != vDesc {
		t.Fatalf("Volume desc should be %q, but is %q", vDesc, v.Description)
	}

	newSize := 15

	vur = VolumeUpdateRequest{Size: &newSize}

	_, _, err = c.Volumes.Update(v.ID, &vur)
	if err != nil {
		t.Fatal(err)
	}

	v, _, err = c.Volumes.Get(v.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if v.Size != newSize {
		t.Fatalf("Volume size should be %q, but is %q", newSize, v.Size)
	}

	newPlan := "storage_2"

	vur = VolumeUpdateRequest{PlanID: &newPlan}

	_, _, err = c.Volumes.Update(v.ID, &vur)
	if err != nil {
		t.Fatal(err)
	}

	v, _, err = c.Volumes.Get(v.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if v.Plan.Slug != newPlan {
		t.Fatalf("Plan should be %q, but is %q", newPlan, v.Plan.Slug)
	}

}

func TestAccVolumeLargeList(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	sp := SnapshotPolicy{
		SnapshotFrequency: "1day",
		SnapshotCount:     3,
	}

	vcr := VolumeCreateRequest{
		Size:             10,
		BillingCycle:     "hourly",
		PlanID:           "storage_1",
		FacilityID:       testFacility(),
		SnapshotPolicies: []*SnapshotPolicy{&sp},
	}

	numOfVolumes := 11
	createdVolumes := make([]Volume, numOfVolumes)
	for i := 0; i < numOfVolumes; i++ {
		vcr.Description = randString8()
		v, _, err := c.Volumes.Create(&vcr, projectID)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Volumes.Delete(v.ID)
		createdVolumes[i] = *v
	}

	for _, volume := range createdVolumes {
		if _, err := waitVolumeActive(volume.ID, c); err != nil {
			t.Fatal(err)
		}
	}

	volumes, _, err := c.Volumes.List(projectID, nil)
	if err != nil {
		t.Fatalf("failed to get list of volumes: %v", err)
	}

	if len(volumes) < numOfVolumes {
		t.Fatalf("failed due to expecting at least %d volumes, but actually got %d", numOfVolumes, len(volumes))
	}

	volumeMap := map[string]Volume{}
	for _, volume := range volumes {
		volumeMap[volume.ID] = volume
	}

	for _, k := range createdVolumes {
		if _, ok := volumeMap[k.ID]; !ok {
			t.Fatalf("failed to find expected volume in list: %s", k.ID)
		}
	}

	perPage := 4
	listOpt := &ListOptions{
		Page:    2,
		PerPage: perPage,
	}

	volumes, _, err = c.Volumes.List(projectID, listOpt)
	if err != nil {
		t.Fatalf("failed to get list of volumes: %v", err)
	}

	if len(volumes) != perPage {
		t.Fatalf("failed due to expecting %d volumes, but actually got %d", perPage, len(volumes))
	}

	// Last test get all volume, 5 per page and includes

	listOpt2 := &ListOptions{
		Includes: []string{"snapshot_policies", "facility"},
		PerPage:  5,
	}
	volumes, _, err = c.Volumes.List(projectID, listOpt2)
	if err != nil {
		t.Fatalf("failed to get list of volumes: %v", err)
	}
	if len(volumes) != numOfVolumes {
		t.Fatalf("failed due to expecting at %d volumes, but actually got %d", numOfVolumes, len(volumes))
	}
	for _, v := range volumes {
		if v.SnapshotPolicies[0].SnapshotCount != 3 {
			t.Fatalf("Wrong SpanshotCount, perhaps it's not included?")
		}
	}

}
