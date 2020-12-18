package packngo

import (
	"path"
)

const (
	connectionBasePath = "/connections"
)

type ConnectionService interface {
	OrganizationCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	ProjectCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	OrganizationList(string, *GetOptions) ([]Connection, *Response, error)
	ProjectList(string, *GetOptions) ([]Connection, *Response, error)
}

type ConnectionServiceOp struct {
	client *Client
}

type connectionsRoot struct {
	Connections []Connection `json:"connections"`
	Meta        meta         `json:"meta"`
}

type Connection struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Redundancy  string   `json:"redundancy,omitempty"`
	Facility    string   `json:"facility,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description *string  `json:"description,omitempty"`
	Project     string   `json:"string,omitempty"`
	Speed       string   `json:"speed,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type ConnectionCreateRequest struct {
	Name        string   `json:"name,omitempty"`
	Redundancy  string   `json:"redundancy,omitempty"`
	Facility    string   `json:"facility,omitempty"`
	Type        string   `json:"type,omitempty"`
	Description *string  `json:"description,omitempty"`
	Project     string   `json:"string,omitempty"`
	Speed       string   `json:"speed,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (s *ConnectionServiceOp) create(apiUrl string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	connection := new(Connection)
	resp, err := s.client.DoRequest("POST", apiUrl, createRequest, connection)
	if err != nil {
		return nil, resp, err
	}

	return connection, resp, err
}

func (s *ConnectionServiceOp) OrganizationCreate(id string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	apiUrl := path.Join(organizationBasePath, id, connectionBasePath)
	return s.create(apiUrl, createRequest)
}

func (s *ConnectionServiceOp) ProjectCreate(id string, createRequest *ConnectionCreateRequest) (*Connection, *Response, error) {
	apiUrl := path.Join(projectBasePath, id, connectionBasePath)
	return s.create(apiUrl, createRequest)
}

func (s *ConnectionServiceOp) list(url string, opts *GetOptions) (connections []Connection, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(url)

	for {
		subset := new(connectionsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		connections = append(connections, subset.Connections...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}

}

func (s *ConnectionServiceOp) OrganizationList(id string, opts *GetOptions) ([]Connection, *Response, error) {
	apiUrl := path.Join(organizationBasePath, id, connectionBasePath)
	return s.list(apiUrl, opts)
}

func (s *ConnectionServiceOp) ProjectList(id string, opts *GetOptions) ([]Connection, *Response, error) {
	apiUrl := path.Join(projectBasePath, id, connectionBasePath)
	return s.list(apiUrl, opts)
}
