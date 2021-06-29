package packngo

import (
	"fmt"
	"testing"
	"time"
)

func waitSubnetRouterActive(id string, c *Client) (*SubnetRouter, error) {
	// 1 minute = 12 * 5sec-retry
	includes := &GetOptions{Includes: []string{"ip_reservation", "virtual_network"}}

	for i := 0; i < 12; i++ {
		r, _, err := c.SubnetRouters.Get(id, includes)
		if err != nil {
			return nil, err
		}
		if r.State == "active" {
			return r, nil
		}
		<-time.After(5 * time.Second)
	}
	return nil, fmt.Errorf("volume %s is still not active after timeout", id)
}

func TestAccSubnetRouterSubnetSize(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	testDesc := "test_desc_" + randString8()

	vcr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: testDesc,
		Metro:       testMetro(),
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

	router, err = waitSubnetRouterActive(router.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	//log.Println(jstr(router))

	routers, _, err := c.SubnetRouters.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(routers) != 1 {
		t.Fatalf("There should be exactly one subnet router in the testing project")
	}

	_, err = c.SubnetRouters.Delete(router.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.ProjectVirtualNetworks.Delete(vlan.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccSubnetRouterExistingReservation(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	testDesc := "test_desc_" + randString8()

	vcr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: testDesc,
		Metro:       testMetro(),
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&vcr)
	if err != nil {
		t.Fatal(err)
	}
	metro := testMetro()

	ipcr := IPReservationRequest{
		Type:                   PublicIPv4,
		Quantity:               8,
		Metro:                  &metro,
		FailOnApprovalRequired: true,
	}
	ipRes, _, err := c.ProjectIPs.Request(projectID, &ipcr)
	if err != nil {
		t.Fatal(err)
	}

	rcr := SubnetRouterCreateRequest{
		VirtualNetworkID: vlan.ID,
		IPReservationID:  ipRes.ID,
	}

	router, _, err := c.SubnetRouters.Create(projectID, &rcr)
	if err != nil {
		t.Fatal(err)
	}

	router, err = waitSubnetRouterActive(router.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	//log.Println(jstr(router))

	routers, _, err := c.SubnetRouters.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(routers) != 1 {
		t.Fatalf("There should be exactly one subnet router in the testing project")
	}

	/*
		    //check for subnet_router attribute of IP reservation:

			ip, _, err := c.ProjectIPs.Get(ipRes.ID, nil)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(jstr(ip))
	*/

	_, err = c.SubnetRouters.Delete(router.ID)
	if err != nil {
		t.Fatal(err)
	}

	deleteProjectIP(t, c, ipRes.ID)

	_, err = c.ProjectVirtualNetworks.Delete(vlan.ID)
	if err != nil {
		t.Fatal(err)
	}
}
