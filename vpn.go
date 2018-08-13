package packngo

import "fmt"

const vpnBasePath = "/user/vpn"

// VPNConfig struct
type VPNConfig struct {
	Config string `json:"config,omitempty"`
}

// VPNService interface defines available VPN functions
type VPNService interface {
	Enable() (*Response, error)
	Disable() (*Response, error)
	Get(code string) (*VPNConfig, *Response, error)
}

// VPNServiceOp implements VPNService
type VPNServiceOp struct {
	client *Client
}

// Enable VPN for current user
func (s *VPNServiceOp) Enable() (resp *Response, err error) {
	return s.client.DoRequest("POST", vpnBasePath, nil, nil)
}

// Disable VPN for current user
func (s *VPNServiceOp) Disable() (resp *Response, err error) {
	return s.client.DoRequest("DELETE", vpnBasePath, nil, nil)

}

// Get returns the client vpn config for the currently logged-in user.
func (s *VPNServiceOp) Get(code string) (config *VPNConfig, resp *Response, err error) {
	config = &VPNConfig{}
	path := fmt.Sprintf("%s?code=%s", vpnBasePath, code)

	resp, err = s.client.DoRequest("GET", path, nil, config)
	if err != nil {
		return nil, resp, err
	}

	return config, resp, err
}
