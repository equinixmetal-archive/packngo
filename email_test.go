package packngo

import (
	"testing"
)

func TestAccCreateEmail(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	updatedAddress := "update@domain.com"
	c := setup(t)

	req := &EmailRequest{
		Address: "test@domain.com",
	}

	ret, _, err := c.Emails.Create(req)
	if err != nil {
		t.Fatal("Create failed", err)
	}

	if req.Address != ret.Address {
		t.Fatal("Address not equal to", req.Address)
	}

	emailID := ret.ID
	req.Address = updatedAddress

	email, _, err := c.Emails.Get(emailID, nil)
	if err != nil {
		t.Fatal("Get failed", err)
	}

	ret, _, err = c.Emails.Update(email.ID, req)
	if err != nil {
		t.Fatal("Update failed:", err)
	}

	if req.Address == ret.Address {
		t.Fatal("Address not updated")
	}

	_, err = c.Emails.Delete(email.ID)
	if err != nil {
		t.Fatal("Delete failed:", err)
	}
}
