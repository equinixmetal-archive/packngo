package packngo

import (
	"log"
	"testing"
)

func TestAccOrgList(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer organizationTeardown(c)

	rs := testProjectPrefix + randString8()
	ocr := OrganizationCreateRequest{
		Name:        rs,
		Description: "Managed by Packngo.",
		Website:     "http://example.com",
		Twitter:     "foo",
	}
	org, _, err := c.Organizations.Create(&ocr)
	if err != nil {
		t.Fatal(err)
	}
	if org.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, org.Name)
	}
	ol, _, err := c.Organizations.List(&ListOptions{PerPage: 1})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, o := range ol {
		if o.ID == org.ID {
			found = true
			break
		}
	}
	if !found {
		log.Println("Couldn't find created test org in org listing")
	}

	_, err = c.Organizations.Delete(org.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccOrgBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer organizationTeardown(c)

	rs := testProjectPrefix + randString8()
	ocr := OrganizationCreateRequest{
		Name:        rs,
		Description: "Managed by Packngo.",
		Website:     "http://example.com",
		Twitter:     "foo",
	}
	p, _, err := c.Organizations.Create(&ocr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
	}

	rs = testProjectPrefix + randString8()
	oDesc := "Managed by Packngo."
	oWeb := "http://quux.example.com"
	oTwi := "bar"
	pur := OrganizationUpdateRequest{
		Name:        &rs,
		Description: &oDesc,
		Website:     &oWeb,
		Twitter:     &oTwi,
	}
	org, _, err := c.Organizations.Update(p.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if org.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}
	gotOrg, _, err := c.Organizations.Get(org.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotOrg.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotOrg.Name)
	}

	_, _, err = c.Organizations.ListEvents(org.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Organizations.Delete(org.ID)
	if err != nil {
		t.Fatal(err)
	}

}
func TestAccOrgListPaymentMethods(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	// setup
	c := setup(t)
	defer organizationTeardown(c)

	rs := testProjectPrefix + randString8()
	ocr := OrganizationCreateRequest{
		Name:        rs,
		Description: "Managed by Packngo.",
		Website:     "http://example.com",
		Twitter:     "foo",
	}
	org, _, err := c.Organizations.Create(&ocr)
	if err != nil {
		t.Fatal(err)
	}

	// tests
	pms, _, err := c.Organizations.ListPaymentMethods(org.ID)

	if err != nil {
		t.Fatal("error: ", err)
	}

	if len(pms) != 0 {
		t.Fatal("the new test org should have no payment methods")
	}

	// teardown
	_, err = c.Organizations.Delete(org.ID)
	if err != nil {
		t.Fatal(err)
	}

}
