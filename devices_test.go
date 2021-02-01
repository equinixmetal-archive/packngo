package packngo

import (
	"fmt"
	"path"
	"reflect"
	"testing"
	"time"
)

var (
	// managementIPS is a block of public ipv4 and ipv6, with private ipv4. all
	// with management enabled. This is standard in Layer3 configurations.
	managementIPS = []*IPAddressAssignment{
		{
			IpAddressCommon: IpAddressCommon{
				Management:    true,
				AddressFamily: 4,
				Public:        true,
			},
		},
		{
			IpAddressCommon: IpAddressCommon{
				Management:    true,
				AddressFamily: 6,
				Public:        true,
			},
		},
		{
			IpAddressCommon: IpAddressCommon{
				Management:    true,
				AddressFamily: 4,
				Public:        false,
			},
		},
	}
)

func waitDeviceActive(t *testing.T, c *Client, id string) *Device {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		d, _, err := c.Devices.Get(id, nil)
		if err != nil {
			t.Fatal(err)
			return nil
		}
		if d.State == "active" {
			return d
		}
		if d.State == "failed" {
			t.Fatal(fmt.Errorf("device %s provisioning failed", id))
			return nil
		}
	}

	t.Fatal(fmt.Errorf("device %s is still not active after timeout", id))
	return nil
}

func deleteDevice(t *testing.T, c *Client, id string, force bool) {
	if _, err := c.Devices.Delete(id, force); err != nil {
		t.Fatal(err)
	}
}

func deleteSSHKey(t *testing.T, c *Client, id string) {
	if _, err := c.SSHKeys.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func deleteVolume(t *testing.T, c *Client, id string) {
	if _, err := c.Volumes.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func deleteVolumeAttachments(t *testing.T, c *Client, id string) {
	if _, err := c.VolumeAttachments.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func deleteProjectIP(t *testing.T, c *Client, id string) {
	if _, err := c.ProjectIPs.Remove(id); err != nil {
		t.Fatal(err)
	}
}

func deleteProjectVirtualNetwork(t *testing.T, c *Client, id string) {
	if _, err := c.ProjectVirtualNetworks.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccDeviceUpdate(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
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

	d = waitDeviceActive(t, c, dID)

	if len(d.RootPassword) == 0 {
		t.Fatal("root_password is empty or non-existent")
	}
	newHN := randString8()
	ur := DeviceUpdateRequest{Hostname: &newHN}

	newD, _, err := c.Devices.Update(dID, &ur)
	if err != nil {
		t.Fatal(err)
	}

	if newD.Hostname != newHN {
		t.Fatalf("hostname of test device should be %s, but is %s", newHN, newD.Hostname)
	}
	for _, ipa := range newD.Network {
		if !ipa.Management {
			t.Fatalf("management flag for all the IP addresses in a new device should be True: was %s", ipa)
		}
	}
}

func TestAccDeviceBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
		Description:  "test",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	dID := d.ID

	d = waitDeviceActive(t, c, dID)

	if len(d.ShortID) == 0 {
		t.Fatal("Device should have shortID")
	}
	if len(d.SwitchUUID) == 0 {
		t.Fatal("Device should have switch UUID")
	}
	networkType := d.GetNetworkType()
	if networkType != NetworkTypeL3 {
		t.Fatal("network_type should be 'layer3'")
	}

	if d.User != "root" {
		t.Fatal("user should be 'root'")
	}
	if d.Description == nil || *d.Description != cr.Description {
		t.Fatal("description is empty or non-existent")
	}
	if len(d.RootPassword) == 0 {
		t.Fatal("root_password is empty or non-existent")
	}
	networkInfo := d.GetNetworkInfo()

	for _, ipa := range d.Network {
		if !ipa.Management {
			t.Fatalf("management flag for all the IP addresses in a new device should be True: was %s", ipa)
		}
		if ipa.Public && (ipa.AddressFamily == 4) {
			if ipa.Address != networkInfo.PublicIPv4 {
				t.Fatalf("strange public IPv4 from GetNetworkInfo, should be %s, is %s", ipa.Address, networkInfo.PublicIPv4)

			}
		}
	}
	dl, _, err := c.Devices.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(dl) != 1 {
		t.Fatalf("Device List should contain exactly one device, was: %v", dl)
	}

}

func TestAccDevicePXE(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()
	pxeURL := "https://boot.netboot.xyz"
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:      "pxe-" + hn,
		Facility:      []string{fac},
		Plan:          testPlan(),
		ProjectID:     projectID,
		BillingCycle:  "hourly",
		OS:            "custom_ipxe",
		IPXEScriptURL: pxeURL,
		AlwaysPXE:     true,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	// Check that settings were persisted
	if !d.AlwaysPXE {
		t.Fatal("always_pxe should be true")
	}
	if d.IPXEScriptURL != pxeURL {
		t.Fatalf("ipxe_script_url should be \"%s\"", pxeURL)
	}

	// Check that we can update PXE options
	pxeURL = "http://boot.netboot.xyz"
	bFalse := false
	d, _, err = c.Devices.Update(d.ID,
		&DeviceUpdateRequest{
			AlwaysPXE: &bFalse,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if d.AlwaysPXE {
		t.Fatalf("always_pxe should have been updated to false")
	}
	d, _, err = c.Devices.Update(d.ID,
		&DeviceUpdateRequest{
			IPXEScriptURL: &pxeURL,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if d.IPXEScriptURL != pxeURL {
		t.Fatalf("ipxe_script_url should have been updated to \"%s\"", pxeURL)
	}
}

func TestAccDeviceAssignGlobalIP(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           testOS,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	req := IPReservationRequest{
		Type:        "global_ipv4",
		Quantity:    1,
		Description: "packngo test",
	}

	reservation, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}

	af := AddressStruct{Address: fmt.Sprintf("%s/%d", reservation.Address, reservation.CIDR)}

	assignment, _, err := c.DeviceIPs.Assign(d.ID, &af)
	if err != nil {
		t.Fatal(err)
	}

	if assignment.Management {
		t.Error("Management flag for assignment resource must be False")
	}

	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// If the Quantity in the IPReservationRequest is >1, this test won't work.
	// The assignment CIDR would then have to be extracted from the reserved
	// block.
	reservation, _, err = c.ProjectIPs.Get(reservation.ID, nil)
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

	for _, ipa := range d.Network {
		if ipa.Href == assignment.Href {
			return
		}
	}
	t.Fatalf("assignment %s should be listed in device %s", assignment, d)

	if assignment.AssignedTo.Href != d.Href {
		t.Fatalf("device %s should be listed in assignment %s",
			d, assignment)
	}

	_, err = c.DeviceIPs.Unassign(assignment.ID)
	if err != nil {
		t.Fatal(err)
	}

	// reload reservation, now without any assignment
	reservation, _, err = c.ProjectIPs.Get(reservation.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(reservation.Assignments) != 0 {
		t.Fatalf("reservation %s shoud be without assignments. Was %v",
			reservation, reservation.Assignments)
	}

	// reload device, now without the assigned floating IP
	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, ipa := range d.Network {
		if ipa.Href == assignment.Href {
			t.Fatalf("assignment %s shoud be not listed in device %s anymore",
				assignment, d)
		}
	}
}

func TestAccDeviceCreateWithReservedIP(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()

	fac := testFacility()

	req := IPReservationRequest{
		Type:        "public_ipv4",
		Quantity:    2,
		Description: "packngo test",
		Facility:    &fac,
	}

	reservation, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteProjectIP(t, c, reservation.ID)

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           testOS,
		IPAddresses: []IPAddressCreateRequest{
			// NOTE: only one public IPv4 entry is allowed here
			{AddressFamily: 4, Public: false},
			{AddressFamily: 4, Public: true,
				Reservations: []string{reservation.ID}, CIDR: 31},
		},
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	d = waitDeviceActive(t, c, d.ID)

	defer deleteDevice(t, c, d.ID, false)

	reservation, _, err = c.ProjectIPs.Get(reservation.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(reservation.Assignments) != 1 {
		t.Fatalf("reservation %s should have exactly 1 assignment", reservation)
	}
}

func TestAccDeviceAssignIP(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()

	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           testOS,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	req := IPReservationRequest{
		Type:        PublicIPv4,
		Quantity:    1,
		Description: "packngo test",
		Facility:    &fac,
	}

	reservation, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}

	af := AddressStruct{Address: fmt.Sprintf("%s/%d", reservation.Address, reservation.CIDR)}

	assignment, _, err := c.DeviceIPs.Assign(d.ID, &af)
	if err != nil {
		t.Fatal(err)
	}

	if assignment.Management {
		t.Error("Management flag for assignment resource must be False")
	}

	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// check that the IP assignment is retrievable via the IP-by-device endpoint
	assignments, _, err := c.DeviceIPs.List(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	var matchedAssignment bool
	for _, ip := range assignments {
		if ip.String() == assignment.String() {
			matchedAssignment = true
			break
		}
	}
	if !matchedAssignment {
		t.Fatal("newly assigned IP not found")
	}

	// If the Quantity in the IPReservationRequest is >1, this test won't work.
	// The assignment CIDR would then have to be extracted from the reserved
	// block.
	reservation, _, err = c.ProjectIPs.Get(reservation.ID, nil)
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

	for _, ipa := range d.Network {
		if ipa.Href == assignment.Href {
			return
		}
	}
	t.Fatalf("assignment %s should be listed in device %s", assignment, d)

	if assignment.AssignedTo.Href != d.Href {
		t.Fatalf("device %s should be listed in assignment %s",
			d, assignment)
	}

	_, err = c.DeviceIPs.Unassign(assignment.ID)
	if err != nil {
		t.Fatal(err)
	}

	// reload reservation, now without any assignment
	reservation, _, err = c.ProjectIPs.Get(reservation.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(reservation.Assignments) != 0 {
		t.Fatalf("reservation %s shoud be without assignments. Was %v",
			reservation, reservation.Assignments)
	}

	// reload device, now without the assigned floating IP
	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, ipa := range d.Network {
		if ipa.Href == assignment.Href {
			t.Fatalf("assignment %s shoud be not listed in device %s anymore",
				assignment, d)
		}
	}
}

func TestAccDeviceAttachVolumeForceDelete(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           testOS,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	// defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	vcr := VolumeCreateRequest{
		Size:         10,
		BillingCycle: "hourly",
		PlanID:       "storage_1",
		FacilityID:   testFacility(),
	}

	v, _, err := c.Volumes.Create(&vcr, projectID)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVolume(t, c, v.ID)

	v, err = waitVolumeActive(v.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.VolumeAttachments.Create(v.ID, d.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.Volumes.Get(v.ID,
		&GetOptions{Includes: []string{"attachments.device"}})
	if err != nil {
		t.Fatal(err)
	}

	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	defer deleteDevice(t, c, d.ID, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccDeviceAttachVolume(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           testOS,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	vcr := VolumeCreateRequest{
		Size:         10,
		BillingCycle: "hourly",
		PlanID:       "storage_1",
		FacilityID:   testFacility(),
	}

	v, _, err := c.Volumes.Create(&vcr, projectID)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteVolume(t, c, v.ID)

	v, err = waitVolumeActive(v.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	a, _, err := c.VolumeAttachments.Create(v.ID, d.ID)
	if err != nil {
		t.Fatal(err)
	}

	if path.Base(a.Volume.Href) != v.ID {
		t.Fatalf("wrong volume href in the attachment: %s, should be %s", a.Volume.Href, v.ID)
	}

	if path.Base(a.Device.Href) != d.ID {
		t.Fatalf("wrong device href in the attachment: %s, should be %s", a.Device.Href, d.ID)
	}

	v, _, err = c.Volumes.Get(v.ID,
		&GetOptions{Includes: []string{"attachments.device"}})
	if err != nil {
		t.Fatal(err)
	}

	d, _, err = c.Devices.Get(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if v.Attachments[0].Device.ID != d.ID {
		t.Fatalf("wrong device linked in volume attachment: %s, should be %s", v.Attachments[0].Device.ID, d.ID)
	}
	if path.Base(d.Volumes[0].Href) != v.ID {
		t.Fatalf("wrong volume linked in device.volumes: %s, should be %s", d.Volumes[0].Href, v.ID)
	}

	defer deleteVolumeAttachments(t, c, a.ID)
}

func TestAccDeviceSpotInstance(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()

	testSPM := 0.04
	testTerm := &Timestamp{Time: time.Now().Add(time.Hour - (time.Minute * 10))}
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:        hn,
		Facility:        []string{fac},
		Plan:            testPlan(),
		OS:              "coreos_stable",
		ProjectID:       projectID,
		BillingCycle:    "hourly",
		SpotInstance:    true,
		SpotPriceMax:    testSPM,
		TerminationTime: testTerm,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	if !d.SpotInstance {
		t.Fatal("spot_instance is false, should be true")
	}

	if d.SpotPriceMax != testSPM {
		t.Fatalf("spot_price_max is %f, should be %f", d.SpotPriceMax, testSPM)
	}

	if !d.TerminationTime.Time.Truncate(time.Minute).Equal(testTerm.Time.Truncate(time.Minute)) {
		t.Fatalf("termination_time is %s, should be %s",
			d.TerminationTime.Time.Local(), testTerm.Time.Local())
	}
}

func TestAccDeviceCustomData(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	fac := testFacility()

	customData := map[string]interface{}{"custom1": "data", "custom2": map[string]interface{}{"nested": "data"}}

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
		CustomData:   &customData,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	dID := d.ID

	_ = waitDeviceActive(t, c, dID)

	device, _, err := c.Devices.Get(dID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(customData, device.CustomData) {
		t.Fatalf("CustomData of new device should be %+v, was %+v", customData, device.CustomData)
	}

	updateCustomData := map[string]interface{}{"hi": "earth"}
	_, _, err = c.Devices.Update(dID, &DeviceUpdateRequest{
		CustomData: &updateCustomData,
	})
	if err != nil {
		t.Fatal(err)
	}

	device, _, err = c.Devices.Get(dID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(updateCustomData, device.CustomData) {
		t.Fatalf("CustomData of updated device should be %+v, was %+v", updateCustomData, device.CustomData)
	}

	updateCustomData = map[string]interface{}{}
	_, _, err = c.Devices.Update(dID, &DeviceUpdateRequest{
		CustomData: &updateCustomData,
	})
	if err != nil {
		t.Fatal(err)
	}
	// It's not possible to unset customdata:
	// https://github.com/packethost/packngo/pull/225#issuecomment-744331937
	/*
		device, _, err = c.Devices.Get(dID, nil)
		if err != nil {
			t.Fatal(err)
		}
		if device.CustomData != nil {
			t.Fatalf("CustomData of updated device should be erased to %+v, was %+v", updateCustomData, device.CustomData)
		}
	*/
}

func TestAccListDeviceEvents(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	initialCustomData := `{"hello":"world"}`
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
		CustomData:   initialCustomData,
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	d = waitDeviceActive(t, c, d.ID)

	events, _, err := c.Devices.ListEvents(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) == 0 {
		t.Fatal("Device events not returned")
	}
}

func TestAccDeviceSSHKeys(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()
	userKey := createKey(t, c, "")
	defer deleteSSHKey(t, c, userKey.ID)
	projectKey := createKey(t, c, projectID)
	defer deleteSSHKey(t, c, projectKey.ID)

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{testFacility()},
		Plan:         testPlan(),
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
	_ = waitDeviceActive(t, c, dID)

	d, _, err = c.Devices.Get(dID, &GetOptions{Includes: []string{"ssh_keys"}})
	if err != nil {
		t.Fatal(err)
	}
	userKeyIn := false
	projectKeyIn := false
	for _, k := range d.SSHKeys {
		if k.ID == userKey.ID {
			userKeyIn = true
		}
		if k.ID == projectKey.ID {
			projectKeyIn = true
		}
	}
	if !userKeyIn {
		t.Fatalf("User SSH Key %+v is not present at device", userKey)
	}
	if !projectKeyIn {
		t.Fatalf("Project SSH Key %+v is not present at device", projectKey)
	}
}

func TestAccDeviceListedSSHKeys(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	hn := randString8()
	userKey := createKey(t, c, "")
	defer deleteSSHKey(t, c, userKey.ID)
	projectKey := createKey(t, c, projectID)
	defer deleteSSHKey(t, c, projectKey.ID)
	projectKey2 := createKey(t, c, projectID)
	defer deleteSSHKey(t, c, projectKey2.ID)
	cr := DeviceCreateRequest{
		Hostname:       hn,
		Facility:       []string{testFacility()},
		Plan:           testPlan(),
		OS:             testOS,
		ProjectID:      projectID,
		BillingCycle:   "hourly",
		ProjectSSHKeys: []string{projectKey.ID},
	}
	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)
	dID := d.ID
	_ = waitDeviceActive(t, c, dID)

	d, _, err = c.Devices.Get(dID, &GetOptions{Includes: []string{"ssh_keys"}})
	if err != nil {
		t.Fatal(err)
	}
	userKeyIn := false
	projectKeyIn := false
	projectKey2In := false
	for _, k := range d.SSHKeys {
		if k.ID == userKey.ID {
			userKeyIn = true
		}
		if k.ID == projectKey.ID {
			projectKeyIn = true
		}
		if k.ID == projectKey2.ID {
			projectKey2In = true
		}
	}
	if userKeyIn {
		t.Fatalf("User SSH Key %+v should not be at device", userKey)
	}
	if !projectKeyIn {
		t.Fatalf("Project SSH Key %+v is not present at device", projectKey)
	}
	if projectKey2In {
		t.Fatalf("Project SSH Key %+v is not present at device", projectKey2)
	}
}

func TestAccDeviceCreateFacilities(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	facilities := []string{"nrt1", "dfw1", "fra1"}

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Plan:         testPlan(),
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
		Facility:     facilities,
		Features:     map[string]string{},
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	dID := d.ID

	d = waitDeviceActive(t, c, dID)

	placedInRequestedFacility := false
	for _, fac := range facilities {
		if d.Facility.Code == fac {
			placedInRequestedFacility = true
		}
	}
	if !placedInRequestedFacility {
		t.Fatal("Did not properly assign facility to device")
	}

}

func TestAccDeviceIPAddresses(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()
	fac := testFacility()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     []string{fac},
		Plan:         testPlan(),
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
		IPAddresses: []IPAddressCreateRequest{
			{AddressFamily: 4, Public: false},
		},
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)

	dID := d.ID

	d = waitDeviceActive(t, c, dID)
	networkType := d.GetNetworkType()
	if networkType != NetworkTypeL3 {
		t.Fatal("network_type should be 'layer3'")
	}

	if len(d.RootPassword) == 0 {
		t.Fatal("root_password is empty or non-existent")
	}

	ni := d.GetNetworkInfo()
	if ni.PrivateIPv4 == "" {
		t.Fatal("Device should have private IPv4 present")
	}
	if ni.PublicIPv4 != "" {
		t.Fatal("Device should not have public IPv4 present")
	}

	dl, _, err := c.Devices.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(dl) != 1 {
		t.Fatalf("Device List should contain exactly one device, was: %v", dl)
	}

}

func TestDevice_NumOfBonds(t *testing.T) {
	tests := []struct {
		name  string
		ports []Port
		want  int
	}{
		{
			name:  "NumOfBonds-none",
			ports: []Port{},
			want:  0,
		}, {
			name:  "NumOfBonds-missed",
			ports: []Port{{Type: "foo"}},
			want:  0,
		}, {
			name:  "NumOfBonds-one",
			ports: []Port{{Type: "NetworkBondPort"}},
			want:  1,
		}, {
			name:  "NumOfBonds-two",
			ports: []Port{{Type: "NetworkBondPort"}, {Type: "NetworkBondPort"}},
			want:  2,
		}, {
			name:  "NumOfBonds-mixed-one",
			ports: []Port{{Type: "NetworkBondPort"}, {Type: "foo"}},
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				NetworkPorts: tt.ports,
			}
			if got := d.NumOfBonds(); got != tt.want {
				t.Errorf("Device.NumOfBonds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDevice_GetNetworkType(t *testing.T) {
	type fields struct {
		Network      []*IPAddressAssignment
		Plan         *Plan
		NetworkPorts []Port
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "GetNetworkType_bm0_provisioning", // t1.small.x86
			fields: fields{Plan: &Plan{Slug: "baremetal_0"},
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer2-bonded",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}},
			},
			want: NetworkTypeL3,
		},
		{
			name: "GetNetworkType_bm0_ready", // t1.small.x86 post-provision
			fields: fields{Plan: &Plan{Slug: "baremetal_0"},
				Network: managementIPS,
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer3",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}},
			},
			want: NetworkTypeL3,
		},
		{
			name: "GetNetworkType_bm1", // c1.small.x86
			fields: fields{Plan: &Plan{Slug: "baremetal_1"},
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer2-bonded",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}}},
			want: NetworkTypeL3,
		},
		{
			name: "GetNetworkType_bm1e_provisioning", // x1.small.x86
			fields: fields{Plan: &Plan{Slug: "baremetal_1e"},
				Network: managementIPS,
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer2-bonded",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}}},
			want: NetworkTypeHybrid,
		},
		{
			name: "GetNetworkType_bm1e_provisioning", // x1.small.x86
			fields: fields{Plan: &Plan{Slug: "baremetal_1e"},
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "hybrid",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: false,
					},
				}}},
			want: NetworkTypeHybrid,
		},
		{
			// a configuration that should not be observed, one nic, no bonds.
			name:   "GetNetworkType_OneNic_NoBonds",
			fields: fields{NetworkPorts: []Port{{Bond: nil}}},
			want:   NetworkTypeL2Individual,
		},
		{
			name: "GetNetworkType_TwoNics_NoBonds", // c3-medium l2-individual
			fields: fields{NetworkPorts: []Port{{
				Type: "NetworkBondPort",
				Name: "bond0",
				Data: PortData{
					Bonded: false,
				},
				NetworkType: "layer2-individual",
			}, {
				Type: "NetworkPort",
				Name: "eth0",
				Data: PortData{
					Bonded: false,
				},
			}, {
				Type: "NetworkPort",
				Name: "eth1",
				Data: PortData{
					Bonded: false,
				},
			}}},
			want: NetworkTypeL2Individual,
		},
		{
			name: "GetNetworkType_FourNics_NotBonded", // n2-xlarge l2-individual
			fields: fields{NetworkPorts: []Port{{
				Type: "NetworkBondPort",
				Name: "bond0",
				Data: PortData{
					Bonded: false,
				},
				NetworkType: "layer2-individual",
			}, {
				Type: "NetworkBondPort",
				Name: "bond1",
				Data: PortData{
					Bonded: false,
				},
				NetworkType: "layer2-individual",
			}, {
				Type: "NetworkPort",
				Name: "eth0",
				Data: PortData{
					Bonded: false,
				},
			}, {
				Type: "NetworkPort",
				Name: "eth1",
				Data: PortData{
					Bonded: false,
				},
			}, {
				Type: "NetworkPort",
				Name: "eth2",
				Data: PortData{
					Bonded: false,
				},
			}, {
				Type: "NetworkPort",
				Name: "eth3",
				Data: PortData{
					Bonded: false,
				},
			}}},
			want: NetworkTypeL2Individual,
		},
		{
			name: "GetNetworkType_TwoNics_Bonded_WithManagement", // c3-medium l3
			fields: fields{
				Network: managementIPS,
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer3",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}}},
			want: NetworkTypeL3,
		},
		{
			name: "GetNetworkType_FourNics_Bonded_WithManagement", // n2-xlarge l3
			fields: fields{
				Network: managementIPS,
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer3",
				}, {
					Type: "NetworkBondPort",
					Name: "bond1",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer3",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth2",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth3",
					Data: PortData{
						Bonded: true,
					},
				},
				}},
			want: NetworkTypeL3,
		},

		{
			name: "GetNetworkType_FourNics_Bonded_NoIPs", // n2-xlarge l2-bonded
			fields: fields{
				Network: []*IPAddressAssignment{},
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer2-bonded",
				}, {
					Type: "NetworkBondPort",
					Name: "bond1",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer2-bonded",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth2",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth3",
					Data: PortData{
						Bonded: true,
					},
				},
				}},
			want: NetworkTypeL2Bonded,
		},
		{
			name: "GetNetworkType_TwoNics_Bonded", // c3-medium l2-bonded
			fields: fields{NetworkPorts: []Port{{
				Type: "NetworkBondPort",
				Name: "bond0",
				Data: PortData{
					Bonded: true,
				},
				NetworkType: "layer2-bonded",
			}, {
				Type: "NetworkPort",
				Name: "eth0",
				Data: PortData{
					Bonded: true,
				},
			}, {
				Type: "NetworkPort",
				Name: "eth1",
				Data: PortData{
					Bonded: true,
				},
			}}},
			want: NetworkTypeL2Bonded,
		},
		{
			name: "GetNetworkType_TwoNics_OneBonded", // c3-medium hybrid
			fields: fields{NetworkPorts: []Port{{
				Type: "NetworkBondPort",
				Name: "bond0",
				Data: PortData{
					Bonded: true,
				},
				NetworkType: "hybrid",
			}, {
				Type: "NetworkPort",
				Name: "eth0",
				Data: PortData{
					Bonded: true,
				},
			}, {
				Type: "NetworkPort",
				Name: "eth1",
				Data: PortData{
					Bonded: false,
				},
			}}},
			want: NetworkTypeHybrid,
		},
		{
			name: "GetNetworkType_FourNics_Bonded", // n2-xlarge hybrid
			fields: fields{
				Network: managementIPS,
				NetworkPorts: []Port{{
					Type: "NetworkBondPort",
					Name: "bond0",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer3", // the API reports layer3
				}, {
					Type: "NetworkBondPort",
					Name: "bond1",
					Data: PortData{
						Bonded: true,
					},
					NetworkType: "layer2-individual",
				}, {
					Type: "NetworkPort",
					Name: "eth0",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth1",
					Data: PortData{
						Bonded: false,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth2",
					Data: PortData{
						Bonded: true,
					},
				}, {
					Type: "NetworkPort",
					Name: "eth3",
					Data: PortData{
						Bonded: false,
					},
				}}},
			want: NetworkTypeHybrid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Device{
				Network:      tt.fields.Network,
				Plan:         tt.fields.Plan,
				NetworkPorts: tt.fields.NetworkPorts,
			}
			got := d.GetNetworkType()

			if got != tt.want {
				t.Errorf("Device.GetNetworkType() = %v, want %v", got, tt.want)
			}
		})
	}
}
