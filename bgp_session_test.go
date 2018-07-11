package packngo

import (
	"testing"
)

var deviceID string
var projectID string
var sessionID string

func TestAccCreateBGPSession(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, _ := setupWithProject(t)
	hn := randString8()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     testFacility(),
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

	bgpSession, _, err := c.BGPSessions.Create(deviceID, CreateBGPSessionRequest{AddressFamily: "ipv4"})
	if err != nil {
		t.Fatal(err)
	}

	sessionID = bgpSession.ID
}

func TestAccListBGPSessionsByDevice(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	sessions, _, err := c.Devices.ListBGPSessions(deviceID, nil)
	if err != nil {
		t.Fatal(err)
	}

	var check *BGPSession
	for _, s := range sessions {
		if s.ID == sessionID {
			check = &s
			break
		}
	}

	if check == nil {
		t.Fatal("BGP Session not returned.")
	}
}

func TestAccListBGPSessionsByProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	sessions, _, err := c.Projects.ListBGPSessions(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	var check *BGPSession
	for _, s := range sessions {
		if s.ID == sessionID {
			check = &s
			break
		}
	}

	if check == nil {
		t.Fatal("BGP Session not returned.")
	}
}

func TestAccGetBgpSession(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	session, _, err := c.BGPSessions.Get(sessionID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if session == nil {
		t.Fatal("Session not retrieved")
	}
}

func TestAccDeleteBgpSession(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	_, err := c.BGPSessions.Delete(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	session, _, err := c.BGPSessions.Get(sessionID, nil)
	if session != nil {
		t.Fatal("Session not deleted")
	}
	if err == nil {
		t.Fatal("Session not deleted")
	}
}

func TestAccCleanup(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	c.Devices.Delete(deviceID)
	projectTeardown(c)
}
