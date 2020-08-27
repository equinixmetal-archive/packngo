package packngo

import (
	"fmt"
	"testing"
)

func TestAccFacilities(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

	l, _, err := c.Facilities.List(
		&ListOptions{
			GetOptions: GetOptions{Includes: []string{"address"}},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(l) == 0 {
		t.Fatal(fmt.Errorf("Expected to get non-zero facilities"))

	}
	for _, f := range l {
		if f.Code == "" {
			t.Fatal(fmt.Errorf("facility %+v has empty Code (slug) attr", f))
		}
	}
}
