package packngo

import (
	"testing"
)

func TestAccCheckCapacity(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	input := &CapacityInput{
		[]ServerInfo{
			{
				Facility: "ams1",
				Plan:     "baremetal_0",
				Quantity: 1},
		},
	}

	_, err := c.CapacityService.Check(input)
	if err != nil {
		t.Fatal("Requested check should have passed for:", input.Servers[0].Facility, input.Servers[0].Plan, input.Servers[0].Quantity)
	}

	list, _, err := c.CapacityService.List()
	if err != nil {
		t.Fatal("List of capacities not fetched")
	}

	for k, v := range *list {
		if v["baremetal_2a2"].Level == "unavailable" {
			input.Servers[0].Plan = v["baremetal_2a2"].Level
			input.Servers[0].Facility = k
			break
		}
	}

	_, err = c.CapacityService.Check(input)
	if err == nil {
		t.Fatal("Requested check should have failed for:", input.Servers[0].Facility, input.Servers[0].Plan, input.Servers[0].Quantity)
	}
}
