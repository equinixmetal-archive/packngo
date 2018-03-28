package packngo

import (
	"sync"
	"testing"
)

func TestAccProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer projectTeardown(c)

	rs := testProjectPrefix + randString8()
	pcr := ProjectCreateRequest{Name: rs}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
	}
	rs = testProjectPrefix + randString8()
	pur := ProjectUpdateRequest{ID: p.ID, Name: rs}
	p, _, err = c.Projects.Update(&pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}
	gotProject, _, err := c.Projects.Get(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotProject.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotProject.Name)
	}

	if gotProject.PaymentMethod.URL == "" {
		t.Fatalf("Empty payment_method: %v", gotProject)
	}
	_, err = c.Projects.Delete(p.ID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestAccCreateOrgProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer projectTeardown(c)

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	rs := testProjectPrefix + randString8()

	orgPath := "/organizations/" + u.DefaultOrganizationID
	pcr := ProjectCreateRequest{Name: rs}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Organization.URL != orgPath {
		t.Fatalf("Expected new project to be part of org %s, not %v", orgPath, p.Organization)
	}
}

func TestAccListVolumesLargeList(t *testing.T) {
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

	var wg sync.WaitGroup
	numOfVolumes := 11
	volumes := make([]Volume, numOfVolumes)
	for i := 0; i < numOfVolumes; i++ {
		vcr.Description = randString8()
		v, _, err := c.Volumes.Create(&vcr, projectID)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Volumes.Delete(v.ID)
		volumes[i] = *v
	}

	wg.Add(numOfVolumes)
	for _, volume := range volumes {
		go func(volume Volume) {
			defer wg.Done()
			_, err := waitVolumeActive(volume.ID, c)
			if err != nil {
				t.Fatal(err)
			}
		}(volume)
	}
	wg.Wait()

	volumes, _, err := c.Projects.ListVolumes(projectID, nil)
	if err != nil {
		t.Fatalf("failed to get list of sshvolumes: %v", err)
	}

	if len(volumes) < numOfVolumes {
		t.Fatalf("failed due to expecting at least %d volumes, but actually got %d", numOfVolumes, len(volumes))
	}

	volumeMap := map[string]Volume{}
	for _, volume := range volumes {
		volumeMap[volume.ID] = volume
	}

	for _, k := range volumes {
		if _, ok := volumeMap[k.ID]; !ok {
			t.Fatalf("failed to find expected volume in list: %s", k.ID)
		}
	}

	perPage := 4
	listOpt := &ListOptions{
		Page:    2,
		PerPage: perPage,
	}

	volumes, _, err = c.Projects.ListVolumes(projectID, listOpt)
	if err != nil {
		t.Fatalf("failed to get list of sshvolumes: %v", err)
	}

	if len(volumes) != perPage {
		t.Fatalf("failed due to expecting %d volumes, but actually got %d", perPage, len(volumes))
	}
}
