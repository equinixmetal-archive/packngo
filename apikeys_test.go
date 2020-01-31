package packngo

import (
	"testing"
)

func createNTestAPIKeys(n int, projectID string, c *Client, t *testing.T) (idMap map[string]struct{}) {
	idMap = make(map[string]struct{})
	for i := 0; i < n; i++ {
		req := APIKeyCreateRequest{
			Description: "PACKNGO_TEST_KEY_DELETE_ME-" + randString8(),
			ReadOnly:    true,
			ProjectID:   projectID,
		}
		key, _, err := c.APIKeys.Create(&req)
		if err != nil {
			t.Fatalf("errored posting key: %v", err)
		}
		idMap[key.ID] = struct{}{}
	}
	return idMap
}

func TestAccAPIKeyListProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	nKeys := 10

	keyIDs := createNTestAPIKeys(nKeys, projectID, c, t)

	if len(keyIDs) != nKeys {
		t.Fatalf("Helper function was supposed to create %d keys, created %d", nKeys, len(keyIDs))
	}

	keyList, _, err := c.APIKeys.ProjectList(projectID, nil)
	if err != nil {
		t.Fatalf("Error getting list of Project keys %s", err)
	}

	if len(keyList) != nKeys {
		t.Fatalf("Listing should return %d keys, returned %d", nKeys, len(keyList))
	}
	for kID := range keyIDs {
		_, err := c.APIKeys.Delete(kID)

		if err != nil {
			t.Logf("Failed to delete testing key %s", kID)
		}
		delete(keyIDs, kID)
	}
	if len(keyIDs) != 0 {
		t.Fatal("Error when sweeping keys after ProjectList test - Seems that not all test keys were deleted")
	}

}

func TestAccAPIKeyListUser(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c := setup(t)

	nKeys := 10

	keyIDs := createNTestAPIKeys(nKeys, "", c, t)

	if len(keyIDs) != nKeys {
		t.Fatalf("Helper function was supposed to create %d keys, created %d", nKeys, len(keyIDs))
	}

	keyList, _, err := c.APIKeys.UserList(nil)
	if err != nil {
		t.Fatalf("Error getting list of User keys %s", err)
	}

	if len(keyList) < nKeys {
		t.Fatalf("Listing should return at least %d keys, returned %d", nKeys, len(keyList))
	}
	for kID := range keyIDs {
		_, err := c.APIKeys.Delete(kID)

		if err != nil {
			t.Logf("Failed to delete testing key %s", kID)
		}
		delete(keyIDs, kID)
	}
	if len(keyIDs) != 0 {
		t.Fatal("Error when sweeping keys after UserList test - Seems that not all test keys were deleted")
	}

}

func TestAccAPIKeyCreateUser(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c := setup(t)
	req := APIKeyCreateRequest{
		Description: "PACKNGO_TEST_KEY_DELETE_ME-" + randString8(),
		ReadOnly:    true,
	}

	key, _, err := c.APIKeys.Create(&req)
	if err != nil {
		t.Fatalf("errored posting key: %v", err)
	}
	if len(key.User.URL) == 0 {
		t.Error("new Key doesn't have User URL set")
	}
	if len(key.Token) == 0 {
		t.Error("new Key doesn't have token set")
	}

	if key.Description != req.Description {
		t.Fatalf("returned key label does not match, want: %v, got: %v", req.Description, key.Description)
	}

	gotKey, err := c.APIKeys.UserGet(key.ID, nil)
	if err != nil {
		t.Fatalf("Error getting created User API key: %s", err)
	}
	if gotKey.Token != key.Token {
		t.Fatalf("Strange mismatch in tokens of the same test key")
	}
	_, err = c.APIKeys.Delete(key.ID)
	if err != nil {
		t.Fatalf("error deleting key")
	}
}

func TestAccAPIKeyCreateProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	req := APIKeyCreateRequest{
		Description: "PACKNGO_TEST_KEY_DELETE_ME-" + randString8(),
		ReadOnly:    true,
		ProjectID:   projectID,
	}

	key, _, err := c.APIKeys.Create(&req)
	if err != nil {
		t.Fatalf("errored posting key: %v", err)
	}
	if len(key.Project.URL) == 0 {
		t.Error("new Key doesn't have Project URL set")
	}
	if len(key.Token) == 0 {
		t.Error("new Key doesn't have token set")
	}

	if key.Description != req.Description {
		t.Fatalf("returned key label does not match, want: %v, got: %v", req.Description, key.Description)
	}

	gotKey, err := c.APIKeys.ProjectGet(projectID, key.ID, nil)
	if err != nil {
		t.Fatalf("Error getting created Project API key: %s", err)
	}
	if gotKey.Token != key.Token {
		t.Fatalf("Strange mismatch in tokens of the same test key")
	}
	_, err = c.APIKeys.Delete(key.ID)
	if err != nil {
		t.Fatalf("error deleting key")
	}
}
