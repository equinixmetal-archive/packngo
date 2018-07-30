package packngo

import (
	"fmt"
	"testing"
)

func TestAccVPN(t *testing.T) {

	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)

	_, err := c.VPN.Enable()
	if err != nil {
		t.Fatal(err)
	}

	config, _, err := c.VPN.Get("ewr1")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(config.Config)

	_, err = c.VPN.Disable()
	if err != nil {
		t.Fatal(err)
	}
}
