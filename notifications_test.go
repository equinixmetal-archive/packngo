package packngo

import (
	"testing"
)

func TestAccListNotifications(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	notifications, _, err := c.Notifications.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(notifications) == 0 {
		t.Fatal("Notifications are empty")
	}
}
