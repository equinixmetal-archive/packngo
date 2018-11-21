package packngo

import (
	"testing"
)

func TestAccUserCurrent(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	if u.DefaultOrganizationID == "" {
		t.Fatal("Expected DefaultOrganizationID should not be empty")
	}
}

func TestAccUsersList(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)

	us, _, err := c.Users.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(us) == 0 {
		t.Fatal("At least the current user shall be listed")
	}
	u, _, err := c.Users.Get(us[0].ID, &GetOptions{Includes: []string{"emails"}})

	if err != nil {
		t.Fatal(err)
	}
	if u.Emails[0].Address == "" {
		t.Fatal("User email should have been included.")
	}

}
