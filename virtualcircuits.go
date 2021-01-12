package packngo

import "path"

const (
	virtualCircuitBasePath = "/virtual-circuits"
	vcStatusActive         = "active"
	vcStatusWaiting        = "waiting_on_customer_vlan"
	//vcStatusActivating     = "activating"
	//vcStatusDeactivating   = "deactivating"
)

type VirtualCircuitService interface {
	Get(string, *GetOptions) (*VirtualCircuit, *Response, error)
	Events(string, *GetOptions) ([]Event, *Response, error)
	Delete(string) (*Response, error)
	ConnectVLAN(string, string, *GetOptions) (*VirtualCircuit, *Response, error)
	RemoveVLAN(string, *GetOptions) (*VirtualCircuit, *Response, error)
}

type VCUpdateRequest struct {
	VirtualNetworkID *string `json:"vnid"`
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
	Status         string          `json:"status,omitempty"`
	VNID           int             `json:"vnid,omitempty"`
	NniVNID        int             `json:"nni_vnid,omitempty"`
	NniVLAN        int             `json:"nni_vlan,omitempty"`
	Project        *Project        `json:"project,omitempty"`
	VirtualNetwork *VirtualNetwork `json:"virtual_network,omitempty"`
}

func (s *VirtualCircuitServiceOp) ConnectVLAN(vcID, vlanID string, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	endpointPath := path.Join(virtualCircuitBasePath, vcID)
	apiPathQuery := opts.WithQuery(endpointPath)
	vc := new(VirtualCircuit)
	updateReq := VCUpdateRequest{VirtualNetworkID: &vlanID}
	resp, err := s.client.DoRequest("PUT", apiPathQuery, updateReq, vc)
	if err != nil {
		return nil, resp, err
	}
	return vc, resp, err
}

func (s *VirtualCircuitServiceOp) RemoveVLAN(vcID string, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	endpointPath := path.Join(virtualCircuitBasePath, vcID)
	apiPathQuery := opts.WithQuery(endpointPath)
	vc := new(VirtualCircuit)
	updateReq := VCUpdateRequest{VirtualNetworkID: nil}
	resp, err := s.client.DoRequest("PUT", apiPathQuery, updateReq, vc)
	if err != nil {
		return nil, resp, err
	}
	return vc, resp, err
}

func (s *VirtualCircuitServiceOp) Events(id string, opts *GetOptions) ([]Event, *Response, error) {
	apiPath := path.Join(virtualCircuitBasePath, id, eventBasePath)
	return listEvents(s.client, apiPath, opts)
}

func (s *VirtualCircuitServiceOp) Get(id string, opts *GetOptions) (*VirtualCircuit, *Response, error) {
	endpointPath := path.Join(virtualCircuitBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	vc := new(VirtualCircuit)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, vc)
	if err != nil {
		return nil, resp, err
	}
	return vc, resp, err
}

func (s *VirtualCircuitServiceOp) Delete(id string) (*Response, error) {
	apiPath := path.Join(virtualCircuitBasePath, id)
	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
