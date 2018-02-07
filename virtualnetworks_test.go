package packngo

import (
	"testing"
)

func TestAccVirtualNetworks(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	l, _, err := c.ProjectVirtualNetworks.List(projectID)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.VirtualNetworks) != 0 {
		t.Fatal("Newly created project should not have any vlans")

	}
	l, _, err = c.ProjectVirtualNetworks.List(projectID)
	if err != nil {
		t.Fatal(err)
	}

	testDesc := "test_desc_" + randString8()

	cr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: testDesc,
		Facility:    testFacility(),
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	if vlan.Description != testDesc {
		t.Fatal("Wrong description string in created VLAN")
	}

	l, _, err = c.ProjectVirtualNetworks.List(projectID)
	if err != nil {
		t.Fatal(err)
	}

	if len(l.VirtualNetworks) != 1 {
		t.Fatal("At this point, there should be exactly 1 VLAN in the project")
	}

	_, err = c.ProjectVirtualNetworks.Delete(l.VirtualNetworks[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	l, _, err = c.ProjectVirtualNetworks.List(projectID)
	if err != nil {
		t.Fatal(err)
	}
	if len(l.VirtualNetworks) != 0 {
		t.Fatal("The test project should not have any VLANs now")
	}

	// TODO: Test several bad inputs to ensure rejection without adverse affects
	// Create virtual network with bad POST body parameters
	// Ensure create failed
}
