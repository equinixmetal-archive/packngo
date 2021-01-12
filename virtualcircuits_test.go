package packngo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func waitVirtualCircuitStatus(t *testing.T, c *Client, id, status string) *VirtualCircuit {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		vc, _, err := c.VirtualCircuits.Get(id, nil)
		if err != nil {
			t.Fatal(err)
			return nil
		}
		if vc.Status == status {
			return vc
		}
	}
	t.Fatal(fmt.Errorf("Virtual Circuit %s is still not %s after timeout", id, status))
	return nil
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
	if vc.Status == vcStatusActive {
		vc, _, err = c.VirtualCircuits.RemoveVLAN(vc.ID, &GetOptions{Includes: []string{"project"}})
		if err != nil {
			t.Fatal(err)
		}
		waitVirtualCircuitStatus(t, c, vc.ID, vcStatusWaiting)
	}

	fac := conn.Facility.Code
	projectID := vc.Project.ID

	cr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "VLAN for VirtualCircuiti test",
		Facility:    fac,
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err := c.ProjectVirtualNetworks.Delete(vlan.ID)
		if err != nil {
			t.Log("Err when removing testing VLAN:", err)
		}
	}()

	vc, _, err = c.VirtualCircuits.ConnectVLAN(vc.ID, vlan.ID,
		&GetOptions{Includes: []string{"virtual_network"}})
	if err != nil {
		t.Fatal(err)
	}
	waitVirtualCircuitStatus(t, c, vc.ID, vcStatusActive)

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
	_, _, err = c.VirtualCircuits.RemoveVLAN(vc.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

}
