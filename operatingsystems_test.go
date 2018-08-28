package packngo

import (
	"testing"
)

func TestAccOS(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	_, _, err := c.OperatingSystems.List()

	if err != nil {
		t.Fatal(err)
	}
}
