package packngo

import "fmt"

// API documentation https://www.packet.net/developers/api/organizations/
const organizationBasePath = "/organizations"

// OrganizationService interface defines available organization methods
type OrganizationService interface {
	List() ([]Organization, *Response, error)
	Get(string) (*Organization, *Response, error)
	Create(*OrganizationCreateRequest) (*Organization, *Response, error)
	Update(*OrganizationUpdateRequest) (*Organization, *Response, error)
	Delete(string) (*Response, error)
}

type organizationsRoot struct {
	Organizations []Organization `json:"organizations"`
}

// Organization represents a Packet organization
type Organization struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	Website      string    `json:"website,omitempty"`
	Twitter      string    `json:"twitter,omitempty"`
	Created      string    `json:"created_at,omitempty"`
	Updated      string    `json:"updated_at,omitempty"`
	Address      Address   `json:"address,omitempty"`
	TaxID        string    `json:"tax_id,omitempty"`
	MainPhone    string    `json:"main_phone,omitempty"`
	BillingPhone string    `json:"billing_phone,omitempty"`
	CreditAmount float64   `json:"credit_amount,omitempty"`
	Logo         string    `json:"logo,omitempty"`
	LogoThumb    string    `json:"logo_thumb,omitempty"`
	Projects     []Project `json:"projects,omitempty"`
	URL          string    `json:"href,omitempty"`
	Users        []User    `json:"members,omitempty"`
	Owners       []User    `json:"owners,omitempty"`
}

func (o Organization) String() string {
	return Stringify(o)
}

// OrganizationCreateRequest type used to create a Packet organization
type OrganizationCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Twitter     string `json:"twitter"`
	Logo        string `json:"logo"`
}

func (o OrganizationCreateRequest) String() string {
	return Stringify(o)
}

// OrganizationUpdateRequest type used to update a Packet organization
type OrganizationUpdateRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Twitter     string `json:"twitter"`
	Logo        string `json:"logo"`
}

func (o OrganizationUpdateRequest) String() string {
	return Stringify(o)
}

// OrganizationServiceOp implements OrganizationService
type OrganizationServiceOp struct {
	client *Client
}

// List returns the user's organizations
func (s *OrganizationServiceOp) List() ([]Organization, *Response, error) {
	root := new(organizationsRoot)

	resp, err := s.client.DoRequest("GET", organizationBasePath, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Organizations, resp, err
}

// Get returns a organization by id
func (s *OrganizationServiceOp) Get(organizationID string) (*Organization, *Response, error) {
	path := fmt.Sprintf("%s/%s", organizationBasePath, organizationID)
	organization := new(Organization)

	resp, err := s.client.DoRequest("GET", path, nil, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Create creates a new organization
func (s *OrganizationServiceOp) Create(createRequest *OrganizationCreateRequest) (*Organization, *Response, error) {
	organization := new(Organization)

	resp, err := s.client.DoRequest("POST", organizationBasePath, createRequest, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Update updates an organization
func (s *OrganizationServiceOp) Update(updateRequest *OrganizationUpdateRequest) (*Organization, *Response, error) {
	path := fmt.Sprintf("%s/%s", organizationBasePath, updateRequest.ID)
	organization := new(Organization)

	resp, err := s.client.DoRequest("PATCH", path, updateRequest, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, err
}

// Delete deletes an organizationID
func (s *OrganizationServiceOp) Delete(organizationID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", organizationBasePath, organizationID)

	return s.client.DoRequest("DELETE", path, nil, nil)
}
