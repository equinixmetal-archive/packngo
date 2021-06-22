package packngo

import (
	"path"
)

const subnetRouterBasePath = "/subnet-routers"

// DevicePortService handles operations on a port which belongs to a particular device
type ProjectSubnetRouterService interface {
	List(projectID string, opts *ListOptions) ([]SubnetRouter, *Response, error)
	Create(projectID string, input *SubnetRouterCreateRequest) (*SubnetRouter, *Response, error)
	Get(subnetRouterID string, opts *GetOptions) (*SubnetRouter, *Response, error)
	Delete(subnetRouterID string) (*Response, error)
}

type SubnetRouter struct {
	ID             string                `json:"id"`
	State          string                `json:"state"`
	Project        *Project              `json:"project,omitempty"`
	VirtualNetwork *VirtualNetwork       `json:"virtual_network,omitempty"`
	IPReservation  *IPAddressReservation `json:"ip_reservation,omitempty"`
	Href           string                `json:"href"`
	Created        string                `json:"created_at,omitempty"`
	Updated        string                `json:"updated_at,omitempty"`
}

type ProjectSubnetRouterServiceOp struct {
	client *Client
}

func (s *ProjectSubnetRouterServiceOp) List(projectID string, opts *ListOptions) (subnetRouters []SubnetRouter, resp *Response, err error) {
	type subnetRoutersRoot struct {
		SubnetRouters []SubnetRouter `json:"subnet_routers"`
		Meta          meta           `json:"meta"`
	}

	endpointPath := path.Join(projectBasePath, projectID, subnetRouterBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(subnetRoutersRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		subnetRouters = append(subnetRouters, subset.SubnetRouters...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}

}

type SubnetRouterCreateRequest struct {
	VirtualNetworkID      string `json:"virtual_network"`
	IPReservationID       string `json:"ip_reservation,omitempty"`
	PrivateIPv4SubnetSize int    `json:"private_ipv4_subnet_size,omitempty"`
}

func (s *ProjectSubnetRouterServiceOp) Get(subnetRouterID string, opts *GetOptions) (*SubnetRouter, *Response, error) {
	endpointPath := path.Join(subnetRouterBasePath, subnetRouterID)
	apiPathQuery := opts.WithQuery(endpointPath)
	subnetRouter := new(SubnetRouter)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, subnetRouter)
	if err != nil {
		return nil, resp, err
	}

	return subnetRouter, resp, err
}

func (s *ProjectSubnetRouterServiceOp) Create(projectID string, input *SubnetRouterCreateRequest) (*SubnetRouter, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, subnetRouterBasePath)
	output := new(SubnetRouter)

	resp, err := s.client.DoRequest("POST", apiPath, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

func (s *ProjectSubnetRouterServiceOp) Delete(subnetRouterID string) (*Response, error) {
	apiPath := path.Join(subnetRouterBasePath, subnetRouterID)

	resp, err := s.client.DoRequest("DELETE", apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
