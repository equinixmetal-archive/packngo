package packngo

import "testing"

func TestValidateUUID(t *testing.T) {
	if ValidateUUID("") == nil {
		t.Error("UUID should be invalid")
	}
	if ValidateUUID("a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0") != nil {
		t.Error("UUID should be valid")
	}
	if ValidateUUID("a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0") == nil {
		t.Error("UUID should be invalid")
	}
	if ValidateUUID("a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0g0a0a") == nil {
		t.Error("UUID should be invalid")
	}
	if ValidateUUID("a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a") == nil {
		t.Error("UUID should be invalid")
	}
	if ValidateUUID("a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0a") == nil {
		t.Error("UUID should be invalid")
	}
	if ValidateUUID("a0a0a0a0-a0a0-a0a0-a0a0-a0a0a0a0a0a0aa") == nil {
		t.Error("UUID should be invalid")
	}
}
