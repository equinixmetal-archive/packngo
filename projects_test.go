package packngo

import (
	"testing"
)

func TestAccProjectBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer projectTeardown(c)

	rs := testProjectPrefix + randString8()
	pcr := ProjectCreateRequest{Name: rs}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
	}
	rs = testProjectPrefix + randString8()
	pur := ProjectUpdateRequest{Name: &rs}
	p, _, err = c.Projects.Update(p.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}
	gotProject, _, err := c.Projects.Get(p.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotProject.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotProject.Name)
	}

	if gotProject.PaymentMethod.URL == "" {
		t.Fatalf("Empty payment_method: %v", gotProject)
	}

	_, err = c.Projects.Delete(p.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccProjectExtra(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer projectTeardown(c)
	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	rs := testProjectPrefix + randString8()
	pcr := ProjectCreateRequest{Name: rs}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
	}
	rs = testProjectPrefix + randString8()
	pur := ProjectUpdateRequest{Name: &rs}
	p, _, err = c.Projects.Update(p.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}
	gotProject, _, err := c.Projects.Get(p.ID, &GetOptions{Includes: []string{"members"}})
	if err != nil {
		t.Fatal(err)
	}

	for _, user := range gotProject.Users {
		if user.ID == "" {
			t.Fatal("Project user details not returned.")
		} else if user.ID == u.ID {
			break
		}
	}

	_, err = c.Projects.Delete(p.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccCreateOrgProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer projectTeardown(c)

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	rs := testProjectPrefix + randString8()

	orgPath := "/organizations/" + u.DefaultOrganizationID
	pcr := ProjectCreateRequest{Name: rs}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Organization.URL != orgPath {
		t.Fatalf("Expected new project to be part of org %s, not %v", orgPath, p.Organization)
	}
}

func TestAccCreateNonDefaultOrgProject(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c := setup(t)
	defer organizationTeardown(c)
	defer projectTeardown(c)

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	orgName := testProjectPrefix + randString8()
	ocr := OrganizationCreateRequest{
		Name:        orgName,
		Description: "Managed by Packngo.",
		Website:     "http://example.com",
		Twitter:     "foo",
	}
	org, _, err := c.Organizations.Create(&ocr)
	if err != nil {
		t.Fatal(err)
	}

	rs := testProjectPrefix + randString8()

	if org.ID == u.DefaultOrganizationID {
		t.Fatalf("Expected new organization %s to not have same ID as Default org %s", org.ID, u.DefaultOrganizationID)
	}

	orgPath := "/organizations/" + org.ID
	pcr := ProjectCreateRequest{Name: rs, OrganizationID: org.ID}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}

	if p.Organization.URL != orgPath {
		t.Fatalf("Expected new project to be part of org %s, not %v", orgPath, p.Organization)
	}

	defaultOrgPath := "/organizations/" + u.DefaultOrganizationID
	if p.Organization.URL == defaultOrgPath {
		t.Fatalf("Expected new project to not be part of org %s", orgPath)
	}
}

func TestAccListProjects(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)

	defer projectTeardown(c)

	rs := testProjectPrefix + randString8()

	u, _, err := c.Users.Current()
	if err != nil {
		t.Fatal(err)
	}

	orgPath := "/organizations/" + u.DefaultOrganizationID

	pcr := ProjectCreateRequest{Name: rs}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}

	if p.Organization.URL != orgPath {
		t.Fatalf("Expected new project to be part of org %s, not %v", orgPath, p.Organization)
	}

	listOpt := &ListOptions{
		Includes: []string{"members"},
	}
	projs, _, err := c.Projects.List(listOpt)
	if err != nil {
		t.Fatal(err)
	}

	for _, proj := range projs {
		if proj.ID == p.ID {
			if proj.Users[0].ID != u.ID {
				t.Fatal("Project user details not returned.")
			}
			break
		}
	}
}

func TestAccProjectListPagination(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c := setup(t)
	defer projectTeardown(c)
	for i := 0; i < 3; i++ {
		pcr := ProjectCreateRequest{
			Name: testProjectPrefix + randString8(),
		}
		_, _, err := c.Projects.Create(&pcr)
		if err != nil {
			t.Fatal(err)
		}
	}
	listOpts := &ListOptions{
		Page:    1,
		PerPage: 3,
	}

	projects, _, err := c.Projects.List(listOpts)
	if err != nil {
		t.Fatalf("failed to get list of projects: %v", err)
	}
	// The user account that runs this test probably have some projects on
	// his own, keep it in mind when improving/extending this test.
	if len(projects) != 3 {
		t.Fatalf("exactly 3 projects should have been fetched: %v", err)
	}

	listOpts = &ListOptions{
		Page:    2,
		PerPage: 1,
	}

	projects, _, err = c.Projects.List(listOpts)
	if err != nil {
		t.Fatalf("failed to get list of projects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("only 1 project should have been fetched: %v", err)
	}

}

func TestAccProjectListEvents(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	// create volume to generate some events

	sp := SnapshotPolicy{
		SnapshotFrequency: "1day",
		SnapshotCount:     3,
	}

	vcr := VolumeCreateRequest{
		Size:             10,
		BillingCycle:     "hourly",
		PlanID:           "storage_1",
		FacilityID:       testFacility(),
		SnapshotPolicies: []*SnapshotPolicy{&sp},
		Description:      "ahoj!",
	}

	v, _, err := c.Volumes.Create(&vcr, projectID)
	if err != nil {
		t.Fatal(err)
	}
	v, err = waitVolumeActive(v.ID, c)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.Volumes.Delete(v.ID)
	if err != nil {
		t.Fatal(err)
	}

	es, _, err := c.Projects.ListEvents(projectID, &ListOptions{PerPage: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(es) < 2 {
		t.Fatal("At least 2 events should be in project - volume add and del")
	}

}
