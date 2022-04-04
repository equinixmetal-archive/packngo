package packngo

import (
	"path"
)

const (
	vrfBasePath = "/vrfs"
)

type VRFService interface {
	List(projectID string, opts *ListOptions) ([]VRF, *Response, error)
	Create(projectID string, input *VRFCreateRequest) (*VRF, *Response, error)
	Get(vrfID string, opts *GetOptions) (*VRF, *Response, error)
	Delete(vrfID string) (*Response, error)
}

type VRF struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	LocalASN    int      `json:"local_asn,omitempty"`
	IPRanges    []string `json:"ip_ranges,omitempty"`
	Project     *Project `json:"project,omitempty"`
	Metro       *Metro   `json:"metro,omitempty"`
	Href        string   `json:"href"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

type VRFCreateRequest struct {
	//  Metro id or code
	Metro string `json:"metro"`

	// Name is the name of the VRF. It must be unique per project.
	Name string `json:"name"`

	// Description of the VRF to be created.
	Description string `json:"description"`

	// LocalASN is the ASN of the local network.
	LocalASN int `json:"local_asn,omitempty"`

	// IPBlocks is a list of all IPv4 and IPv6 Ranges that will be available to
	// BGP Peers. IPv4 addresses must be /8 or smaller with a minimum size of
	// /29. IPv6 must be /56 or smaller with a minimum size of /64. Ranges must
	// not overlap other ranges within the VRF.
	IPBlocks []string `json:"ip_blocks,omitempty"`
}

type VRFServiceOp struct {
	client *Client
}

func (s *VRFServiceOp) List(projectID string, opts *ListOptions) (vrfs []VRF, resp *Response, err error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	type vrfsRoot struct {
		VRFs []VRF `json:"vrfs"`
		Meta meta  `json:"meta"`
	}

	endpointPath := path.Join(projectBasePath, projectID, vrfBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(vrfsRoot)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		vrfs = append(vrfs, subset.VRFs...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

func (s *VRFServiceOp) Get(vrfID string, opts *GetOptions) (*VRF, *Response, error) {
	if validateErr := ValidateUUID(vrfID); validateErr != nil {
		return nil, nil, validateErr
	}
	endpointPath := path.Join(vrfBasePath, vrfID)
	apiPathQuery := opts.WithQuery(endpointPath)
	metalGateway := new(VRF)

	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, metalGateway)
	if err != nil {
		return nil, resp, err
	}

	return metalGateway, resp, err
}

func (s *VRFServiceOp) Create(projectID string, input *VRFCreateRequest) (*VRF, *Response, error) {
	if validateErr := ValidateUUID(projectID); validateErr != nil {
		return nil, nil, validateErr
	}
	apiPath := path.Join(projectBasePath, projectID, vrfBasePath)
	output := new(VRF)

	resp, err := s.client.DoRequest("POST", apiPath, input, output)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, nil
}

func (s *VRFServiceOp) Delete(vrfID string) (*Response, error) {
	if validateErr := ValidateUUID(vrfID); validateErr != nil {
		return nil, validateErr
	}
	apiPath := path.Join(vrfBasePath, vrfID)

	resp, err := s.client.DoRequest("DELETE", apiPath, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
