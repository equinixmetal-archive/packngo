package packngo

import (
	"log"
	"path"
	"testing"
	"time"
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

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         "baremetal_1e",
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)
	dID := d.ID

	// If you need to test this, run a 1e device in your project in a faciltiy
	// and then comment code from MARK_1 to here and uncomment following.
	// Fill the values from your device, project and facility.

	/*

				c, stopRecord := setup(t)
		defer stopRecord()

				dID := "414f52d3-022a-420d-a521-915fdcc66801"
				projectID := "52000fb2-ee46-4673-93a8-de2c2bdba33b"
				fac := "atl1"
				d := &Device{}
				err := fmt.Errorf("hi")

	*/

	d = waitDeviceActive(t, c, dID)

	nType, err := c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeHybrid {
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
	defer deleteProjectVirtualNetwork(t, c, vlan.ID)
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

func TestAccPortL2HybridL3ConvertTypeC1LA(t *testing.T) {
	// PACKNGO_TEST_FACILITY=nrt1 PACKNGO_TEST_ACTUAL_API=1 go test -v -timeout 30m -run=TestAccPortL2HybridL3ConvertType2A
	testL2HybridL3Convert(t, "c1.large.arm")
}

func TestAccPortL2HybridL3ConvertType2(t *testing.T) {
	testL2HybridL3Convert(t, "baremetal_2")
}

func TestAccPortL2HybridL3ConvertType3(t *testing.T) {
	testL2HybridL3Convert(t, "baremetal_3")
}

func TestAccPortL2HybridL3ConvertTypeS(t *testing.T) {
	testL2HybridL3Convert(t, "baremetal_s")
}

func TestAccPortL2HybridL3ConvertN2(t *testing.T) {
	testL2HybridL3Convert(t, "n2.xlarge.x86")
}

func testL2HybridL3Convert(t *testing.T, plan string) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	log.Println("Testing type", plan, "convert to Hybrid L2 and back to L3")

	// MARK_2

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         plan,
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)
	dID := d.ID

	// If you need to test this, run a ${plan} device in your project in a
	// facility,
	// and then comment code from MARK_2 to here and uncomment following.
	// Fill the values from youri testing device, project and facility.

	/*
				c, stopRecord := setup(t)
		defer stopRecord()

				projectID := "52000fb2-ee46-4673-93a8-de2c2bdba33b"
				dID := "b7515d6b-6a86-4830-ab50-9bca3aa51e1c"
				fac := "dfw2"

				//dID := "131dfaf1-e38a-4963-86d6-dbc2531f85d7"
				//fac := "ewr1"

				d := &Device{}
				err := fmt.Errorf("hi")
	*/

	d = waitDeviceActive(t, c, dID)

	nType, err := c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeL3 {
		t.Fatalf("New %s device should be in network type L3", plan)
	}

	d, err = c.DevicePorts.DeviceToNetworkType(d.ID, NetworkTypeHybrid)
	if err != nil {
		t.Fatal(err)
	}

	nType, err = c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeHybrid {
		t.Fatal("the device should now be in network type L2 Bonded")
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
	defer deleteProjectVirtualNetwork(t, c, vlan.ID)

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

	d, err = c.DevicePorts.DeviceToNetworkType(d.ID, NetworkTypeL3)
	if err != nil {
		t.Fatal(err)
	}

	nType, err = c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeL3 {
		t.Fatal("the device should now be back in network type L3")
	}

}

func TestAccPortL2L3ConvertType2A(t *testing.T) {
	// run possible as:
	// PACKNGO_TEST_FACILITY=nrt1 PACKNGO_TEST_ACTUAL_API=1 go test -v -timeout 20m -run=TestAccPort2A
	testL2L3Convert(t, "baremetal_2a")
}

func TestAccPortL2L3ConvertType2(t *testing.T) {
	// PACKNGO_TEST_FACILITY=nrt1 PACKNGO_TEST_ACTUAL_API=1 go test -v -run=TestAccPort2
	testL2L3Convert(t, "baremetal_2")
}

func TestAccPortL2L3ConvertType3(t *testing.T) {
	testL2L3Convert(t, "baremetal_3")
}

func TestAccPortL2L3ConvertTypeS(t *testing.T) {
	testL2L3Convert(t, "baremetal_s")
}

func TestAccPortL2L3ConvertN2(t *testing.T) {
	testL2L3Convert(t, "n2.xlarge.x86")
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
		Facility:     []string{fac},
		Plan:         plan,
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)
	dID := d.ID
	/*

		// If you need to test this, run a ${plan} device in your project in a
		// facility,
		// and then comment code from MARK_2 to here and uncomment following.
		// Fill the values from youri testing device, project and facility.

		c, stopRecord := setup(t)
		defer stopRecord()

		//	dID := "131dfaf1-e38a-4963-86d6-dbc2531f85d7"
		dID := "b7515d6b-6a86-4830-ab50-9bca3aa51e1c"
		projectID := "52000fb2-ee46-4673-93a8-de2c2bdba33b"
		fac := "dfw2"

		d := &Device{}
		err := fmt.Errorf("hi")
	*/

	d = waitDeviceActive(t, c, dID)

	nType, err := c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeL3 {
		t.Fatalf("New %s device should be in network type L3", plan)
	}

	d, err = c.DevicePorts.DeviceToNetworkType(d.ID, NetworkTypeL2Bonded)
	if err != nil {
		t.Fatal(err)
	}

	nType, err = c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeL2Bonded {
		t.Fatal("the device should now be in network type L2 Bonded")
	}

	bond0, err := c.DevicePorts.GetPortByName(d.ID, "bond0")
	if err != nil {
		t.Fatal(err)
	}

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
	defer deleteProjectVirtualNetwork(t, c, vlan.ID)

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

	if len(p.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to the port at this time")
	}

	d, err = c.DevicePorts.DeviceToNetworkType(d.ID, NetworkTypeL3)
	if err != nil {
		t.Fatal(err)
	}

	nType, err = c.DevicePorts.DeviceNetworkType(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if nType != NetworkTypeL3 {
		t.Fatal("the device now should be back in network type L3")
	}

}

func deviceToNetworkType(t *testing.T, c *Client, deviceID, targetNetworkType string) {
	oldt, err := c.DevicePorts.DeviceNetworkType(deviceID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Converting", oldt, "=>", targetNetworkType, "...")
	_, err = c.DevicePorts.DeviceToNetworkType(deviceID, targetNetworkType)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(oldt, "=>", targetNetworkType, "OK")
	time.Sleep(15 * time.Second)
}

/*
func TestXXX(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, stopRecord := setup(t)
defer stopRecord()
	//	deviceID := "131dfaf1-e38a-4963-86d6-dbc2531f85d7"
	deviceID := "b7515d6b-6a86-4830-ab50-9bca3aa51e1c"
	d, _, err := c.Devices.Get(deviceID, nil)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(d)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL3)
	deviceToNetworkType(t, c, deviceID, NetworkTypeHybrid)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL3)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Individual)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL3)

}
*/

func TestAccPortNetworkStateTransitions(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     "networktypetest",
		Facility:     []string{fac},
		Plan:         "m1.xlarge.x86",
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}
	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)
	deviceID := d.ID

	d = waitDeviceActive(t, c, deviceID)

	networkType, err := d.GetNetworkType()
	if err != nil {
		t.Fatal(err)
	}

	if networkType != NetworkTypeL2Bonded {
		deviceToNetworkType(t, c, deviceID, NetworkTypeL2Bonded)
	}

	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Individual)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL3)
	deviceToNetworkType(t, c, deviceID, NetworkTypeHybrid)

	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Bonded)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL3)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Bonded)

	deviceToNetworkType(t, c, deviceID, NetworkTypeHybrid)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Individual)
	deviceToNetworkType(t, c, deviceID, NetworkTypeHybrid)

	deviceToNetworkType(t, c, deviceID, NetworkTypeL3)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Individual)
	deviceToNetworkType(t, c, deviceID, NetworkTypeL2Bonded)
}

func TestAccPortNativeVlan(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     "networktypetest",
		Facility:     []string{fac},
		Plan:         "baremetal_2",
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}
	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	deviceID := d.ID

	d = waitDeviceActive(t, c, deviceID)
	defer deleteDevice(t, c, d.ID, false)

	deviceToNetworkType(t, c, deviceID, NetworkTypeHybrid)

	vncr := VirtualNetworkCreateRequest{
		ProjectID: projectID,
		Facility:  fac,
	}

	vlan1, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteProjectVirtualNetwork(t, c, vlan1.ID)
	vlan2, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteProjectVirtualNetwork(t, c, vlan2.ID)

	eth1, err := c.DevicePorts.GetPortByName(d.ID, "eth1")
	if err != nil {
		t.Fatal(err)
	}
	if eth1.NativeVirtualNetwork != nil {
		t.Fatal("Native virtual network on fresh device should be nil")
	}
	par1 := PortAssignRequest{
		PortID:           eth1.ID,
		VirtualNetworkID: vlan1.ID}
	_, _, err = c.DevicePorts.Assign(&par1)
	if err != nil {
		t.Fatal(err)
	}
	par2 := PortAssignRequest{
		PortID:           eth1.ID,
		VirtualNetworkID: vlan2.ID}
	_, _, err = c.DevicePorts.Assign(&par2)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.DevicePorts.AssignNative(&par1)
	if err != nil {
		t.Fatal(err)
	}
	eth1, err = c.DevicePorts.GetPortByName(deviceID, "eth1")
	if err != nil {
		t.Fatal(err)
	}
	if eth1.NativeVirtualNetwork != nil {
		if path.Base(eth1.NativeVirtualNetwork.Href) != vlan1.ID {
			t.Fatal("Wrong native virtual network at the test device")
		}
	} else {
		t.Fatal("No native virtual network at the test device")
	}
	_, _, err = c.DevicePorts.UnassignNative(eth1.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.DevicePorts.Unassign(&par2)
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := c.DevicePorts.Unassign(&par1)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(p)
}
