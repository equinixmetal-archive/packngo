package packngo

import (
	"testing"
)

func TestAccVPN(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)

	_, err := c.VPN.Enable()
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = c.VPN.Get("ewr1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.VPN.Disable()
	if err != nil {
		t.Fatal(err)
	}
}
