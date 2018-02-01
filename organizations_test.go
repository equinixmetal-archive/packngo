package packngo

import "testing"

func TestAccOrg(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer organizationTeardown(c)

	rs := testProjectPrefix + randString8()
	ocr := OrganizationCreateRequest{
		Name:        rs,
		Description: "Managed by Terraform.",
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
	pur := OrganizationUpdateRequest{
		ID:          p.ID,
		Name:        rs,
		Description: "Managed by Terraform.",
		Website:     "http://quux.example.com",
		Twitter:     "bar",
	}
	p, _, err = c.Organizations.Update(&pur)
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
func TestOrgListPaymentMethods(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	// setup
	c := setup(t)
	defer organizationTeardown(c)

	rs := testProjectPrefix + randString8()
	ocr := OrganizationCreateRequest{
		Name:        rs,
		Description: "Managed by Terraform.",
		Website:     "http://example.com",
		Twitter:     "foo",
	}
	org, _, err := c.Organizations.Create(&ocr)
	if err != nil {
		t.Fatal(err)
	}

	// tests
	_, _, err = c.Organizations.ListPaymentMethods(org.ID)
	if err != nil {
		t.Fatal("error: ", err)
	}

	// teardown
	_, err = c.Organizations.Delete(org.ID)
	if err != nil {
		t.Fatal(err)
	}

}
