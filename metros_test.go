package packngo

import (
	"fmt"
	"testing"
)

func TestAccMetrosList(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

	l, _, err := c.Metros.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) == 0 {
		t.Fatal(fmt.Errorf("Expected to get at least one metro"))

	}
	for _, m := range l {
		if m.Code == "" {
			t.Fatal(fmt.Errorf("metro %+v has empty Code (slug) attr", m))
		}
	}
}
