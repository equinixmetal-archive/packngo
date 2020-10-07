package packngo

import (
	"testing"
)

func TestAccVPN(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	if u.TwoFactor == "" {
		t.Fatal("VPN can't be used with with disabled 2FA")
	}

	if u.VPN {
		t.Fatal("You must disable VPN in your account before this test")
	}

	_, err = c.VPN.Enable()
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.VPN.Get(testFacility(), nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.VPN.Disable()
	if err != nil {
		t.Fatal(err)
	}
}
