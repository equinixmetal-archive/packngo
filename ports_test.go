package packngo

import (
	"fmt"
	"log"
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

	/*

		c, projectID, teardown := setupWithProject(t)
		defer teardown()

		hn := randString8()

		fac := testFacility()

		cr := DeviceCreateRequest{
			Hostname:     hn,
			Facility:     fac,
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
	*/

	// If you need to test this, run a 1e device in your project in a faciltiy
	// and then comment code from MARK_1 to here and uncomment following.
	// Fill the values from your device, project and facility.

	c := setup(t)

	dID := "414f52d3-022a-420d-a521-915fdcc66801"
	projectID := "52000fb2-ee46-4673-93a8-de2c2bdba33b"
	fac := "atl1"
	d := &Device{}
	err := fmt.Errorf("hi")

	d, err = waitDeviceActive(dID, c)
	if err != nil {
		t.Fatal(err)
	}

	nType, err := c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkHybrid {
		t.Fatal("New 1E device should be in Hybrid Network Type")
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
		Facility:  fac,
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.ProjectVirtualNetworks.Delete(vlan.ID)
	par := PortAssignRequest{
		PortID:           eth1.ID,
		VirtualNetworkID: vlan.ID}

	p, _, err := c.DevicePorts.Assign(&par)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.AttachedVirtualNetworks) != 1 {
		t.Fatal("Exactly one vlan should be attached to a eth1 at this point")
	}

	if path.Base(p.AttachedVirtualNetworks[0].Href) != vlan.ID {
		t.Fatal("mismatch in the UUID of the assigned VLAN")
	}

	p, _, err = c.DevicePorts.Unassign(&par)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to the port at this time")
	}
}

func TestAccPort2A(t *testing.T) {
	// run possible as:
	// PACKNGO_TEST_FACILITY=nrt1 PACKNGO_TEST_ACTUAL_API=1 go test -v -timeout 20m -run=TestAccPort2A
	testL2L3Convert(t, "baremetal_2a")
}

func TestAccPort2(t *testing.T) {
	// PACKNGO_TEST_FACILITY=nrt1 PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccPort2
	testL2L3Convert(t, "baremetal_2")
}

func TestAccPort3(t *testing.T) {
	testL2L3Convert(t, "baremetal_3")
}

func TestAccPortS(t *testing.T) {
	testL2L3Convert(t, "baremetal_s")
}

func testL2L3Convert(t *testing.T, plan string) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	log.Println("Testing type", plan, "convert to L2 and back to L3")

	// MARK_2

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     fac,
		Plan:         plan,
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

	// If you need to test this, run a ${plan} device in your project in a
	// facility,
	// and then comment code from MARK_2 to here and uncomment following.
	// Fill the values from youri testing device, project and facility.

	/*

		c := setup(t)

		dID := "873fecd1-85a0-4998-ae61-13cfb4f199a8"
		projectID := "52000fb2-ee46-4673-93a8-de2c2bdba33b"
		fac := "nrt1"

		d := &Device{}
		err := fmt.Errorf("hi")
	*/

	d, err = waitDeviceActive(dID, c)
	if err != nil {
		t.Fatal(err)
	}

	nType, err := c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkL3 {
		t.Fatalf("New %s device should be in network type L3", plan)
	}

	d, err = c.DevicePorts.DeviceToLayerTwo(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	nType, err = c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkL2Bonded {
		t.Fatalf("the device shoud be in network type L2 Bonded", plan)
	}

	bond0, err := c.DevicePorts.GetBondPort(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	/*

		bond0, err := c.DevicePorts.GetBondPort(d.ID)
		if err != nil {
			t.Fatal(err)
		}

		bond0, _, err = c.DevicePorts.ConvertToLayerTwo(bond0.ID)
		if err != nil {
			t.Fatal(err)
		}
	*/
	// would be cool to test if the device is indeed in L2 networking mode
	// at this point but I don't know a way to tell from the API

	if len(bond0.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to a bond0 in the begining of this test")
	}

	vncr := VirtualNetworkCreateRequest{
		ProjectID: projectID,
		Facility:  fac,
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.ProjectVirtualNetworks.Delete(vlan.ID)

	par := PortAssignRequest{
		PortID:           bond0.ID,
		VirtualNetworkID: vlan.ID}
	p, _, err := c.DevicePorts.Assign(&par)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.AttachedVirtualNetworks) != 1 {
		t.Fatal("Exactly one vlan should be attached to a bond0 at this point")
	}

	if path.Base(p.AttachedVirtualNetworks[0].Href) != vlan.ID {
		t.Fatal("mismatch in the UUID of the assigned VLAN")
	}

	p, _, err = c.DevicePorts.Unassign(&par)
	if err != nil {
		t.Fatal(err)
	}

	/*

		bond0, _, err = c.DevicePorts.Bond(&BondRequest{
			PortID: bond0.ID, BulkEnable: false})

		if err != nil {
			t.Fatal(err)
		}
	*/

	if len(p.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to the port at this time")
	}

	d, err = c.DevicePorts.DeviceToLayerThree(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	nType, err = c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkL3 {
		t.Fatalf("the %s device should be back in network type L3", plan)
	}

}
