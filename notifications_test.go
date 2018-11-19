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

	notification, _, err := c.Notifications.Get(notifications[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if notification == nil {
		t.Fatal("Notification not returned")
	}

	readNotification, _, err := c.Notifications.MarkAsRead(notifications[0].ID)
	if err != nil {
		t.Fatal(err)
	}

	if readNotification.Read == false {
		t.Fatal("Notification was not marked as read")
	}

}
