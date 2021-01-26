package packngo

import (
	"log"
	"testing"
)

func plansInFacilities(t *testing.T, plans []Plan) {
	if len(plans) == 0 {
		t.Fatal("Empty plans listing from the API")
	}
	avail := map[string][]string{}
	for _, p := range plans {
		for _, f := range p.AvailableIn {
			if _, ok := avail[f.Code]; !ok {
				avail[f.Code] = []string{p.Slug}
			} else {
				avail[f.Code] = append(avail[f.Code], p.Slug)
			}
		}
		if p.Pricing.Hour < 0.0 {
			t.Fatalf("Strange pricing for %s %s", p.Name, p.Slug)
		}
	}

	for f, ps := range avail {
		if len(ps) == 0 {
			t.Fatalf("no plans available in facility %s", f)
		}
		// prints plans available in facility
		log.Printf("%s: %+v\n", f, ps)
	}
}

func TestAccPlansBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()
	l, _, err := c.Plans.List(&ListOptions{Includes: []string{"available_in"}})
	if err != nil {
		t.Fatal(err)
	}
	plansInFacilities(t, l)
}

func TestAccPlansProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	l, _, err := c.Plans.ProjectList(projectID, &ListOptions{Includes: []string{"available_in"}})
	if err != nil {
		t.Fatal(err)
	}
	plansInFacilities(t, l)
}

func TestAccPlansOrganization(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()

	user, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	l, _, err := c.Plans.OrganizationList(user.DefaultOrganizationID, &ListOptions{Includes: []string{"available_in"}})
	if err != nil {
		t.Fatal(err)
	}
	plansInFacilities(t, l)
}
