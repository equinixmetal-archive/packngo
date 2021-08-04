package packngo

import (
	"path"
	"reflect"
	"testing"
)

func TestAccPublicIPReservation(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	quantityToMask := map[int]int{
		1: 32, 2: 31, 4: 30, 8: 29, 16: 28,
	}

	testFac := testFacility()
	quantity := 2

	ipList, _, err := c.ProjectIPs.List(projectID, &ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(ipList) != 0 {
		t.Fatalf("There should be no reservations a new project, existing list: %s", ipList)
	}

	customData := map[string]interface{}{"custom1": "data", "custom2": map[string]interface{}{"nested": "data"}}
	tags := []string{"Tag1", "Tag2"}

	req := IPReservationRequest{
		Type:                   PublicIPv4,
		Quantity:               quantity,
		Facility:               &testFac,
		CustomData:             &customData,
		Tags:                   tags,
		FailOnApprovalRequired: true,
	}

	res, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}

	if res.CIDR != quantityToMask[quantity] {
		t.Fatalf(
			"CIDR prefix length for requested reservation should be %d, was %d",
			quantityToMask[quantity], res.CIDR)
	}

	if path.Base(res.Project.GetHref()) != projectID {
		t.Fatalf("Wrong project linked in reserved block: %s", res.Project.Href)
	}

	if res.Management {
		t.Fatal("Management flag of new reservation block must be False")
	}

	if res.Facility.Code != testFac {
		t.Fatalf(
			"Facility of new reservation should be %s, was %s", testFac,
			res.Facility.Code)
	}

	if !reflect.DeepEqual(customData, res.CustomData) {
		t.Fatalf("CustomData of new reservation should be %+v, was %+v", customData, res.CustomData)
	}

	if !reflect.DeepEqual(tags, res.Tags) {
		t.Fatalf(
			"Tags of new reservation should be %+v, was %+v", tags, res.Tags)
	}

	ipList, _, err = c.ProjectIPs.List(projectID, &ListOptions{})
	if len(ipList) != 1 {
		t.Fatalf("There should be only one reservation, was: %s", ipList)
	}
	if err != nil {
		t.Fatal(err)
	}

	globalPtr := ipList[0].Global
	if globalPtr != false {
		t.Fatalf("The reserved IP should not be global")
	}

	sameRes, _, err := c.ProjectIPs.Get(res.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sameRes.ID != res.ID {
		t.Fatalf("re-requested test reservation should be %s, is %s",
			res, sameRes)
	}

	availableAddresses, _, err := c.ProjectIPs.AvailableAddresses(
		res.ID, &AvailableRequest{CIDR: 32})
	if err != nil {
		t.Fatal(err)
	}
	if len(availableAddresses) != quantity {
		t.Fatalf("New block should have %d available addresses, got %s",
			quantity, availableAddresses)
	}

	deleteProjectIP(t, c, res.ID)

	_, _, err = c.ProjectIPs.Get(res.ID, nil)
	if err == nil {
		t.Fatalf("Reservation %s should be deleted at this point", res)
	}
}

func TestAccGlobalIPReservation(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	quantityToMask := map[int]int{
		1: 32, 2: 31, 4: 30, 8: 29, 16: 28,
	}

	quantity := 1

	ipList, _, err := c.ProjectIPs.List(projectID, &ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(ipList) != 0 {
		t.Fatalf("There should be no reservations a new project, existing list: %s", ipList)
	}

	description := "packngo test"
	req := IPReservationRequest{
		Type:        GlobalIPv4,
		Quantity:    quantity,
		Description: description,
	}

	res, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}

	if (res.Description == nil) || (*(res.Description) != description) {
		t.Fatalf("Description should be %s, was %v", description, res.Description)
	}

	if res.CIDR != quantityToMask[quantity] {
		t.Fatalf(
			"CIDR prefix length for requested reservation should be %d, was %d",
			quantityToMask[quantity], res.CIDR)
	}

	if path.Base(res.Project.GetHref()) != projectID {
		t.Fatalf("Wrong project linked in reserved block: %s", res.Project.Href)
	}

	if res.Management {
		t.Fatal("Management flag of new reservation block must be False")
	}

	if res.Facility != nil {
		t.Fatalf("Facility of new reservation should be nil")
	}

	ipList, _, err = c.ProjectIPs.List(projectID, &ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(ipList) != 1 {
		t.Fatalf("There should be only one reservation, was: %s", ipList)
	}

	if !ipList[0].Global {
		t.Fatalf("The reserved IP should be global")
	}

	sameRes, _, err := c.ProjectIPs.Get(res.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sameRes.ID != res.ID {
		t.Fatalf("re-requested test reservation should be %s, is %s",
			res, sameRes)
	}

	availableAddresses, _, err := c.ProjectIPs.AvailableAddresses(res.ID, &AvailableRequest{CIDR: 32})
	if err != nil {
		t.Fatal(err)
	}
	if len(availableAddresses) != quantity {
		t.Fatalf("New block should have %d available addresses, got %s",
			quantity, availableAddresses)
	}

	_, err = c.ProjectIPs.Remove(res.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.ProjectIPs.Get(res.ID, nil)
	if err == nil {
		t.Fatalf("Reservation %s should be deleted at this point", res)
	}
}

func TestAccPublicIPReservationFailFast(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	testFac := testFacility()
	// this should be an absurdly high number
	quantity := 256

	customData := map[string]interface{}{"custom1": "data", "custom2": map[string]interface{}{"nested": "data"}}

	req := IPReservationRequest{
		Type:                   PublicIPv4,
		Quantity:               quantity,
		Facility:               &testFac,
		CustomData:             &customData,
		FailOnApprovalRequired: true,
	}

	_, resp, err := c.ProjectIPs.Request(projectID, &req)
	if err == nil {
		t.Fatal("should have had an error 422")
	}
	if resp == nil {
		t.Fatal("unexpected response was nil")
	}
	if resp.StatusCode != 422 {
		t.Fatalf("received response code %d instead of expected %d", resp.StatusCode, 422)
	}
}

func TestAccPublicMetroIPReservation(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()
	quantityToMask := map[int]int{
		1: 32, 2: 31, 4: 30, 8: 29, 16: 28,
	}

	testMetro := testMetro()
	quantity := 2

	ipList, _, err := c.ProjectIPs.List(projectID, &ListOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(ipList) != 0 {
		t.Fatalf("There should be no reservations a new project, existing list: %s", ipList)
	}

	customData := map[string]interface{}{"custom1": "data", "custom2": map[string]interface{}{"nested": "data"}}
	tags := []string{"Tag1", "Tag2"}

	req := IPReservationRequest{
		Type:                   PublicIPv4,
		Quantity:               quantity,
		Metro:                  &testMetro,
		CustomData:             &customData,
		Tags:                   tags,
		FailOnApprovalRequired: true,
	}

	res, _, err := c.ProjectIPs.Request(projectID, &req)
	if err != nil {
		t.Fatal(err)
	}

	if res.CIDR != quantityToMask[quantity] {
		t.Fatalf(
			"CIDR prefix length for requested reservation should be %d, was %d",
			quantityToMask[quantity], res.CIDR)
	}

	if path.Base(res.Project.GetHref()) != projectID {
		t.Fatalf("Wrong project linked in reserved block: %s", res.Project.Href)
	}

	if res.Management {
		t.Fatal("Management flag of new reservation block must be False")
	}

	if res.Metro.Code != testMetro {
		t.Fatalf(
			"Metro of new reservation should be %s, was %s", testMetro,
			res.Facility.Code)
	}

	if !reflect.DeepEqual(customData, res.CustomData) {
		t.Fatalf("CustomData of new reservation should be %+v, was %+v", customData, res.CustomData)
	}

	if !reflect.DeepEqual(tags, res.Tags) {
		t.Fatalf(
			"Tags of new reservation should be %+v, was %+v", tags, res.Tags)
	}

	ipList, _, err = c.ProjectIPs.List(projectID, &ListOptions{})
	if len(ipList) != 1 {
		t.Fatalf("There should be only one reservation, was: %s", ipList)
	}
	if err != nil {
		t.Fatal(err)
	}

	if ipList[0].Global {
		t.Fatalf("The reserved IP should not be global")
	}

	sameRes, _, err := c.ProjectIPs.Get(res.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if sameRes.ID != res.ID {
		t.Fatalf("re-requested test reservation should be %s, is %s",
			res, sameRes)
	}

	availableAddresses, _, err := c.ProjectIPs.AvailableAddresses(
		res.ID, &AvailableRequest{CIDR: 32})
	if err != nil {
		t.Fatal(err)
	}
	if len(availableAddresses) != quantity {
		t.Fatalf("New block should have %d available addresses, got %s",
			quantity, availableAddresses)
	}

	deleteProjectIP(t, c, res.ID)

	_, _, err = c.ProjectIPs.Get(res.ID, nil)
	if err == nil {
		t.Fatalf("Reservation %s should be deleted at this point", res)
	}
}
