package packngo

import (
	"testing"
)

func TestAccConnectionProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	connReq := ConnectionCreateRequest{
		Name:       "testconn",
		Redundancy: ConnectionRedundant,
		Facility:   "ny5",
		Type:       ConnectionShared,
		Project:    projectID,
	}

	conn, _, err := c.Connections.ProjectCreate(projectID, &connReq)
	if err != nil {
		t.Fatal(err)
	}

	createdConnID := conn.ID

	conn, _, err = c.Connections.Get(conn.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if conn.ID != createdConnID {
		t.Fatalf("connection obtained over GET has different ID than created connection (%s vs %s)", conn.ID, createdConnID)
	}

	ports, _, err := c.Connections.Ports(conn.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) == 0 {
		t.Fatal("New connections should have nonzero ports")
	}

	port, _, err := c.Connections.Port(conn.ID, ports[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if port.ID != ports[0].ID {
		t.Fatalf("Mismatch when getting Connection Port, ID should be %s, was %s", ports[0].ID, port.ID)
	}

	_, _, err = c.Connections.PortEvents(conn.ID, port.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	vcs, _, err := c.Connections.VirtualCircuits(conn.ID, port.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(vcs) > 0 {
		vc, _, err := c.VirtualCircuits.Get(vcs[0].ID, nil)
		if err != nil {
			t.Fatal(err)
		}
		_, _, err = c.VirtualCircuits.Events(vc.ID, nil)
		if err != nil {
			t.Fatal(err)
		}
		/*
			        fails with "Virtual Circuits on shared connections may not be deleted."
					_, err = c.Connections.DeleteVirtualCircuit(vc.ID)
					if err != nil {
						t.Fatal(err)
					}
		*/
	}

	conns, _, err := c.Connections.ProjectList(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	found := false

	for _, c := range conns {
		if c.ID == conn.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("The test Project Connection with ID %s was not created", conn.ID)
	}

	connEvents, _, err := c.Connections.Events(conn.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(connEvents) == 0 {
		t.Fatal("There should be some events for the test connection")
	}

	_, err = c.Connections.Delete(conn.ID, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccConnectionOrganization(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	connReq := ConnectionCreateRequest{
		Name:       "testconn",
		Redundancy: ConnectionRedundant,
		Facility:   "ny5",
		Type:       ConnectionShared,
		Project:    projectID,
	}
	user, _, err := c.Users.Current()

	if err != nil {
		t.Fatal(err)
	}

	conn, _, err := c.Connections.OrganizationCreate(user.DefaultOrganizationID, &connReq)
	if err != nil {
		t.Fatal(err)
	}

	conns, _, err := c.Connections.OrganizationList(user.DefaultOrganizationID, nil)
	if err != nil {
		t.Fatal(err)
	}

	updReq := ConnectionUpdateRequest{Redundancy: ConnectionPrimary}
	conn, _, err = c.Connections.Update(conn.ID, &updReq, nil)
	if err != nil {
		t.Fatal(err)
	}

	if conn.Redundancy != ConnectionPrimary {
		t.Fatalf("Updated connection should be primary")
	}

	found := false

	for _, c := range conns {
		if c.ID == conn.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("The test Organization Connection with ID %s was not created", conn.ID)
	}

	_, err = c.Connections.Delete(conn.ID, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccConnectionFabricTokenRedundant(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	cr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "vlan1",
		Metro:       testMetro(),
	}

	vlan1, _, err := c.ProjectVirtualNetworks.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	cr.Description = "vlan2"
	vlan2, _, err := c.ProjectVirtualNetworks.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	connReq := ConnectionCreateRequest{
		Name:             "testconn_redundant",
		Redundancy:       ConnectionRedundant,
		Metro:            testMetro(),
		Type:             ConnectionShared,
		ServiceTokenType: FabricServiceTokenASide,
		VLANs:            []int{vlan1.VXLAN, vlan2.VXLAN},
		ContactEmail:     "nobody@hmail.com",
		Speed:            500000000,
		Project:          projectID,
	}

	co, _, err := c.Connections.ProjectCreate(projectID, &connReq)
	if err != nil {
		t.Fatal(err)
	}

	co, _, err = c.Connections.Get(co.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if co.Ports[0].VirtualCircuits[0].VNID != vlan1.VXLAN {
		t.Fatalf("VNID of first port is not the same as VLAN1 VNID")
	}

	if co.Ports[1].VirtualCircuits[0].VNID != vlan2.VXLAN {
		t.Fatalf("VNID of second port is not the same as VLAN2 VNID")
	}

	maybeVlan1, _, err := c.ProjectVirtualNetworks.GetByVXLAN(projectID, vlan1.VXLAN, nil)
	if err != nil {
		t.Fatal(err)
	}
	if maybeVlan1.ID != vlan1.ID {
		t.Fatalf("VLAN1 VNID does not match VLAN ID fetched by vxlan")
	}

	_, err = c.Connections.Delete(co.ID, true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.ProjectVirtualNetworks.Delete(vlan1.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.ProjectVirtualNetworks.Delete(vlan2.ID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestAccConnectionFabricTokenSingle(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	cr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "vlan1",
		Metro:       testMetro(),
	}

	vlan1, _, err := c.ProjectVirtualNetworks.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	connReq := ConnectionCreateRequest{
		Name:             "testconn_single",
		Redundancy:       ConnectionPrimary,
		Metro:            testMetro(),
		Type:             ConnectionShared,
		ServiceTokenType: FabricServiceTokenASide,
		VLANs:            []int{vlan1.VXLAN},
		ContactEmail:     "nobody@hmail.com",
		Speed:            500000000,
		Project:          projectID,
	}

	co, _, err := c.Connections.ProjectCreate(projectID, &connReq)
	if err != nil {
		t.Fatal(err)
	}

	co, _, err = c.Connections.Get(co.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if co.Ports[0].VirtualCircuits[0].VNID != vlan1.VXLAN {
		t.Fatalf("VNID of first port is not the same as VLAN1 VNID")
	}

	maybeVlan1, _, err := c.ProjectVirtualNetworks.GetByVXLAN(projectID, vlan1.VXLAN, nil)
	if err != nil {
		t.Fatal(err)
	}
	if maybeVlan1.ID != vlan1.ID {
		t.Fatalf("VLAN1 VNID does not match VLAN ID fetched by vxlan")
	}

	_, err = c.Connections.Delete(co.ID, true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.ProjectVirtualNetworks.Delete(vlan1.ID)
	if err != nil {
		t.Fatal(err)
	}

}
