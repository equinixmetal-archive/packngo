package packngo

import "fmt"

const transferRequestBasePath = "/transfers"

// TransferRequestsService interface defines available transfer request functions
type TransferRequestsService interface {
	Get(string, *ListOptions) (*TransferRequest, *Response, error)
	List(string, *ListOptions) ([]TransferRequest, *Response, error)
	Accept(string) (*Response, error)
	Decline(string) (*Response, error)
	TransferProject(string, string) (*Response, error)
}

// TransferRequestsServiceOp implements TransferRequestsService
type TransferRequestsServiceOp struct {
	client *Client
}

type transferRequestRoot struct {
	Transfers []TransferRequest `json:"transfers,omitempty"`
	Meta      meta              `json:"meta"`
}

// TransferRequest struct
type TransferRequest struct {
	ID                 string       `json:"id,omitempty"`
	CreatedAt          Timestamp    `json:"created_at,omitempty"`
	UpdatedAt          Timestamp    `json:"updated_at,omitempty"`
	TargetOrganization Organization `json:"target_organization,omitempty"`
	Project            Project      `json:"project,omitempty"`
	Href               string       `json:"href,omitempty"`
}

// TransferProject allows organization owners can transfer their projects to other organizations.
func (s *TransferRequestsServiceOp) TransferProject(projectID, organizationID string) (resp *Response, err error) {
	path := fmt.Sprintf("%s/%s%s", projectBasePath, projectID, transferRequestBasePath)

	body := map[string]string{}
	body["target_organization_id"] = organizationID

	resp, err = s.client.DoRequest("POST", path, body, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// List retrieves all project transfer requests from or to an organization
func (s *TransferRequestsServiceOp) List(organizationID string, listOpt *ListOptions) (transfers []TransferRequest, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s%s?%s", organizationBasePath, organizationID, transferRequestBasePath, params)
	for {
		subset := new(transferRequestRoot)
		resp, err = s.client.DoRequest("GET", path, nil, subset)
		if err != nil {
			return nil, resp, err
		}
		transfers = append(transfers, subset.Transfers...)

		if subset.Meta.Next != nil && (listOpt == nil || listOpt.Page == 0) {
			path = subset.Meta.Next.Href
			if params != "" {
				path = fmt.Sprintf("%s&%s", path, params)
			}
			continue
		}

		return
	}
}

// Get returns a single transfer request.
func (s *TransferRequestsServiceOp) Get(transferRequestID string, listOpt *ListOptions) (transferRequest *TransferRequest, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s?%s", transferRequestBasePath, transferRequestID, params)
	resp, err = s.client.DoRequest("GET", path, nil, transferRequest)
	if err != nil {
		return nil, resp, err
	}

	return transferRequest, resp, err
}

// Accept a transfer request
func (s *TransferRequestsServiceOp) Accept(projectID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", transferRequestBasePath, projectID)

	return s.client.DoRequest("PUT", path, nil, nil)
}

// Decline a transfer request
func (s *TransferRequestsServiceOp) Decline(projectID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", transferRequestBasePath, projectID)

	return s.client.DoRequest("DELETE", path, nil, nil)
}
