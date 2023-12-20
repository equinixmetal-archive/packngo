package packngo

import (
	"net/http"

	"github.com/packethost/packngo/href"
)

type includer interface {
	DefaultIncludes() []string
}

type clienter interface {
	GetClient() requestDoer
}

type serviceOper interface {
	includer
	clienter
}
type serviceOp struct {
	serviceOper

	client *Client
}

// Hydrate fetches and populates a resource based on the resources Href.
func (s *serviceOp) Hydrate(resource href.Hrefer, opts *GetOptions) (*Response, error) {
	x := s.DefaultIncludes()
	opts = opts.Including(x...)

	apiPathQuery := opts.WithQuery(resource.GetHref())

	return s.GetClient().DoRequest(http.MethodGet, apiPathQuery, nil, resource)
}
