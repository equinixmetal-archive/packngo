package packngo

import (
	"path"
)

type ConnectionRedundancy string
type ConnectionType string

const (
	connectionBasePath                           = "/connections"
	virtualCircuitsBasePath                      = "/virtual-circuits"
	ConnectionShared        ConnectionType       = "shared"
	ConnectionDedicated     ConnectionType       = "dedicated"
	ConnectionRedundant     ConnectionRedundancy = "redundant"
	ConnectionPrimary       ConnectionRedundancy = "primary"
)

type ConnectionService interface {
	OrganizationCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	ProjectCreate(string, *ConnectionCreateRequest) (*Connection, *Response, error)
	OrganizationList(string, *GetOptions) ([]Connection, *Response, error)
	ProjectList(string, *GetOptions) ([]Connection, *Response, error)
	Delete(string) (*Response, error)
	Get(string, *GetOptions) (*Connection, *Response, error)
	Events(string, *GetOptions) ([]Event, *Response, error)
	PortEvents(string, string, *GetOptions) ([]Event, *Response, error)
	VirtualCircuitEvents(string, *GetOptions) ([]Event, *Response, error)
	Ports(string, *GetOptions) ([]ConnectionPort, *Response, error)
	Port(string, string, *GetOptions) (*ConnectionPort, *Response, error)
	VirtualCircuits(string, string, *GetOptions) ([]ConnectionVirtualCircuit, *Response, error)
	VirtualCircuit(string, *GetOptions) (*ConnectionVirtualCircuit, *Response, error)
	DeleteVirtualCircuit(string) (*Response, error)
}

type ConnectionServiceOp struct {
	client *Client
}

type connectionPortsRoot struct {
	Ports []ConnectionPort `json:"ports"`
}

type virtualCircuitsRoot struct {
	VirtualCircuits []ConnectionVirtualCircuit `json:"virtual_circuits"`
	Meta            meta                       `json:"meta"`
}

type connectionsRoot struct {
	Connections []Connection `json:"interconnections"`
	Meta        meta         `json:"meta"`
}

type ConnectionVirtualCircuit struct {
	ID      string          `json:"id"`
	Name    string          `json:"name,omitempty"`
	Status  string          `json:"status,omitempty"`
	VNID    string          `json:"vnid,omitempty"`
	NniVNID string          `json:"nni_vnid,omitempty"`
	NniVLAN string          `json:"nni_vlan,omitempty"`
	Project *Project        `json:"project,omitempty"`
	Port    *ConnectionPort `json:"port,omitempty"`
}

type ConnectionPort struct {
	ID              string                     `json:"id"`
	Name            string                     `json:"name,omitempty"`
	Status          string                     `json:"status,omitempty"`
	Role            string                     `json:"role,omitempty"`
	Speed           string                     `json:"speed,omitempty"`
	Organization    *Organization              `json:"organization,omitempty"`
	VirtualCircuits []ConnectionVirtualCircuit `json:"virtual_circuits,omitempty"`
	LinkStatus      string                     `json:"link_status,omitempty"`
	Href            string                     `json:"href,omitempty"`
}

type Connection struct {
	ID           string               `json:"id"`
	Name         string               `json:"name,omitempty"`
	Status       string               `json:"status,omitempty"`
	Redundancy   ConnectionRedundancy `json:"redundancy,omitempty"`
	Facility     *Facility            `json:"facility,omitempty"`
	Type         ConnectionType       `json:"type,omitempty"`
	Description  string               `json:"description,omitempty"`
	Project      *Project             `json:"project,omitempty"`
	Organization *Organization        `json:"organization,omitempty"`
	Speed        string               `json:"speed,omitempty"`
	Token        string               `json:"token,omitempty"`
	Tags         []string             `json:"tags,omitempty"`
	Ports        []ConnectionPort     `json:"ports,omitempty"`
}

type ConnectionCreateRequest struct {
	Name        string               `json:"name,omitempty"`
	Redundancy  ConnectionRedundancy `json:"redundancy,omitempty"`
	Facility    string               `json:"facility,omitempty"`
	Type        ConnectionType       `json:"type,omitempty"`
	Description *string              `json:"description,omitempty"`
	Project     string               `json:"project,omitempty"`
	Speed       string               `json:"speed,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
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

func (s *ConnectionServiceOp) Delete(id string) (*Response, error) {
	apiPath := path.Join(connectionBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

func (s *ConnectionServiceOp) Port(connID, portID string, opts *GetOptions) (*ConnectionPort, *Response, error) {
	endpointPath := path.Join(connectionBasePath, connID, portBasePath, portID)
	apiPathQuery := opts.WithQuery(endpointPath)
	port := new(ConnectionPort)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, port)
	if err != nil {
		return nil, resp, err
	}
	return port, resp, err
}

func (s *ConnectionServiceOp) Get(id string, opts *GetOptions) (*Connection, *Response, error) {
	endpointPath := path.Join(connectionBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	connection := new(Connection)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, connection)
	if err != nil {
		return nil, resp, err
	}
	return connection, resp, err
}

func (s *ConnectionServiceOp) Ports(connID string, opts *GetOptions) ([]ConnectionPort, *Response, error) {
	endpointPath := path.Join(connectionBasePath, connID, portBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	ports := new(connectionPortsRoot)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, ports)
	if err != nil {
		return nil, resp, err
	}
	return ports.Ports, resp, nil

}

func (s *ConnectionServiceOp) Events(id string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(connectionBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) PortEvents(connID, portID string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(connectionBasePath, connID, portBasePath, portID, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) VirtualCircuitEvents(id string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(virtualCircuitsBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *ConnectionServiceOp) VirtualCircuits(connID, portID string, opts *GetOptions) (vcs []ConnectionVirtualCircuit, resp *Response, err error) {
	endpointPath := path.Join(connectionBasePath, connID, portBasePath, portID, virtualCircuitsBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	for {
		subset := new(virtualCircuitsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		vcs = append(vcs, subset.VirtualCircuits...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

func (s *ConnectionServiceOp) VirtualCircuit(id string, opts *GetOptions) (*ConnectionVirtualCircuit, *Response, error) {
	endpointPath := path.Join(virtualCircuitsBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	vc := new(ConnectionVirtualCircuit)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, vc)
	if err != nil {
		return nil, resp, err
	}
	return vc, resp, err
}

func (s *ConnectionServiceOp) DeleteVirtualCircuit(id string) (*Response, error) {
	apiPath := path.Join(virtualCircuitsBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
