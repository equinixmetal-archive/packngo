package packngo

import (
	"fmt"
	"testing"
	"time"
)

func waitDeviceActive(id string, c *Client) (*Device, error) {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		d, _, err := c.Devices.Get(id)
		if err != nil {
			return nil, err
		}
		if d.State == "active" {
			return d, nil
		}
	}
	return nil, fmt.Errorf("device %s is still not active after timeout", id)
}

func TestAccDeviceBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     "ewr1",
		Plan:         "baremetal_0",
		OS:           "ubuntu_16_04",
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	dID := d.ID

	d, err = waitDeviceActive(dID, c)
	if err != nil {
		t.Fatal(err)
	}

	if len(d.RootPassword) == 0 {
		t.Fatal("root_password is empty or non-existent")
	}

	newHN := randString8()
	ur := DeviceUpdateRequest{Hostname: newHN}

	newD, _, err := c.Devices.Update(dID, &ur)
	if err != nil {
		t.Fatal(err)
	}

	if newD.Hostname != newHN {
		t.Fatalf("hostname of test device should be %s, but is %s", newHN, newD.Hostname)
	}

	_, err = c.Devices.Delete(dID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccDevicePXE(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()
	pxeURL := "https://boot.netboot.xyz"

	cr := DeviceCreateRequest{
		Hostname:      "pxe-" + hn,
		Facility:      "ewr1",
		Plan:          "baremetal_0",
		ProjectID:     projectID,
		BillingCycle:  "hourly",
		OS:            "custom_ipxe",
		IPXEScriptURL: "https://boot.netboot.xyz",
		AlwaysPXE:     true,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	d, err = waitDeviceActive(d.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	if !d.AlwaysPXE {
		t.Fatal("always_pxe should be set")
	}
	if d.IPXEScriptURL != pxeURL {
		t.Fatalf("ipxe_script_url should be %s", pxeURL)
	}
	_, err = c.Devices.Delete(d.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccDeviceAssignIP(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()

	testFac := "ewr1"

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     testFac,
		Plan:         "baremetal_0",
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           "ubuntu_16_04",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	d, err = waitDeviceActive(d.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	req := IPReservationRequest{
		Type:     "public_ipv4",
		Quantity: 1,
		Comments: "packngo test",
		Facility: testFac,
	}

	af, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}

	assignment, _, err := c.DeviceIPs.Assign(d.ID, af)
	if err != nil {
		t.Fatal(err)
	}

	d, _, err = c.Devices.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	// If the quantity in the IPReservationRequest is >1, this test won't work.
	// The assignment CIDR would then have to be extracted from the reserved
	// block.
	reservation, _, err := c.ProjectIPs.GetByCIDR(projectID, af.Address)
	if err != nil {
		t.Fatal(err)
	}

	if len(reservation.Assignments) != 1 {
		t.Fatalf("reservation %s should have exactly 1 assignment", reservation)
	}

	if reservation.Assignments[0].Href != assignment.Href {
		t.Fatalf("assignment %s should be listed in reservation resource %s",
			assignment.Href, reservation)

	}

	func() {
		for _, ipa := range d.Network {
			if ipa.Href == assignment.Href {
				return
			}
		}
		t.Fatalf("assignment %s should be listed in device %s", assignment, d)
	}()

	if assignment.AssignedTo.Href != d.Href {
		t.Fatalf("device %s should be listed in assignment %s",
			d, assignment)
	}

	_, err = c.DeviceIPs.Unassign(assignment.ID)
	if err != nil {
		t.Fatal(err)
	}

	// reload reservation, now without any assignment
	reservation, _, err = c.ProjectIPs.Get(reservation.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(reservation.Assignments) != 0 {
		t.Fatalf("reservation %s shoud be without assignments. Was %v",
			reservation, reservation.Assignments)
	}

	// reload device, now without the assigned floating IP
	d, _, err = c.Devices.Get(d.ID)
	if err != nil {
		t.Fatal(err)
	}

	for _, ipa := range d.Network {
		if ipa.Href == assignment.Href {
			t.Fatalf("assignment %s shoud be not listed in device %s anymore",
				assignment, d)
		}
	}

	_, err = c.Devices.Delete(d.ID)
	if err != nil {
		t.Fatal(err)
	}

}
