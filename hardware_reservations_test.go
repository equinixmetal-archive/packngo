package packngo

import (
	"testing"
)

func TestAccListHardwareReservations(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, tearDown := setupWithProject(t)
	defer tearDown()

	hardwareReservations, _, err := c.HardwareReservations.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(hardwareReservations) != 0 {
		t.Fatal("There should not be any hardware reservations.")
	}
}
