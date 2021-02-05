package packngo

import (
	"path"
)

type PortServiceOp struct {
	client *Client
}

// PortService handles operations on a port
type PortService interface {
	Assign(*PortAssignRequest) (*Port, *Response, error)
	Unassign(*PortAssignRequest) (*Port, *Response, error)
	AssignNative(*PortAssignRequest) (*Port, *Response, error)
	UnassignNative(string) (*Port, *Response, error)
	Bond(string, bool) (*Port, *Response, error)
	Disbond(string, bool) (*Port, *Response, error)
	ConvertToLayerTwo(string, string) (*Port, *Response, error)
	ConvertToLayerThree(string, []AddressRequest) (*Port, *Response, error)
	Get(string, *GetOptions) (*Port, *Response, error)
}

var _ PortService = &PortServiceOp{}

// Assign adds a VLAN to a port
func (i *PortServiceOp) Assign(par *PortAssignRequest) (*Port, *Response, error) {
	apiPath := path.Join(portBasePath, par.PortID, "assign")
	return i.portAction(apiPath, par)
}

// AssignNative assigns a virtual network to the port as a "native VLAN"
func (i *PortServiceOp) AssignNative(par *PortAssignRequest) (*Port, *Response, error) {
	apiPath := path.Join(portBasePath, par.PortID, "native-vlan")
	return i.portAction(apiPath, par)
}

// UnassignNative removes native VLAN from the supplied port
func (i *PortServiceOp) UnassignNative(portID string) (*Port, *Response, error) {
	apiPath := path.Join(portBasePath, portID, "native-vlan")
	port := new(Port)

	resp, err := i.client.DoRequest("DELETE", apiPath, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// Unassign removes a VLAN from the port
func (i *PortServiceOp) Unassign(par *PortAssignRequest) (*Port, *Response, error) {
	apiPath := path.Join(portBasePath, par.PortID, "unassign")
	return i.portAction(apiPath, par)
}

// Bond enables bonding for one or all ports
func (i *PortServiceOp) Bond(portID string, bulkEnable bool) (*Port, *Response, error) {
	br := &BondRequest{PortID: portID, BulkEnable: bulkEnable}
	apiPath := path.Join(portBasePath, br.PortID, "bond")
	return i.portAction(apiPath, br)
}

// Disbond disables bonding for one or all ports
func (i *PortServiceOp) Disbond(portID string, bulkEnable bool) (*Port, *Response, error) {
	dr := &DisbondRequest{PortID: portID, BulkDisable: bulkEnable}
	apiPath := path.Join(portBasePath, dr.PortID, "disbond")
	return i.portAction(apiPath, dr)
}

func (i *PortServiceOp) portAction(apiPath string, req interface{}) (*Port, *Response, error) {
	port := new(Port)

	resp, err := i.client.DoRequest("POST", apiPath, req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// ConvertToLayerTwo converts a bond port to Layer 2. IP assignments of the port will be removed.
func (i *PortServiceOp) ConvertToLayerTwo(portID, portName string) (*Port, *Response, error) {
	apiPath := path.Join(portBasePath, portID, "convert", "layer-2")
	port := new(Port)

	resp, err := i.client.DoRequest("POST", apiPath, nil, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// ConvertToLayerThree converts a bond port to Layer 3. VLANs must first be unassigned.
func (i *PortServiceOp) ConvertToLayerThree(portID string, ips []AddressRequest) (*Port, *Response, error) {
	apiPath := path.Join(portBasePath, portID, "convert", "layer-3")
	port := new(Port)

	req := BackToL3Request{
		RequestIPs: ips,
	}

	resp, err := i.client.DoRequest("POST", apiPath, &req, port)
	if err != nil {
		return nil, resp, err
	}

	return port, resp, err
}

// Get returns a port by id
func (s *PortServiceOp) Get(portID string, opts *GetOptions) (*Port, *Response, error) {
	endpointPath := path.Join(portBasePath, portID)
	apiPathQuery := opts.WithQuery(endpointPath)
	port := new(Port)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, port)
	if err != nil {
		return nil, resp, err
	}
	return port, resp, err
}
