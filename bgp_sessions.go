package packngo

import "fmt"

var bgpSessionBasePath = "/bgp/sessions"

// BGPSessionService interface defines available BGP session methods
type BGPSessionService interface {
	ListByDevice(string, listOpt *ListOptions) ([]BGPSession, *Response, error)
	ListByProject(string, listOpt *ListOptions) ([]BGPSession, *Response, error)
	Get(string, *ListOptions) (*BGPSession, *Response, error)
	Create(string, CreateBGPSessionRequest) (*BGPSession, *Response, error)
	Delete(string) (*Response, error)
}

type bgpSessionsRoot struct {
	Sessions []BGPSession `json:"bgp_sessions"`
	Meta     meta         `json:"meta"`
}

// BGPSessionServiceOp implements BgpSessionService
type BGPSessionServiceOp struct {
	client *Client
}

// BgpSession represents a Packet BGP Session
type BGPSession struct {
	ID            string   `json:"id,omitempty"`
	Status        string   `json:"status,omitempty"`
	LearnedRoutes []string `json:"learned_routes,omitempty"`
	AddressFamily string   `json:"address_family,omitempty"`
	Device        Device   `json:"device,omitempty"`
	Href          string   `json:"href,omitempty"`
}

// CreateBGPSessionRequest struct
type CreateBGPSessionRequest struct {
	AddressFamily string `json:"address_family"`
}

// Create function
func (s *BGPConfigServiceOp) Create(deviceID string, request CreateBGPSessionRequest) (*BGPSession, *Response, error) {
	path := fmt.Sprintf("%s/%s/%s", deviceBasePath, deviceID, bgpSessionBasePath)
	session := new(BGPSession)

	resp, err := s.client.DoRequest("POST", path, request, session)
	if err != nil {
		return nil, resp, err
	}

	return session, resp, err
}

// Delete function
func (s *BGPSessionServiceOp) Delete(id string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", bgpSessionBasePath, id)

	return s.client.DoRequest("DELETE", path, nil, nil)
}

// ListByDevice function
func (s *BGPSessionServiceOp) ListByDevice(deviceID string, listOpt *ListOptions) (bgpSessions []BGPSession, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s/%s?%s", deviceBasePath, deviceID, bgpSessionBasePath, params)

	for {
		subset := new(bgpSessionsRoot)

		resp, err = s.client.DoRequest("GET", path, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		bgpSessions = append(bgpSessions, subset.Sessions...)

		if subset.Meta.Next != nil && (listOpt == nil || listOpt.Page == 0) {
			path = subset.Meta.Next.Href
			if params != "" {
				path = fmt.Sprintf("%s&%s", path, params)
			}
			continue
		}

		return
	}

}

// ListByProject function
func (s *BGPSessionServiceOp) ListByProject(projectID string, listOpt *ListOptions) (bgpSessions []BGPSession, resp *Response, err error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}
	path := fmt.Sprintf("%s/%s/%s?%s", projectBasePath, projectID, bgpSessionBasePath, params)

	for {
		subset := new(bgpSessionsRoot)

		resp, err = s.client.DoRequest("GET", path, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		bgpSessions = append(bgpSessions, subset.Sessions...)

		if subset.Meta.Next != nil && (listOpt == nil || listOpt.Page == 0) {
			path = subset.Meta.Next.Href
			if params != "" {
				path = fmt.Sprintf("%s&%s", path, params)
			}
			continue
		}

		return
	}

}
