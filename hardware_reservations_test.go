package packngo

import (
	"testing"
)

func TestAccListHardwareReservations(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	projects, _, err := c.Projects.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) == 0 {
		t.Fatal("No projects returned.")
	}

	projectID := projects[0].ID

	hardwareReservations, _, err := c.HardwareReservations.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	hrID := hardwareReservations[0].ID

	hardwareReservation, _, err := c.HardwareReservations.Get(hrID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if hardwareReservation.ID != hrID {
		t.Fatal("Hardware reservation IDs don't match.")
	}

}
