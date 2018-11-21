package packngo

import (
	"testing"
	"time"
)

func TestAccInstanceBatches(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	req := &BatchCreateRequest{
		Batches: []BatchCreateDevice{
			{
				DeviceCreateRequest: DeviceCreateRequest{
					Hostname:     "test1",
					Plan:         "baremetal_0",
					OS:           "ubuntu_16_04",
					Facility:     []string{"ewr1"},
					BillingCycle: "hourly",
					Tags:         []string{"abc"},
				},
				Quantity: 3,
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
	batch, _, err := c.Batches.Get(batchID, &GetOptions{Includes: []string{"devices"}})

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
