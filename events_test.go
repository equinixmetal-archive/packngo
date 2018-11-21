package packngo

import (
	"testing"
)

func TestAccListEvents(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	events, _, err := c.Events.List(&ListOptions{Page: 1, PerPage: 9})
	if err != nil {
		t.Fatal(err)
	}

	if len(events) == 0 {
		t.Fatal("Events are empty")
	}

	event, _, err := c.Events.Get(events[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if event == nil {
		t.Fatal("Event not returned")
	}
}
