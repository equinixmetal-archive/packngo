package packngo

import (
	"log"
	"testing"
	"time"
)

func TestAccSpotMarketRequestBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	ps := InstanceParameters{
		BillingCycle:    "hourly",
		Plan:            "baremetal_0",
		OperatingSystem: "ubuntu_16_04",
		Hostname:        "test{{index}}",
	}

	cr := SpotMarketRequestCreateRequest{
		DevicesMax:  3,
		DevicesMin:  2,
		FacilityIDs: []string{"ewr1", "ewr1"},
		MaxBidPrice: 0.02,
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

	smrs, _, err := c.SpotMarketRequests.List(projectID)
	if err != nil {
		t.Fatal(err)
	}
	if len(smrs) != 1 {
		t.Fatal("there should be only one SpotMarketRequest")
	}

	smr2, _, err := c.SpotMarketRequests.Get(smr.ID)
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

	ps := InstanceParameters{
		BillingCycle:    "hourly",
		Plan:            "baremetal_0",
		OperatingSystem: "ubuntu_16_04",
		Hostname:        "test{{index}}",
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
			smr, _, err = c.SpotMarketRequests.Get(smr.ID)
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
