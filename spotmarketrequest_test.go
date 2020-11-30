package packngo

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func deleteSpotMarketRequest(t *testing.T, c *Client, id string, force bool) {
	if _, err := c.SpotMarketRequests.Delete(id, force); err != nil {
		t.Fatal(err)
	}
}

func TestAccSpotMarketRequestBasic(t *testing.T) {
	// This test is only going to create the spot market request with
	// max bid price set to half of current spot price, so that the devices
	// are not run at all.
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()
	fac := testFacility()

	ps := SpotMarketRequestInstanceParameters{
		BillingCycle:    "hourly",
		Plan:            testPlan(),
		OperatingSystem: testOS,
		Hostname:        fmt.Sprintf("%s{{index}}", hn),
	}

	prices, _, err := c.SpotMarket.Prices()
	pri := prices[fac][testPlan()]
	if err != nil {
		t.Fatal(err)
	}

	cr := SpotMarketRequestCreateRequest{
		DevicesMax:  3,
		DevicesMin:  2,
		FacilityIDs: []string{fac, testFacilityAlternate},
		MaxBidPrice: pri / 2,
		Parameters:  ps,
	}

	smr, _, err := c.SpotMarketRequests.Create(&cr, projectID)

	if err != nil {
		t.Fatal(err)
	}
	defer deleteSpotMarketRequest(t, c, smr.ID, true)

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

	if smrs[0].Plan.Slug != testPlan() {
		t.Fatalf("Plan should be reported as %s, was %s", testPlan(), smrs[0].Plan.Slug)
	}

	smr2, _, err := c.SpotMarketRequests.Get(smr.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if (smr.ID != smrs[0].ID) || (smrs[0].ID != smr2.ID) {
		t.Fatal("mismatch in the created SpotMarketRequest")
	}
}

// I am not sure if spot-market-requests work in the new DCs. The test works with baremetal_0 in ewr1 i.e.:
//
// PACKNGO_TEST_PLAN=baremetal_0 PACKNGO_TEST_FACILITY=ewr1 PACKNGO_DEBUG=1 PACKNGO_TEST_ACTUAL_API=1 go test -v -timeout=20m -run=TestAccSpotMarketRequestPriceAware

func TestAccSpotMarketRequestPriceAware(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	prices, _, err := c.SpotMarket.Prices()
	if err != nil {
		t.Fatal(err)
	}

	fac := testFacility()
	pri := prices[fac][testPlan()]
	thr := pri * 1.2
	hn := randString8()

	ps := SpotMarketRequestInstanceParameters{
		BillingCycle:    "hourly",
		Plan:            testPlan(),
		OperatingSystem: testOS,
		Hostname:        fmt.Sprintf("%s{{index}}", hn),
	}

	nDevices := 2

	cr := SpotMarketRequestCreateRequest{
		DevicesMax:  nDevices,
		DevicesMin:  nDevices,
		FacilityIDs: []string{fac},
		MaxBidPrice: thr,
		Parameters:  ps,
	}

	smr, _, err := c.SpotMarketRequests.Create(&cr, projectID)
	if err != nil {
		t.Fatal(err)
	}

out:
	for {
		<-time.Tick(5 * time.Second)
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
	log.Println("all devices active")
	deleteSpotMarketRequest(t, c, smr.ID, true)
out2:
	// wait for devices to disappear .. takes ~5 minutes
	for {
		<-time.Tick(5 * time.Second)
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
