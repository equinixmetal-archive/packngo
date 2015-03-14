package packngo

import "fmt"

const projectBasePath = "/projects"

type ProjectService interface {
	List() ([]Project, *Response, error)
	Get(string) (*Project, *Response, error)
	Create(*ProjectCreateRequest) (*Project, *Response, error)
	Update(*ProjectUpdateRequest) (*Project, *Response, error)
}

type ProjectsRoot struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	ID      string   `json:"id"`
	Name    string   `json:"name,omitempty"`
	Created string   `json:"created_at,omitempty"`
	Updated string   `json:"updated_at,omitempty"`
	Users   []User   `json:"members,omitempty"`
	Devices []Device `json:"devices,omitempty"`
	SshKeys []SshKey `json:"ssh_keys,omitempty"`
	Url     string   `json:"href,omitempty"`
}
func (p Project) String() string {
	return Stringify(p)
}

type ProjectCreateRequest struct {
	Name string `json:"name"`
}
func (p ProjectCreateRequest) String() string {
	return Stringify(p)
}

type ProjectUpdateRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
func (p ProjectUpdateRequest) String() string {
	return Stringify(p)
}

type ProjectServiceOp struct {
	client *Client
}

func (s *ProjectServiceOp) List() ([]Project, *Response, error) {
	req, err := s.client.NewRequest("GET", projectBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectsRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Projects, resp, err
}

func (s *ProjectServiceOp) Get(projectId string) (*Project, *Response, error) {
	path := fmt.Sprintf("%s/%s", projectBasePath, projectId)
	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

func (s *ProjectServiceOp) Create(createRequest *ProjectCreateRequest) (*Project, *Response, error) {
	req, err := s.client.NewRequest("POST", projectBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

func (s *ProjectServiceOp) Update(updateRequest *ProjectUpdateRequest) (*Project, *Response, error) {
	path := fmt.Sprintf("%s/%s", projectBasePath, updateRequest.Id)
	req, err := s.client.NewRequest("PATCH", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	project := new(Project)
	resp, err := s.client.Do(req, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}
