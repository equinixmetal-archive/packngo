package packngo

import "fmt"

var bgpConfigBasePath = "/bgp-config"

// BGPConfigService interface defines available BGP config methods
type BGPConfigService interface {
	Get(string, *ListOptions) (*BGPConfig, *Response, error)
	Create(string, CreateBGPSessionRequest) (*BGPConfig, *Response, error)
}

// BGPConfigServiceOp implements BgpConfigService
type BGPConfigServiceOp struct {
	client *Client
}

// CreateBGPConfigRequest struct
type CreateBGPConfigRequest struct {
	DeploymentType string `json:"deployment_type,omitempty"`
	Asn            int    `json:"asn,omitempty"`
	Md5            string `json:"md5,omitempty"`
	UseCase        string `json:"use_case,omitempty"`
}

// BGPConfig represents a Packet BGP Config
type BGPConfig struct {
	ID             string     `json:"id,omitempty"`
	Status         string     `json:"status,omitempty"`
	DeploymentType string     `json:"deployment_type,omitempty"`
	Asn            int32      `json:"asn,omitempty"`
	RouteObject    string     `json:"route_object,omitempty"`
	Md5            string     `json:"md5,omitempty"`
	MaxPrefix      int32      `json:"max_prefix,omitempty"`
	Project        Project    `json:"project,omitempty"`
	CreatedAt      Timestamp  `json:"created_at,omitempty"`
	RequestedAt    Timestamp  `json:"requested_at,omitempty"`
	Session        BGPSession `json:"session,omitempty"`
	Href           string     `json:"href,omitempty"`
}

// Create function
func (s *BGPConfigServiceOp) Create(projectID string, request CreateBGPConfigRequest) (*BGPConfig, *Response, error) {
	path := fmt.Sprintf("%s/%s/%ss", projectBasePath, projectID, bgpConfigBasePath)
	session := new(BGPConfig)

	resp, err := s.client.DoRequest("POST", path, request, session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, err
}

// Get function
func (s *BGPConfigServiceOp) Get(projectID string, listOpt *ListOptions) (bgpConfig *BGPConfig, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s/%s?%s", projectBasePath, projectID, bgpConfigBasePath, params)

	subset := new(BGPConfig)

	resp, err = s.client.DoRequest("GET", path, nil, subset)
	if err != nil {
		return nil, resp, err
	}

	return subset, resp, err
}
