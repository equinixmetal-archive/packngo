package packngo

import (
	"path"
)

type MetalGatewayState string

const (
	metalGatewayBasePath                   = "/metal-gateways"
	MetalGatewayActive   MetalGatewayState = "active"
	MetalGatewayReady    MetalGatewayState = "ready"
	MetalGatewayDeleting MetalGatewayState = "deleting"
)

type MetalGatewayService interface {
	List(projectID string, opts *ListOptions) ([]MetalGateway, *Response, error)
	Create(projectID string, input *MetalGatewayCreateRequest) (*MetalGateway, *Response, error)
	Get(metalGatewayID string, opts *GetOptions) (*MetalGateway, *Response, error)
	Delete(metalGatewayID string) (*Response, error)
}

type MetalGateway struct {
	ID             string                `json:"id"`
	State          MetalGatewayState     `json:"state"`
	Project        *Project              `json:"project,omitempty"`
	VirtualNetwork *VirtualNetwork       `json:"virtual_network,omitempty"`
	IPReservation  *IPAddressReservation `json:"ip_reservation,omitempty"`
	Href           string                `json:"href"`
	CreatedAt      string                `json:"created_at,omitempty"`
	UpdatedAt      string                `json:"updated_at,omitempty"`
}

type MetalGatewayServiceOp struct {
	client *Client
}

func (s *MetalGatewayServiceOp) List(projectID string, opts *ListOptions) (metalGateways []MetalGateway, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	type metalGatewaysRoot struct {
		MetalGateways []MetalGateway `json:"metal_gateways"`
		Meta          meta           `json:"meta"`
	}

	endpointPath := path.Join(projectBasePath, projectID, metalGatewayBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(metalGatewaysRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		metalGateways = append(metalGateways, subset.MetalGateways...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}

}

type MetalGatewayCreateRequest struct {
	VirtualNetworkID      string `json:"virtual_network_id"`
	IPReservationID       string `json:"ip_reservation_id,omitempty"`
	PrivateIPv4SubnetSize int    `json:"private_ipv4_subnet_size,omitempty"`
}

func (s *MetalGatewayServiceOp) Get(metalGatewayID string, opts *GetOptions) (*MetalGateway, *Response, error) {
	if validateErr := ValidateUUID(metalGatewayID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(metalGatewayBasePath, metalGatewayID)
	apiPathQuery := opts.WithQuery(endpointPath)
	metalGateway := new(MetalGateway)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, metalGateway)
	if err != nil {
		return nil, resp, err
	}

	return metalGateway, resp, err
}

func (s *MetalGatewayServiceOp) Create(projectID string, input *MetalGatewayCreateRequest) (*MetalGateway, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(projectBasePath, projectID, metalGatewayBasePath)
	output := new(MetalGateway)

	resp, err := s.client.DoRequest("POST", apiPath, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

func (s *MetalGatewayServiceOp) Delete(metalGatewayID string) (*Response, error) {
	if validateErr := ValidateUUID(metalGatewayID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(metalGatewayBasePath, metalGatewayID)

	resp, err := s.client.DoRequest("DELETE", apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
