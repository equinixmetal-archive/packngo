package packngo

import (
	"fmt"
	"testing"
	"time"
)

func waitMetalGatewayActive(id string, c *Client) (*MetalGateway, error) {
	includes := &GetOptions{Includes: []string{"ip_reservation", "virtual_network"}}

	for i := 0; i < 12; i++ {
		r, _, err := c.MetalGateways.Get(id, includes)
		if err != nil {
			return nil, err
		}
		if r.State == MetalGatewayActive {
			return r, nil
		}
		<-time.After(5 * time.Second)
	}
	return nil, fmt.Errorf("Metal gateway %s is still not active after timeout", id)
}

func TestAccMetalGatewaySubnetSize(t *testing.T) {

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

	rcr := MetalGatewayCreateRequest{
		VirtualNetworkID:      vlan.ID,
		PrivateIPv4SubnetSize: 8,
	}

	router, _, err := c.MetalGateways.Create(projectID, &rcr)
	if err != nil {
		t.Fatal(err)
	}

	router, err = waitMetalGatewayActive(router.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	routers, _, err := c.MetalGateways.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(routers) != 1 {
		t.Fatalf("There should be exactly one metal gateway in the testing project")
	}

	_, err = c.MetalGateways.Delete(router.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.ProjectVirtualNetworks.Delete(vlan.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccMetalGatewayExistingReservation(t *testing.T) {

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

	rcr := MetalGatewayCreateRequest{
		VirtualNetworkID: vlan.ID,
		IPReservationID:  ipRes.ID,
	}

	router, _, err := c.MetalGateways.Create(projectID, &rcr)
	if err != nil {
		t.Fatal(err)
	}

	router, err = waitMetalGatewayActive(router.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	routers, _, err := c.MetalGateways.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(routers) != 1 {
		t.Fatalf("There should be exactly one metal gateway in the testing project")
	}

	_, err = c.MetalGateways.Delete(router.ID)
	if err != nil {
		t.Fatal(err)
	}

	deleteProjectIP(t, c, ipRes.ID)

	_, err = c.ProjectVirtualNetworks.Delete(vlan.ID)
	if err != nil {
		t.Fatal(err)
	}
}
