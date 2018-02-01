package packngo

import (
	"testing"
)

func TestUserCurrent(t *testing.T) {
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
