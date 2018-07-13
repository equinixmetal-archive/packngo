package packngo

const capacityBasePath = "/capacity"

// CapacityService interface defines available capacity methods
type CapacityService interface {
	List() (*CapacityReport, *Response, error)
	Check(*CapacityInput) (*Response, error)
}

// CapacityInput struct
type CapacityInput struct {
	Servers []ServerInfo `json:"servers,omitempty"`
}

// ServerInfo struct
type ServerInfo struct {
	Facility string `json:"facility,omitempty"`
	Plan     string `json:"plan,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

type capacityRoot struct {
	Capacity CapacityReport `json:"capacity,omitempty"`
}

// CapacityReport struct
type CapacityReport struct {
	Ams1 *CapacityPerFacility `json:"ams1,omitempty"`
	Ewr1 *CapacityPerFacility `json:"ewr1,omitempty"`
	Sjc1 *CapacityPerFacility `json:"sjc1,omitempty"`
	Atl1 *CapacityPerFacility `json:"atl1,omitempty"`
	Dfw1 *CapacityPerFacility `json:"dfw1,omitempty"`
	Fra1 *CapacityPerFacility `json:"fra1,omitempty"`
	Iad1 *CapacityPerFacility `json:"iad1,omitempty"`
	Lax1 *CapacityPerFacility `json:"lax1,omitempty"`
	Nrt1 *CapacityPerFacility `json:"nrt1,omitempty"`
	Ord1 *CapacityPerFacility `json:"ord1,omitempty"`
	Sea1 *CapacityPerFacility `json:"sea1,omitempty"`
	Sin1 *CapacityPerFacility `json:"sin1,omitempty"`
	Syd1 *CapacityPerFacility `json:"syd1,omitempty"`
	Yyz1 *CapacityPerFacility `json:"yyz1,omitempty"`
}

// CapacityPerFacility struct
type CapacityPerFacility struct {
	T1SmallX86  *CapacityPerBaremetal `json:"t1.small.x86,omitempty"`
	C1SmallX86  *CapacityPerBaremetal `json:"c1.small.x86,omitempty"`
	M1XlargeX86 *CapacityPerBaremetal `json:"m1.xlarge.x86,omitempty"`
	C1XlargeX86 *CapacityPerBaremetal `json:"c1.xlarge.x86,omitempty"`

	Baremetal0   *CapacityPerBaremetal `json:"baremetal_0,omitempty"`
	Baremetal1   *CapacityPerBaremetal `json:"baremetal_1,omitempty"`
	Baremetal1e  *CapacityPerBaremetal `json:"baremetal_1e,omitempty"`
	Baremetal2   *CapacityPerBaremetal `json:"baremetal_2,omitempty"`
	Baremetal2a  *CapacityPerBaremetal `json:"baremetal_2a,omitempty"`
	Baremetal2a2 *CapacityPerBaremetal `json:"baremetal_2a2,omitempty"`
	Baremetal3   *CapacityPerBaremetal `json:"baremetal_3,omitempty"`
}

// CapacityPerBaremetal struct
type CapacityPerBaremetal struct {
	Level string `json:"level,omitempty"`
}

// CapacityList struct
type CapacityList struct {
	Capacity CapacityReport `json:"capacity,omitempty"`
}

// CapacityServiceOp implements CapacityService
type CapacityServiceOp struct {
	client *Client
}

// List returns a list of facilities and plans with their current capacity.
func (s *CapacityServiceOp) List() (*CapacityReport, *Response, error) {
	root := new(capacityRoot)

	resp, err := s.client.DoRequest("GET", capacityBasePath, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Capacity, nil, nil
}

// Check validates if a deploy can be fulfilled.
func (s *CapacityServiceOp) Check(input *CapacityInput) (resp *Response, err error) {

	return s.client.DoRequest("POST", capacityBasePath, input, nil)

}
