package packngo

import "fmt"

const staffBasePath = "/staff"
const hardwareBasePath = "hardware"

// HardwareService interface defines available hardware device methods
type HardwareService interface {
	List(listOpt *ListOptions) ([]Hardware, *Response, error)
}

// Hardware represents a Packet hardware from the API
type Hardware struct {
	ID                   string                `json:"id"`
	Href                 string                `json:"href,omitempty"`
	Hostname             string                `json:"hostname,omitempty"`
	ModelNumber          string                `json:"model_number,omitempty"`
	State                string                `json:"state,omitempty"`
	Created              string                `json:"created_at,omitempty"`
	Updated              string                `json:"updated_at,omitempty"`
	HardwareManufacturer *HardwareManufacturer `json:"manufacturer,omitempty"`
}

type HardwareManufacturer struct {
	ID      string `json:"id"`
	Created string `json:"created_at,omitempty"`
	Updated string `json:"updated_at,omitempty"`
	Slug    string `json:"slug"`
}

type hardwareRoot struct {
	Hardware []Hardware `json:"hardware"`
	Meta     meta       `json:"meta"`
}

// HardwareServiceOp implements DeviceService
type HardwareServiceOp struct {
	client *Client
}

// List returns hardware devices
func (s *HardwareServiceOp) List(listOpt *ListOptions) (hardware []Hardware, resp *Response, err error) {

	params := urlQuery(listOpt)
	path := fmt.Sprintf("%s/%s?%s", staffBasePath, hardwareBasePath, params)
	for {
		subset := new(hardwareRoot)

		staffHeader := map[string]string{"X-Packet-Staff": "true"}
		resp, err = s.client.DoRequestWithHeader("GET", staffHeader, path, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		hardware = append(hardware, subset.Hardware...)

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
