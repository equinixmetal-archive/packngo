package packngo

import (
	"testing"
)

func TestAccHardwareListOne(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, _, teardown := setupWithProject(t)
	defer teardown()

	fac := testFacility()

	reserved := new(bool)
	opts := &ListOptions{
		FacilityCode: fac,
		Page:         1,
		PerPage:      1,
		State:        "provisionable",
		Plan:         "t1.small.x86",
		Includes:     []string{"manufacturer"},
		Reserved:     reserved,
	}

	hw, _, err := c.Hardware.List(opts)
	if err != nil {
		t.Fatal("error in retrieving hardware list")
	}

	if len(hw) == 0 {
		t.Fatal("expected one hardware item, got zero")
	}

}
