package packngo

import (
	"testing"
)

var batchID string

func TestAccCreateBatch(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	batches := &InstanceBatchCreateRequest{
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

	batch, _, err := c.Batches.Create("6e505da6-c861-487f-adcb-adfe3959da49", batches)
	if err != nil {
		t.Fatal(err)
	}

	if batch != nil {
		batchID = batch.ID
	}
}
func TestAccListBatches(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)
	batches, _, err := c.Batches.List("", nil)

	if err != nil {
		t.Fatal(err)
	}

	if batches == nil {
		t.Fatal("No batches have been created")
	}
}

func TestAccGetBatch(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	batch, _, err := c.Batches.Get(batchID, nil)

	if err != nil {
		t.Fatal(err)
	}

	if batch == nil {
		t.Fatal("Batch not found")
	}
}
