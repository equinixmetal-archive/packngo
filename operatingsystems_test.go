package packngo

import (
	"testing"
)

func TestAccOS(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()
	l, _, err := c.OperatingSystems.List()

	if len(l) == 0 {
		t.Fatal("Empty Operating System listing from the API")
	}

	if err != nil {
		t.Fatal(err)
	}
}
