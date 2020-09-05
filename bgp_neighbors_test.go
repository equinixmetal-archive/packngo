package packngo

import (
	"log"
	"testing"
)

func createBGPDevice(t *testing.T, c *Client, projectID string) *Device {
	hn := randString8()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{testFacility()},
		Plan:         "baremetal_0",
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           "ubuntu_16_04",
	}
	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	d = waitDeviceActive(t, c, d.ID)

	aTrue := true
	_, _, err = c.BGPSessions.Create(d.ID,
		CreateBGPSessionRequest{
			AddressFamily: "ipv4",
			DefaultRoute:  &aTrue})
	if err != nil {
		t.Fatal(err)
	}

	return d
}

func TestAccBGPNeighbors(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	configRequest := CreateBGPConfigRequest{
		DeploymentType: "local",
		Asn:            65000,
		Md5:            "c3RhY2twb2ludDIwMTgK",
	}

	_, err := c.BGPConfig.Create(projectID, configRequest)
	if err != nil {
		t.Fatal(err)
	}

	d := createBGPDevice(t, c, projectID)

	defer c.Devices.Delete(d.ID, false)

	bgpNeighbors, _, err := c.Devices.ListBGPNeighbors(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Testdevice BGP neighbors: \n %+v", bgpNeighbors)

	if len(bgpNeighbors) < 1 {
		t.Fatal("Device with BGP session should have at least one listed BGP neighbor - itself")
	}
	n0 := bgpNeighbors[0]
	if len(n0.RoutesIn) < 1 {
		t.Fatal("Device with BGP session should have at least one route in")
	}
	if len(n0.RoutesOut) < 1 {
		t.Fatal("Device with BGP session should have at least one route out")
	}

}
