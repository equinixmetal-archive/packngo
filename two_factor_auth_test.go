package packngo

import "testing"

func TestAccTFAApp(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	_, err := c.TwoFactorAuth.EnableApp()
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.TwoFactorAuth.DisableApp()
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccTFASms(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	_, err := c.TwoFactorAuth.EnableSms()
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.TwoFactorAuth.DisableSms()
	if err != nil {
		t.Fatal(err)
	}
}
