package packngo

import (
	"fmt"
	"reflect"
	"strings"
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

	resp, err := c.CapacityService.Check(input)
	if err != nil {
		t.Fatal("Requested check should have passed for:", input.Servers[0].Facility, input.Servers[0].Plan, input.Servers[0].Quantity)
	}

	list, _, err := c.CapacityService.List()
	if err != nil {
		t.Fatal("List of capacities not fetched")
	}

	// Getting facility where a plan is unavailable
	s := reflect.ValueOf(list).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i).Interface().(*CapacityPerFacility)
		if f.Baremetal2a2 != nil && f.Baremetal2a2.Level == "unavailable" {
			input.Servers[0].Plan = f.Baremetal2a2.Level
			input.Servers[0].Facility = strings.ToLower(typeOfT.Field(i).Name)
		}
	}

	resp, err = c.CapacityService.Check(input)
	if err == nil {
		t.Fatal("Requested check should have failed for:", input.Servers[0].Facility, input.Servers[0].Plan, input.Servers[0].Quantity)
	}
	fmt.Println(resp)
}
