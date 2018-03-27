package packngo

import "fmt"

const (
	sshKeyBasePath = "/ssh-keys"
)

// SSHKeyService interface defines available device methods
type SSHKeyService interface {
	List(*ListOptions) ([]SSHKey, *Response, error)
	ProjectList(string, *ListOptions) ([]SSHKey, *Response, error)
	Get(string) (*SSHKey, *Response, error)
	Create(*SSHKeyCreateRequest) (*SSHKey, *Response, error)
	Update(*SSHKeyUpdateRequest) (*SSHKey, *Response, error)
	Delete(string) (*Response, error)
}

type sshKeyRoot struct {
	SSHKeys []SSHKey `json:"ssh_keys"`
	Meta    meta     `json:"meta"`
}

// SSHKey represents a user's ssh key
type SSHKey struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Key         string `json:"key"`
	FingerPrint string `json:"fingerprint"`
	Created     string `json:"created_at"`
	Updated     string `json:"updated_at"`
	User        User   `json:"user,omitempty"`
	URL         string `json:"href,omitempty"`
}

func (s SSHKey) String() string {
	return Stringify(s)
}

// SSHKeyCreateRequest type used to create an ssh key
type SSHKeyCreateRequest struct {
	Label     string `json:"label"`
	Key       string `json:"key"`
	ProjectID string `json:"-"`
}

func (s SSHKeyCreateRequest) String() string {
	return Stringify(s)
}

// SSHKeyUpdateRequest type used to update an ssh key
type SSHKeyUpdateRequest struct {
	ID    string `json:"id"`
	Label string `json:"label,omitempty"`
	Key   string `json:"key,omitempty"`
}

func (s SSHKeyUpdateRequest) String() string {
	return Stringify(s)
}

// SSHKeyServiceOp implements SSHKeyService
type SSHKeyServiceOp struct {
	client *Client
}

func (s *SSHKeyServiceOp) list(url string, listOpt *ListOptions) (sshKeys []SSHKey, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
		if params != "" {
			url = fmt.Sprintf("%s?%s", url, params)
		}
	}

	for {
		subset := new(sshKeyRoot)

		resp, err = s.client.DoRequest("GET", url, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		sshKeys = append(sshKeys, subset.SSHKeys...)

		if subset.Meta.Next != nil {
			url = subset.Meta.Next.Href
			if params != "" {
				url = fmt.Sprintf("%s&%s", url, params)
			}
			continue
		}

		return
	}
}

// ProjectList lists ssh keys of a project
func (s *SSHKeyServiceOp) ProjectList(projectID string, listOpt *ListOptions) ([]SSHKey, *Response, error) {
	return s.list(fmt.Sprintf("%s/%s%s", projectBasePath, projectID, sshKeyBasePath), listOpt)
}

// List returns a user's ssh keys
func (s *SSHKeyServiceOp) List(listOpt *ListOptions) ([]SSHKey, *Response, error) {
	return s.list(sshKeyBasePath, listOpt)
}

// Get returns an ssh key by id
func (s *SSHKeyServiceOp) Get(sshKeyID string) (*SSHKey, *Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, sshKeyID)
	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("GET", path, nil, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Create creates a new ssh key
func (s *SSHKeyServiceOp) Create(createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	path := sshKeyBasePath
	if createRequest.ProjectID != "" {
		path = fmt.Sprintf("%s/%s%s", projectBasePath, createRequest.ProjectID, sshKeyBasePath)
	}
	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("POST", path, createRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Update updates an ssh key
func (s *SSHKeyServiceOp) Update(updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Response, error) {
	if updateRequest.Label == "" && updateRequest.Key == "" {
		return nil, nil, fmt.Errorf("You must set either Label or Key string for SSH Key update")
	}
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, updateRequest.ID)

	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("PATCH", path, updateRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Delete deletes an ssh key
func (s *SSHKeyServiceOp) Delete(sshKeyID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, sshKeyID)

	return s.client.DoRequest("DELETE", path, nil, nil)
}
