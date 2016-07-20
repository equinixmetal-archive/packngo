package packngo

import "fmt"

const storageBasePath = "/storage"

// StoraveService interface defines available Storage methods
type StorageService interface {
  Get(string) (*Storage, *Response, error)
  Update(*StorageUpdateRequest) (*Storage, *Response, error)
  Delete(string) (*Response, error)
  Create(*StorageCreateRequest) (*Storage, *Response, error)
}

// Storage represents a storage
type Storage struct {
  ID               string           `json:"id"`
  Name             string           `json:"name,omitempty"`
  Description      string           `json:"description,omitempty"`
  Size             int              `json:"size,omitempty"`
  State            string           `json:"state,omitempty"`
  Locked           bool             `json:"locked,omitempty"`
  BillingCycle     string           `json:"billing_cycle,omitempty"`
  Created          string           `json:"created_at,omitempty"`
	Updated          string           `json:"updated_at,omitempty"`
  Href             string           `json:"href,omitempty"`
  SnapshotPolicies []SnapshotPolicy `json:"snapshot_policies,omitempty"`
  Attachments      []Attachment     `json:"attachments,omitempty"`
  Plan             *Plan            `json:"plan,omitempty"`
  Facility         *Facility        `json:"facility,omitempty"`
  Project          *Project         `json:"project,omitempty"`
}

// SnapshotPolicy used to execute actions on storage
type SnapshotPolicy struct {
  ID    string    `json:"id"`
  Href  string    `json:"href"`
}

// Attachment used to execute actions on storage
type Attachment struct {
  ID    string    `json:"id"`
  Href  string    `json:"href"`
}

func (s Storage) String() string {
	return Stringify(s)
}

// StorageUpdateRequest type used to update a Packet storage
type StorageUpdateRequest struct {
	ID            string   `json:"id"`
	Description   string   `json:"description,omitempty"`
	Size          int      `json:"size,omitempty"`
  Locked        bool     `json:"locked",omitempty`
}

func (p StorageUpdateRequest) String() string {
	return Stringify(p)
}

// StorageServiceOp implements StorageService
type StorageServiceOp struct {
	client *Client
}

// Get returns a storage by id
func (s *StorageServiceOp) Get(storageID string) (*Storage, *Response, error) {
	path := fmt.Sprintf("%s/%s", storageBasePath, storageID)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	storage := new(Storage)
	resp, err := s.client.Do(req, storage)
	if err != nil {
		return nil, resp, err
	}

	return storage, resp, err
}

// Update updates a storage
func (s *StorageServiceOp) Update(updateRequest *StorageUpdateRequest) (*Storage, *Response, error) {
	path := fmt.Sprintf("%s/%s", storageBasePath, updateRequest.ID)
	req, err := s.client.NewRequest("PATCH", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	storage := new(Storage)
	resp, err := s.client.Do(req, storage)
	if err != nil {
		return nil, resp, err
	}

	return storage, resp, err
}

// Delete deletes a storage
func (s *StorageServiceOp) Delete(storageID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", storageBasePath, storageID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// StorageCreateRequest type used to create a Packet storage
type StorageCreateRequest struct {
  Name          string   `json:"name"`
  Size          int      `json:"size"`
  BillingCycle  string   `json:"billing_cycle"`
	ProjectID     string   `json:"project_id"`
  PlanID        string   `json:"plan_id"`
  FacilityID    string   `json:"facility_id"`
}

func (s StorageCreateRequest) String() string {
	return Stringify(s)
}

// Create creates a new storage for a project
func (s *StorageServiceOp) Create(createRequest *StorageCreateRequest) (*Storage, *Response, error) {
	url := fmt.Sprintf("%s/%s%s", projectBasePath, createRequest.ProjectID, storageBasePath)
	req, err := s.client.NewRequest("POST", url, createRequest)
	if err != nil {
		return nil, nil, err
	}

	storage := new(Storage)
	resp, err := s.client.Do(req, storage)
	if err != nil {
		return nil, resp, err
	}

	return storage, resp, err
}
