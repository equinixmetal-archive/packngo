package packngo

import (
	"testing"
	"time"
)

func TestAccInstanceBatches(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	req := &InstanceBatchCreateRequest{
		Batches: []BatchInstance{
			{
				Hostname:        "test1",
				Description:     "test batch",
				Plan:            "baremetal_0",
				OperatingSystem: "ubuntu_16_04",
				Facility:        "ewr1",
				BillingCycle:    "hourly",
				Tags:            []string{"abc"},
				Quantity:        3,
			},
		},
	}

	batches, _, err := c.Batches.Create(projectID, req)
	if err != nil {
		t.Fatal(err)
	}

	var batchID string
	if len(batches) != 0 {
		batchID = batches[0].ID
	}

	batches, _, err = c.Batches.List(projectID, nil)

	if err != nil {
		t.Fatal(err)
	}

	if batches == nil {
		t.Fatal("No batches have been created")
	}
	time.Sleep(5 * time.Second)
	batch, _, err := c.Batches.Get(batchID, &ListOptions{Includes: "devices"})

	if err != nil {
		t.Fatal(err)
	}

	if batch == nil {
		t.Fatal("Batch not found")
	}

	for _, d := range batch.Devices {
		waitDeviceActive(d.ID, c)
	}

	_, err = c.Batches.Delete(batchID, true)

	if err != nil {
		t.Fatal(err)
	}
}
