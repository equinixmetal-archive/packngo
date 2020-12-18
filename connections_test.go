package packngo

import (
	"log"
	"testing"
)

func TestAccConnectionProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	connReq := ConnectionCreateRequest{
		Name:       "testconn",
		Redundancy: "test",
		Facility:   "ewr1",
		Type:       "testtype",
	}

	log.Println("hear")

	conn, _, err := c.Connections.ProjectCreate(projectID, &connReq)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", conn)
}

func TestAccConnectionOrganization(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, stopRecord := setup(t)
	defer stopRecord()

	connReq := ConnectionCreateRequest{
		Name:       "testconn",
		Redundancy: "test",
		Facility:   "ewr1",
		Type:       "testtype",
	}

	user, _, err := c.Users.Current()

	if err != nil {
		t.Fatal(err)
	}

	conn, _, err := c.Connections.OrganizationCreate(user.DefaultOrganizationID, &connReq)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%#v", conn)
}
