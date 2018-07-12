package packngo

import (
	"fmt"
)

const batchBasePath = "/batches"

// BatchService interface defines available batch methods
type BatchService interface {
	Get(batchID string, listOpt *ListOptions) (*Batch, *Response, error)
	List(ProjectID string, listOpt *ListOptions) ([]Batch, *Response, error)
	Create(projectID string, batches *InstanceBatchCreateRequest) ([]Batch, *Response, error)
}

// Batch type
type Batch struct {
	ID                     string     `json:"id"`
	State                  string     `json:"state,omitempty"`
	Quantity               int32      `json:"quantity,omitempty"`
	CreatedAt              *Timestamp `json:"created_at,omitempty"`
	Href                   string     `json:"href,omitempty"`
	Project                Href       `json:"project,omitempty"`
	Instances              Href       `json:"instances,omitempty"`
	Facilities             []Facility `json:"facilities,omitempty"`
	FacilityDiversityLevel int32      `json:"facility_diversity_level,omitempty"`
}

//BatchesList represents collection of batches
type batchesList struct {
	Batches []Batch `json:"batches,omitempty"`
}

// InstanceBatchCreateRequest type used to create batch of device instances
type InstanceBatchCreateRequest struct {
	Batches []BatchInstance `json:"batches"`
}

// BatchInstance type used to describe batch instances
type BatchInstance struct {
	Plan            string     `json:"plan"`
	Hostname        string     `json:"hostname"`
	Facility        string     `json:"facility"`
	BillingCycle    string     `json:"billing_cycle"`
	OperatingSystem string     `json:"operating_system"`
	Quantity        int        `json:"quantity"`
	Hostnames       []string   `json:"hostnames,omitempty"`
	Description     string     `json:"description,omitempty"`
	AlwaysPxe       bool       `json:"always_pxe,omitempty"`
	Userdata        string     `json:"userdata,omitempty"`
	Locked          bool       `json:"locked,omitempty"`
	TerminationTime *Timestamp `json:"termination_time,omitempty"`
	Tags            []string   `json:"tags,omitempty"`
	ProjectSSHKeys  []string   `json:"project_ssh_keys,omitempty"`
	UserSSSHKeys    []string   `json:"user_ssh_keys,omitempty"`
	Features        []string   `json:"features,omitempty"`
	Customdata      string     `json:"customdata,omitempty"`
}

// BatchServiceOp implements BatchService
type BatchServiceOp struct {
	client *Client
}

// Get returns batch details
func (s *BatchServiceOp) Get(batchID string, listOpt *ListOptions) (*Batch, *Response, error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s?%s", batchBasePath, batchID, params)
	batch := new(Batch)

	resp, err := s.client.DoRequest("GET", path, nil, batch)
	if err != nil {
		return nil, resp, err
	}

	return batch, resp, err
}

// List returns batches on a project
func (s *BatchServiceOp) List(projectID string, listOpt *ListOptions) (batches []Batch, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s%s?%s", projectBasePath, projectID, batchBasePath, params)
	subset := new(batchesList)
	resp, err = s.client.DoRequest("GET", path, nil, subset)
	if err != nil {
		return nil, resp, err
	}

	batches = append(batches, subset.Batches...)
	return batches, resp, err
}

// Create function to create batch of device instances
func (s *BatchServiceOp) Create(projectID string, request *InstanceBatchCreateRequest) ([]Batch, *Response, error) {
	path := fmt.Sprintf("%s/%s/devices/batch", projectBasePath, projectID)

	batches := new(batchesList)
	resp, err := s.client.DoRequest("POST", path, request, batches)

	if err != nil {
		return nil, resp, err
	}

	return batches.Batches, resp, err
}
