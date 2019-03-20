package packngo

import "testing"

func TestAccConnectBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)

	defer teardown()
}
