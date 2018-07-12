package packngo

import (
	"fmt"
	"testing"
)

var batchID string

func TestAccCreateBatch(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

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

	batches, _, err := c.Batches.Create("93125c2a-8b78-4d4f-a3c4-7367d6b7cca8", req)
	if err != nil {
		t.Fatal(err)
	}

	if len(batches) != 0 {
		batchID = batches[0].ID
	}
}
func TestAccListBatches(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)
	batches, _, err := c.Batches.List("93125c2a-8b78-4d4f-a3c4-7367d6b7cca8", nil)

	if err != nil {
		t.Fatal(err)
	}

	if batches == nil {
		t.Fatal("No batches have been created")
	}

	fmt.Println(len(batches))
	projects, _, err := c.Projects.List(nil)
	for _, p := range projects {
		fmt.Println(p.PaymentMethod)
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
