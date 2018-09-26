package packngo

import (
	"testing"
	"time"
)

func TestAccInstanceBatches(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	req := &InstanceBatchCreateRequest{
		Batches: []BatchInstance{
			{
				Hostname:        "test1",
				Description:     "test batch",
				Plan:            "baremetal_0",
				OperatingSystem: "coreos_stable",
				Facility:        "ewr1",
				BillingCycle:    "hourly",
				Tags:            []string{"abc"},
				Quantity:        1,
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

	batch, _, err := c.Batches.Get(batchID, &ListOptions{Includes: "devices"})

	if err != nil {
		t.Fatal(err)
	}

	var finished bool

	// Wait for all devices to become 'active'
	for {
		if len(batch.Devices) == 0 {
			break
		}
		for _, d := range batch.Devices {

			dev, _, _ := c.Devices.Get(d.ID)
			if dev.State == "active" {
				finished = true
			} else { //if at least one is not "active" set finished to false and break the loop
				finished = false
				break
			}
		}

		if finished {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	if batch == nil {
		t.Fatal("Batch not found")
	}

	_, err = c.Batches.Delete(batchID, true)

	if err != nil {
		t.Fatal(err)
	}
}
