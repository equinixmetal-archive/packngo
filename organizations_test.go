package packngo

import "testing"

func TestAccOrg(t *testing.T) {
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
	p, _, err = c.Organizations.Update(p.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}
	gotProject, _, err := c.Organizations.Get(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotProject.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotProject.Name)
	}
	_, err = c.Organizations.Delete(p.ID)
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
