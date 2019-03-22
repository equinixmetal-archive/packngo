package packngo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func waitConnectStatus(connectID, projectID, status string, c *Client) (*Connect, error) {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		co, _, err := c.Connects.Get(connectID, projectID, nil)
		if err != nil {
			return nil, err
		}
		if co.Status == status {
			return co, nil
		}
	}
	return nil, fmt.Errorf("Packet Connect %s is still noti %s after timeout", connectID, status)
}

func TestAccConnectBasic(t *testing.T) {

	azureEnvVar := "AZURE_KEY"
	azureKey := os.Getenv(azureEnvVar)
	if len(azureKey) == 0 {
		t.Fatalf("You must set %s", azureEnvVar)
	}

	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	cs, _, err := c.Connects.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(cs)
	if len(cs) != 0 {
		t.Fatalf("There should be no Connect resource")
	}

	vcr := VirtualNetworkCreateRequest{
		ProjectID:   projectID,
		Description: "connectTestVLAN",
		Facility:    "ewr1",
	}
	vlan, _, err := c.ProjectVirtualNetworks.Create(&vcr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.ProjectVirtualNetworks.Delete(vlan.ID)

	ccr := ConnectCreateRequest{
		Name:            "testconn",
		ProjectID:       projectID,
		ProviderID:      AzureProviderID,
		ProviderPayload: azureKey,
		Facility:        "ewr1",
		PortSpeed:       100,
		VLAN:            vlan.VXLAN,
		Description:     "testconn",
		Tags:            []string{"testconn"},
	}

	connect, _, err := c.Connects.Create(&ccr)
	if err != nil {
		t.Fatal(err)
	}
	connect, err = waitConnectStatus(connect.ID, projectID, "PROVISIONED", c)
	if err != nil {
		t.Fatal(err)
	}

	defer c.Connects.Delete(connect.ID, projectID)

	cs, _, err = c.Connects.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	log.Println(cs)
	if len(cs) != 1 {
		t.Fatalf("There should be only 1 Connect resource")
	}

	time.Sleep(30 * time.Second)

	connect, _, err = c.Connects.Deprovision(connect.ID, projectID)
	if err != nil {
		t.Fatal(err)
	}
	connect, err = waitConnectStatus(connect.ID, projectID, "DEPROVISIONED", c)
	if err != nil {
		t.Fatal(err)
	}

}
