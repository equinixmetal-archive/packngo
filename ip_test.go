package packngo

import (
	"strings"
	"testing"
)

func TestAccIPReservation(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	quantityToMask := map[int]string{
		1: "32", 2: "31", 4: "30", 8: "29", 16: "28",
	}

	testFac := "ewr1"
	quantity := 2

	req := IPReservationRequest{
		Type:     "public_ipv4",
		Quantity: quantity,
		Comments: "packngo test",
		Facility: testFac,
	}

	af, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}
	addrMask := strings.Split(af.Address, "/")
	if addrMask[1] != quantityToMask[quantity] {
		t.Errorf(
			"CIDR prefix length for requested reservation should be %s, was %s",
			quantityToMask[quantity], addrMask[1])
	}

	res, _, err := c.ProjectIPs.GetByCIDR(projectID, af.Address)
	if err != nil {
		t.Fatal(err)
	}
	if res.Facility.Code != testFac {
		t.Errorf(
			"Facility of new reservation should be %s, was %s", testFac,
			res.Facility.Code)
	}

	ipList, _, err := c.ProjectIPs.List(projectID)
	if len(ipList) != 1 {
		t.Errorf("There should be only one reservation, was: %s", ipList)
	}
	if err != nil {
		t.Fatal(err)
	}

	sameRes, _, err := c.ProjectIPs.Get(res.ID)
	if err != nil {
		t.Fatal(err)
	}
	if sameRes.ID != res.ID {
		t.Errorf("re-requested test reservation should be %s, is %s",
			res, sameRes)
	}

	availableAddresses, _, err := c.ProjectIPs.AvailableAddresses(
		res.ID, &AvailableRequest{CIDR: 32})
	if err != nil {
		t.Fatal(err)
	}
	if len(availableAddresses) != quantity {
		t.Errorf("New block should have %d available addresses, got %s",
			quantity, availableAddresses)
	}

	_, err = c.ProjectIPs.Remove(res.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.ProjectIPs.Get(res.ID)
	if err == nil {
		t.Errorf("Reservation %s should be deleted at this point", res)
	}
}
