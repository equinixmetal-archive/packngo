package packngo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func removeVirtualNetwork(t *testing.T, c *Client, id string) {
	_, err := c.ProjectVirtualNetworks.Delete(id)
	if err != nil {
		t.Log("Err when removing testing VLAN:", err)
	}
}

func removeVirtualCircuit(t *testing.T, c *Client, id string) {
	_, err := c.VirtualCircuits.Delete(id)
	if err != nil {
		t.Log("Err when removing testing VirtualCircuit:", err)
	}
}

func waitVirtualCircuitStatus(t *testing.T, c *Client, id string, status VCStatus, errStati []string) (*VirtualCircuit, error) {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		vc, _, err := c.VirtualCircuits.Get(id, nil)
		if err != nil {
			return nil, err
		}
		if vc.Status == status {
			return vc, nil
		}
		if contains(errStati, string(vc.Status)) {
			return nil, fmt.Errorf("VirtualCircuit %s ended up in status %s", vc.ID, vc.Status)
		}
	}
	return nil, fmt.Errorf("Virtual Circuit %s is still not %s after timeout", id, status)
}

// this test needs an existing Dedicated Connection. Pass the ID in env var
// "PACKNGO_TEST_CONNECTION"
func TestAccVirtualCircuitDedicated(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	testConnectionEnvVar := "PACKNGO_TEST_CONNECTION"

	cid := os.Getenv(testConnectionEnvVar)
	if cid == "" {
		t.Skipf("%s is not set", testConnectionEnvVar)
	}

	conn, _, err := c.Connections.Get(cid, &GetOptions{Includes: []string{"facility"}})
	if err != nil {
		t.Fatal(err)
	}

	fac := conn.Facility.Code

	vncr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "VLAN for VirtualCircuit test",
		Facility:    fac,
	}
	vlan, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer removeVirtualNetwork(t, c, vlan.ID)

	vncr2 := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "VLAN for VirtualCircuit test2",
		Facility:    fac,
	}
	vlan2, _, err := c.ProjectVirtualNetworks.Create(&vncr2)
	if err != nil {
		t.Fatal(err)
	}
	defer removeVirtualNetwork(t, c, vlan2.ID)

	primaryPort := conn.PortByRole(ConnectionPortPrimary)

	cr := VCCreateRequest{
		VirtualNetworkID: vlan.ID,
		NniVLAN:          889,
		Name:             "TestDedicatedConn1",
	}

	vc, _, err := c.VirtualCircuits.Create(projectID, conn.ID, primaryPort.ID, &cr, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer removeVirtualCircuit(t, c, vc.ID)

	_, err = waitVirtualCircuitStatus(t, c, vc.ID, VCStatusActive,
		[]string{string(VCStatusActivationFailed)})
	if err != nil {
		t.Fatal(err)
	}

	cr2 := VCCreateRequest{
		VirtualNetworkID: vlan2.ID,
		NniVLAN:          891,
		Name:             "TestDedicatedConn2",
	}

	vc2, _, err := c.VirtualCircuits.Create(projectID, conn.ID, primaryPort.ID, &cr2, nil)
	if err != nil {
		t.Fatal(err)
	}

	defer removeVirtualCircuit(t, c, vc2.ID)

	_, err = waitVirtualCircuitStatus(t, c, vc2.ID, VCStatusActive,
		[]string{string(VCStatusActivationFailed)})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.VirtualCircuits.Update(vc.ID, &VCUpdateRequest{VirtualNetworkID: nil}, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.VirtualCircuits.Update(vc2.ID, &VCUpdateRequest{VirtualNetworkID: nil}, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = waitVirtualCircuitStatus(t, c, vc.ID, VCStatusWaiting,
		[]string{string(VCStatusDeactivationFailed)})
	if err != nil {
		t.Fatal(err)
	}

	_, err = waitVirtualCircuitStatus(t, c, vc2.ID, VCStatusWaiting,
		[]string{string(VCStatusDeactivationFailed)})
	if err != nil {
		t.Fatal(err)
	}
}

// this test needs an existing Shared Connection. Pass the ID in env var
// "PACKNGO_TEST_CONNECTION"
func TestAccVirtualCircuitShared(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()
	testConnectionEnvVar := "PACKNGO_TEST_CONNECTION"

	cid := os.Getenv(testConnectionEnvVar)
	if cid == "" {
		t.Skipf("%s is not set", testConnectionEnvVar)
	}

	conn, _, err := c.Connections.Get(cid, &GetOptions{Includes: []string{"facility"}})
	if err != nil {
		t.Fatal(err)
	}
	vc, _, err := c.VirtualCircuits.Get(conn.Ports[0].VirtualCircuits[0].ID,
		&GetOptions{Includes: []string{"project"}})
	if err != nil {
		t.Fatal(err)
	}
	if vc.Status == VCStatusActive {
		vc, _, err = c.VirtualCircuits.Update(
			vc.ID,
			&VCUpdateRequest{VirtualNetworkID: nil},
			&GetOptions{Includes: []string{"project"}},
		)
		if err != nil {
			t.Fatal(err)
		}
		_, err = waitVirtualCircuitStatus(t, c, vc.ID, VCStatusWaiting,
			[]string{string(VCStatusDeactivationFailed)})
		if err != nil {
			t.Fatal(err)
		}
	}

	fac := conn.Facility.Code
	projectID := vc.Project.ID

	cr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "VLAN for VirtualCircuit test",
		Facility:    fac,
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer removeVirtualNetwork(t, c, vlan.ID)

	vc, _, err = c.VirtualCircuits.Update(
		vc.ID,
		&VCUpdateRequest{VirtualNetworkID: &(vlan.ID)},
		&GetOptions{Includes: []string{"virtual_network"}},
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = waitVirtualCircuitStatus(t, c, vc.ID, VCStatusActive,
		[]string{string(VCStatusActivationFailed)})
	if err != nil {
		t.Fatal(err)
	}

	vc, _, err = c.VirtualCircuits.Get(vc.ID,
		&GetOptions{Includes: []string{"virtual_network,project"}})
	if err != nil {
		t.Fatal(err)
	}

	if vc.VirtualNetwork.ID != vlan.ID {
		t.Fatalf("ID of assigned vlan from the virtual circuit should be %s, was %s",
			vlan.ID, vc.VirtualNetwork.ID)
	}

	if vc.VNID != vlan.VXLAN {
		t.Fatalf("Numerical ID of assigned vlan from the virtual circuit should be %d, was %d",
			vlan.VXLAN, vc.VNID)
	}
	_, _, err = c.VirtualCircuits.Update(vc.ID, &VCUpdateRequest{VirtualNetworkID: nil}, nil)
	log.Println(vc.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = waitVirtualCircuitStatus(t, c, vc.ID, VCStatusWaiting,
		[]string{string(VCStatusDeactivationFailed)})
	if err != nil {
		t.Fatal(err)
	}
}
