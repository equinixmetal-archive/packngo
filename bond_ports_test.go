package packngo

import (
	"log"
	"testing"
	"time"
)

func deviceBondToNetworkType(t *testing.T, c *Client, deviceID, bondPortName, targetNetworkType string) {
	d, _, err := c.Devices.Get(deviceID, nil)
	if err != nil {
		t.Fatal(err)
	}
	oldType := d.GetBondNetworkType(bondPortName)
	log.Println("Converting", bondPortName, oldType, "=>", targetNetworkType, "...")
	_, err = c.DevicePorts.BondToNetworkType(deviceID, bondPortName, targetNetworkType)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(oldType, "=>", targetNetworkType, "OK")
	// why sleep here??
	time.Sleep(5 * time.Second)
}

func TestAccPortBondNetworkStateTransitions(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     "bondNetworkTypeTest",
		Facility:     []string{fac},
		Plan:         "c3.small.x86",
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

	bond := "bond0"

	networkType := d.GetBondNetworkType(bond)
	if networkType != NetworkTypeL3 {
		t.Fatal("network_type should be 'layer3'")
	}

	if networkType != NetworkTypeL2Bonded {
		deviceBondToNetworkType(t, c, deviceID, bond, NetworkTypeL2Bonded)
	}

	deviceBondToNetworkType(t, c, deviceID, bond, NetworkTypeL2Individual)
}

func TestAccPortBondNetworkStateHybridBonded(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     "NetworkTypeHybridBondedTest",
		Facility:     []string{fac},
		Plan:         "c3.small.x86",
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

	bond := "bond0"

	networkType := d.GetBondNetworkType(bond)
	if networkType != NetworkTypeL3 {
		t.Fatal("network_type should be 'layer3'")
	}

	testDesc := "test_desc_" + randString8()

	vncr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: testDesc,
		Facility:    testFacility(),
	}

	vlan, _, err := c.ProjectVirtualNetworks.Create(&vncr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteProjectVirtualNetwork(t, c, vlan.ID)

	bondPort, _ := d.GetPortByName(bond)

	par := PortAssignRequest{
		PortID:           bondPort.ID,
		VirtualNetworkID: vlan.ID}

	p, _, err := c.DevicePorts.Assign(&par)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v\n", p)

	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	newBondNetworkType := d.GetBondNetworkType(bond)
	if newBondNetworkType != NetworkTypeHybridBonded {
		t.Fatalf("After bond0 is in L3 and VLAN is assigned to bond0, it's network type should be %s. Was %s", NetworkTypeHybridBonded, newBondNetworkType)
	}

	_, _, err = c.DevicePorts.Unassign(&par)
	if err != nil {
		t.Fatal(err)
	}

}
