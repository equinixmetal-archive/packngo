package packngo

import (
	"log"
	"testing"
)

func TestAccPlans(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	l, _, err := c.Plans.List(&ListOptions{Includes: []string{"available_in"}})

	avail := map[string][]string{}
	for _, p := range l {
		for _, f := range p.AvailableIn {
			if _, ok := avail[f.Code]; !ok {
				avail[f.Code] = []string{p.Slug}
			} else {
				avail[f.Code] = append(avail[f.Code], p.Slug)
			}
		}
		if p.Pricing.Hour < 0.0 {
			t.Fatalf("strange pricing for %s %s", p.Name, p.Slug)
		}
	}

	for f, ps := range avail {
		if len(ps) == 0 {
			t.Fatalf("no plans available in facility %s", f)
		}
		// prints plans available in facility
		log.Printf("%s: %+v\n", f, ps)
	}

	if len(l) == 0 {
		t.Fatal("Empty plans listing from the API")
	}

	if err != nil {
		t.Fatal(err)
	}
}
