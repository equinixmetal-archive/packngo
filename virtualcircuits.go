package packngo

import "path"

const (
	virtualCircuitBasePath = "/virtual-circuits"

	// VC is being create but not ready yet
	VCStatusPending = "pending"

	// VC is ready with a VLAN
	VCStatusActive = "active"

	// VC is ready without a VLAN
	VCStatusWaiting = "waiting_on_customer_vlan"

	// VC is being deleted
	VCStatusDeleting = "deleting"

	// not sure what the following states mean, or whether they exist
	// someone from the API side could check
	VCStatusActivating         = "activating"
	VCStatusDeactivating       = "deactivating"
	VCStatusActivationFailed   = "activation_failed"
	VCStatusDeactivationFailed = "dactivation_failed"
)

type VirtualCircuitService interface {
	Create(string, string, string, *VCCreateRequest, *GetOptions) (*VirtualCircuit, *Response, error)
	Get(string, *GetOptions) (*VirtualCircuit, *Response, error)
	Events(string, *GetOptions) ([]Event, *Response, error)
	Delete(string) (*Response, error)
	Update(string, *VCUpdateRequest, *GetOptions) (*VirtualCircuit, *Response, error)
}

type VCUpdateRequest struct {
	Name             *string   `json:"name,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
	Description      *string   `json:"description,omitempty"`
	VirtualNetworkID *string   `json:"vnid,omitempty"`

	// Speed is a bps representation of the VirtualCircuit throughput. This is informational only, the field is a user-controlled description of the speed. It may be presented as a whole number with a bps, mpbs, or gbps suffix (or the respective initial).
	Speed string `json:"speed,omitempty"`
}

type VCCreateRequest struct {
	VirtualNetworkID string   `json:"vnid"`
	NniVLAN          int      `json:"nni_vlan,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	Tags             []string `json:"tags,omitempty"`

	// Speed is a bps representation of the VirtualCircuit throughput. This is informational only, the field is a user-controlled description of the speed. It may be presented as a whole number with a bps, mpbs, or gbps suffix (or the respective initial).
	Speed string `json:"speed,omitempty"`
}

type VirtualCircuitServiceOp struct {
	client *Client
}

type virtualCircuitsRoot struct {
	VirtualCircuits []VirtualCircuit `json:"virtual_circuits"`
	Meta            meta             `json:"meta"`
}

type VirtualCircuit struct {
	ID             string          `json:"id"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Speed          string          `json:"speed,omitempty"`
	Status         string          `json:"status,omitempty"`
	VNID           int             `json:"vnid,omitempty"`
	NniVNID        int             `json:"nni_vnid,omitempty"`
	NniVLAN        int             `json:"nni_vlan,omitempty"`
	Project        *Project        `json:"project,omitempty"`
	Port           *ConnectionPort `json:"port,omitempty"`
	VirtualNetwork *VirtualNetwork `json:"virtual_network,omitempty"`
	Tags           []string        `json:"tags,omitempty"`
}

func (s *VirtualCircuitServiceOp) do(method, apiPathQuery string, req interface{}) (*VirtualCircuit, *Response, error) {
	vc := new(VirtualCircuit)
	resp, err := s.client.DoRequest(method, apiPathQuery, req, vc)
	if err != nil {
		return nil, resp, err
	}
	return vc, resp, err
}

func (s *VirtualCircuitServiceOp) Update(vcID string, req *VCUpdateRequest, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	if validateErr := ValidateUUID(vcID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(virtualCircuitBasePath, vcID)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("PUT", apiPathQuery, req)
}

func (s *VirtualCircuitServiceOp) Events(id string, opts *GetOptions) ([]Event, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(virtualCircuitBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *VirtualCircuitServiceOp) Get(id string, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(virtualCircuitBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("GET", apiPathQuery, nil)
}

func (s *VirtualCircuitServiceOp) Delete(id string) (*Response, error) {
	if validateErr := ValidateUUID(id); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(virtualCircuitBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}

func (s *VirtualCircuitServiceOp) Create(projectID, connID, portID string, request *VCCreateRequest, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(connID); validateErr != nil {
		return nil, nil, validateErr
	}
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(projectBasePath, projectID, connectionBasePath, connID, portBasePath, portID, virtualCircuitBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)
	return s.do("POST", apiPathQuery, request)
}
