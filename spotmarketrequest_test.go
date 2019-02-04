package packngo

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestAccSpotMarketRequestBasic(t *testing.T) {
	// This test is only going to create the spot market request with
	// max bid price set to half of current spot price, so that the devices
	// are not run at all.
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	ps := SpotMarketRequestInstanceParameters{
		BillingCycle:    "hourly",
		Plan:            "baremetal_0",
		OperatingSystem: "rancher",
		Hostname:        fmt.Sprintf("%s{{index}}", hn),
	}

	prices, _, err := c.SpotMarket.Prices()
	pri := prices["ewr1"]["baremetal_0"]
	if err != nil {
		t.Fatal(err)
	}

	cr := SpotMarketRequestCreateRequest{
		DevicesMax:  3,
		DevicesMin:  2,
		FacilityIDs: []string{"ewr1", "sjc1"},
		MaxBidPrice: pri / 2,
		Parameters:  ps,
	}

	smr, _, err := c.SpotMarketRequests.Create(&cr, projectID)

	if err != nil {
		t.Fatal(err)
	}
	defer c.SpotMarketRequests.Delete(smr.ID, true)

	if smr.Project.ID != projectID {
		t.Fatal("Strange project ID in SpotMarketReuqest:", smr.Project.ID)
	}

	smrs, _, err := c.SpotMarketRequests.List(projectID,
		&ListOptions{Includes: []string{"devices,project,plan"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(smrs) != 1 {
		t.Fatal("there should be only one SpotMarketRequest")
	}

	if smrs[0].Plan.Slug != "t1.small.x86" {
		t.Fatal("Plan should be reported as t1.small.x86 (aka baremetal_0).")
	}

	smr2, _, err := c.SpotMarketRequests.Get(smr.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if (smr.ID != smrs[0].ID) || (smrs[0].ID != smr2.ID) {
		t.Fatal("mismatch in the created SpotMarketRequest")
	}
}

func TestAccSpotMarketRequestPriceAware(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	prices, _, err := c.SpotMarket.Prices()
	if err != nil {
		t.Fatal(err)
	}

	pri := prices["ewr1"]["baremetal_0"]
	thr := pri * 2.0
	hn := randString8()

	ps := SpotMarketRequestInstanceParameters{
		BillingCycle:    "hourly",
		Plan:            "baremetal_0",
		OperatingSystem: "rancher",
		Hostname:        fmt.Sprintf("%s{{index}}", hn),
	}

	nDevices := 3

	cr := SpotMarketRequestCreateRequest{
		DevicesMax:  nDevices,
		DevicesMin:  nDevices,
		FacilityIDs: []string{"ewr1"},
		MaxBidPrice: thr,
		Parameters:  ps,
	}

	smr, _, err := c.SpotMarketRequests.Create(&cr, projectID)
	if err != nil {
		t.Fatal(err)
	}

out:
	for {
		select {
		case <-time.Tick(5 * time.Second):
			smr, _, err = c.SpotMarketRequests.Get(
				smr.ID,
				&GetOptions{Includes: []string{"devices"}},
			)
			if err != nil {
				t.Fatal(err)
			}
			activeDevs := 0
			for _, d := range smr.Devices {
				if d.State == "active" {
					activeDevs++
				}
			}
			if activeDevs == nDevices {
				break out
			}
		}
	}
	log.Println("all devices active")
	c.SpotMarketRequests.Delete(smr.ID, true)
out2:
	// wait for devices to disappear .. takes ~5 minutes
	for {
		select {
		case <-time.Tick(5 * time.Second):
			ds, _, err := c.Devices.List(projectID, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(ds) > 0 {
				log.Println(len(ds), "devices still exist")
			} else {
				break out2
			}

		}
	}
}
