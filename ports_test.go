package packngo

import (
	"path"
	"testing"
)

// run this test as
// PACKNGO_TEST_FACILITY=atl1 PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccPort1E
// .. you can choose another facility, but there must be Type 1E available

func TestAccPort1E(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	// MARK_1

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     testFacility(),
		Plan:         "baremetal_1e",
		OS:           "ubuntu_16_04",
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID)
	dID := d.ID

	// If you need to test this, run a 1e device in your project in a faciltiy
	// and then comment code from MARK_1 to here and uncomment following.
	// Fill the values from your device, project and facility.

	/*
		c := setup(t)

		dID := "b904ec58-4da7-438f-b4d2-e1d0c2f2eeeb"
		projectID := "52000fb2-ee46-4673-93a8-de2c2bdba33b"
		fac := "atl1"
	*/

	d, err = waitDeviceActive(dID, c)
	if err != nil {
		t.Fatal(err)
	}

	eth1, err := c.DevicePorts.GetPortByName(d.ID, "eth1")

	if err != nil {
		t.Fatal(err)
	}

	if len(eth1.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to a eth1 in the begining of this test")
	}

	vncr := VirtualNetworkCreateRequest{
		ProjectID: projectID,
		Facility:  testFacility(),
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.ProjectVirtualNetworks.Delete(vlan.ID)

	p, _, err := c.DevicePorts.Assign(eth1.ID, vlan.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.AttachedVirtualNetworks) != 1 {
		t.Fatal("Exactly one vlan should be attached to a eth1 at this point")
	}

	if path.Base(p.AttachedVirtualNetworks[0].Href) != vlan.ID {
		t.Fatal("mismatch in the UUID of the assigned VLAN")
	}

	p, _, err = c.DevicePorts.Unassign(eth1.ID, vlan.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to the port at this time")
	}

	// TODO: Figure out how to mock test this.
	// attempt to assign virtual network to a non-bonded port                     (assert failure)
	// assign virtual network to bonded port                                      (assert success)
	// attempt to assign same virtual network to previous port                    (assert failure)
	// unassign virtual network from bonded port                                  (assert success)
	// attempt to unassign the same virtual network from the previous bonded port (assert failure)
}
