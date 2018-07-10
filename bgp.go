package packngo

import "fmt"

var bgpBasePath = "/bgp/sessions"

// BGPService interface defines available device methods
type BGPService interface {
	List(listOpt *ListOptions) ([]BgpSession, *Response, error)
	Get(string, *ListOptions) (*BgpSession, *Response, error)
	Create(string, CreateBGPSessionRequest) (*BgpSession, *Response, error)
	Delete(string) (*Response, error)
}

// BGPServiceOp implements DeviceService
type BGPServiceOp struct {
	client *Client
}

// BgpSession represents a Packet BGP Session
type BgpSession struct {
	ID            string   `json:"id,omitempty"`
	Status        string   `json:"status,omitempty"`
	LearnedRoutes []string `json:"learned_routes,omitempty"`
	AddressFamily string   `json:"address_family,omitempty"`
	Device        Device   `json:"device,omitempty"`
	Href          string   `json:"href,omitempty"`
}

// BgpConfig represents a Packet BGP Config
type BgpConfig struct {
	ID             string      `json:"id,omitempty"`
	Status         string      `json:"status,omitempty"`
	DeploymentType string      `json:"deployment_type,omitempty"`
	Asn            int32       `json:"asn,omitempty"`
	RouteObject    string      `json:"route_object,omitempty"`
	Md5            string      `json:"md5,omitempty"`
	MaxPrefix      int32       `json:"max_prefix,omitempty"`
	Project        Project     `json:"project,omitempty"`
	CreatedAt      Timestamp   `json:"created_at,omitempty"`
	RequestedAt    Timestamp   `json:"requested_at,omitempty"`
	Session        interface{} `json:"session,omitempty"`
	Href           string      `json:"href,omitempty"`
}

// CreateBGPSessionRequest struct
type CreateBGPSessionRequest struct {
	AddressFamily string `json:"address_family"`
}

// Create function
func (s *BGPServiceOp) Create(deviceID string, request CreateBGPSessionRequest) (*BgpSession, *Response, error) {
	path := fmt.Sprintf("%s/%s/%s", deviceBasePath, deviceID, bgpBasePath)
	session := new(BgpSession)

	resp, err := s.client.DoRequest("POST", path, request, session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, err
}

// Delete function
func (s *BGPServiceOp) Delete(id string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", bgpBasePath, id)

	return s.client.DoRequest("DELETE", path, nil, nil)
}
