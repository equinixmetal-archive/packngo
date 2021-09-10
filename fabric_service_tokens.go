package packngo

import "path"

type FabricServiceTokenType string

const (
	fabricServiceTokenBasePath                        = "/fabric-service-tokens"
	FabricServiceTokenASide    FabricServiceTokenType = "a_side"
	FabricServiceTokenZSide    FabricServiceTokenType = "z_side"
)

// FabricServiceTokenService interface defines available metro methods
type FabricServiceTokenService interface {
	Get(string, *GetOptions) (*FabricServiceToken, *Response, error)
}

// FabricServiceToken represents an Equinix Metal metro
type FabricServiceToken struct {
	ID               string                 `json:"id"`
	Role             string                 `json:"role,omitempty"`
	State            string                 `json:"state,omitempty"`
	MaxAllowedSpeed  int                    `json:"max_allowed_speed,omitempty"`
	ServiceTokenType FabricServiceTokenType `json:"service_token_type,omitempty"`
	Connection       *Connection            `json:"interconnection,omitempty"`
	ConnectionPort   *ConnectionPort        `json:"interconnection_port,omitempty"`
	Organization     *Organization          `json:"organization,omitempty"`
}

func (f FabricServiceToken) String() string {
	return Stringify(f)
}

// FabricServiceTokenServiceOp implements FabricServiceTokenService
type FabricServiceTokenServiceOp struct {
	client *Client
}

func (s *FabricServiceTokenServiceOp) Get(id string, opts *GetOptions) (*FabricServiceToken, *Response, error) {
	endpointPath := path.Join(fabricServiceTokenBasePath, id)
	apiPathQuery := opts.WithQuery(endpointPath)
	fst := new(FabricServiceToken)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, fst)
	if err != nil {
		return nil, resp, err
	}
	return fst, resp, err
}
