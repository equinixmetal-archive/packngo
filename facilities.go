package packngo 

const facilityBasePath = "/facilities"

type FacilityService interface {
	List() ([]Facility, *Response, error)
}

type facilityRoot struct {
	Facilities []Facility `json:"facilities"`
}

type Facility struct {
	Id        string   `json:"id"`
	Name      string   `json:"name,omitempty"`
	Code      string   `json:"code,omitempty"`
	Features  []string `json:"features,omitempty"`
	Address   *Address `json:"address,omitempty"`
  Url       string   `json:"href,omitempty"`
}
func (f Facility) String() string {
	return Stringify(f)
}

type Address struct {
	Id string `json:"id,omitempty"`
}
func (a Address) String() string {
	return Stringify(a)
}

type FacilityServiceOp struct {
	client *Client
}

func (s *FacilityServiceOp) List() ([]Facility, *Response, error) {
	req, err := s.client.NewRequest("GET", facilityBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(facilityRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Facilities, resp, err
}
