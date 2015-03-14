package packngo

import "fmt"

const sshKeyBasePath = "/ssh-keys"

type SshKeyService interface {
	List() ([]SshKey, *Response, error)
	Get(string) (*SshKey, *Response, error)
	Create(*SshKeyCreateRequest) (*SshKey, *Response, error)
	Update(*SshKeyUpdateRequest) (*SshKey, *Response, error)
	Delete(string) (*Response, error)
}

type SshKeyRoot struct {
	SshKeys []SshKey `json:"ssh_keys"`
}

type SshKey struct {
	ID          string    `json:"id"`
	Label       string    `json:"label"`
  Key         string    `json:"key"`
	FingerPrint string    `json:"fingerprint"`
	Created     string    `json:"created_at"`
	Updated     string    `json:"updated_at"`
	User        User      `json:"user,omitempty"`
	Url         string    `json:"href,omitempty"`
}
func (s SshKey) String() string {
	return Stringify(s)
}

type SshKeyCreateRequest struct {
	Label string   `json:"label"`
	Key   string   `json:"key"`
}
func (s SshKeyCreateRequest) String() string {
	return Stringify(s)
}

type SshKeyUpdateRequest struct {
	Id    string   `json:"id"`
	Label string   `json:"label"`
	Key   string   `json:"key"`
}
func (s SshKeyUpdateRequest) String() string {
	return Stringify(s)
}

type SshKeyServiceOp struct {
	client *Client
}

func (s *SshKeyServiceOp) List() ([]SshKey, *Response, error) {
	req, err := s.client.NewRequest("GET", sshKeyBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SshKeyRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SshKeys, resp, err
}

func (s *SshKeyServiceOp) Get(sshKeyId string) (*SshKey, *Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, sshKeyId)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	sshKey := new(SshKey)
	resp, err := s.client.Do(req, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

func (s *SshKeyServiceOp) Create(createRequest *SshKeyCreateRequest) (*SshKey, *Response, error) {
	req, err := s.client.NewRequest("POST", sshKeyBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	sshKey := new(SshKey)
	resp, err := s.client.Do(req, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

func (s *SshKeyServiceOp) Update(updateRequest *SshKeyUpdateRequest) (*SshKey, *Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, updateRequest.Id)
	req, err := s.client.NewRequest("PATCH", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	sshKey := new(SshKey)
	resp, err := s.client.Do(req, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

func (s *SshKeyServiceOp) Delete(sshKeyID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, sshKeyID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}
