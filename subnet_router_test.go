package packngo

import (
	"log"
	"testing"
)

func TestAccSubnetRouterSubnetSize(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	testDesc := "test_desc_" + randString8()

	vcr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: testDesc,
		Facility:    testFacility(),
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&vcr)
	if err != nil {
		t.Fatal(err)
	}

	rcr := SubnetRouterCreateRequest{
		VirtualNetworkID:      vlan.ID,
		PrivateIPv4SubnetSize: 8,
	}

	router, _, err := c.SubnetRouters.Create(projectID, &rcr)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(jstr(router))

	routers, _, err := c.SubnetRouters.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(len(routers))

	_, err = c.SubnetRouters.Delete(router.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.ProjectVirtualNetworks.Delete(vlan.ID)
	if err != nil {
		t.Fatal(err)
	}
}
