package packngo

import (
	"testing"
)

func TestAccListHardwareReservations(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	projectID := "93125c2a-8b78-4d4f-a3c4-7367d6b7cca8"
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
