package packngo

import (
	"testing"
)

func TestAccUserCurrent(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

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

	c, stopRecord := setup(t)
	defer stopRecord()

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

/*
func strPtr(s string) *string {
	return &s
}

func TestAccUserUpdate(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	if u.DefaultOrganizationID == "" {
		t.Fatal("Expected DefaultOrganizationID should not be empty")
	}
	testFirst := "test firstname"
	testLast := "test lasstname"

	uur := UserUpdateRequest{
		FirstName: strPtr(testFirst),
		LastName:  strPtr(testLast),
	}
	uu, _, err := c.Users.Update(&uur)
	if err != nil {
		t.Fatal(err)
	}

	if uu.LastName != testLast {
		t.Fatalf("Updated last name should be %s, was %s", testLast, uu.LastName)
	}
	if uu.FirstName != testFirst {
		t.Fatalf("Updated first name should be %s, was %s", testFirst, uu.FirstName)
	}

	orig := UserUpdateRequest{
		FirstName: strPtr(u.FirstName),
		LastName:  strPtr(u.LastName),
	}
	_, _, err = c.Users.Update(&orig)
	if err != nil {
		t.Fatal(err)
	}
}
*/
