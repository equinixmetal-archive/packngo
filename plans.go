package packngo 

const planBasePath = "/plans"

type PlanService interface {
	List() ([]Plan, *Response, error)
}

type planRoot struct {
	Plans []Plan `json:"plans"`
}

type Plan struct {
	Id        string    `json:"id"`
	Slug      string    `json:"slug,omitempty"`
	Name      string    `json:"name,omitempty"`
	Line      string    `json:"line,omitempty"`
	Specs     *Specs    `json:"specs,omitempty"`
	Pricing   *Pricing  `json:"pricing,omitempty"`
}
func (p Plan) String() string {
	return Stringify(p)
}

type Specs struct {
	Cpus      *Cpus     `json:"cpus,omitempty"`
	Memory    *Memory   `json:"memory,omitempty"`
	Drives    *Drives   `json:"drives,omitempty"`
  Nics      *Nics     `json:"nics,omitempty"`
	Features  *Features `json:"features,omitempty"`
}
func (s Specs) String() string {
	return Stringify(s)
}

type Cpus struct {
	Count int    `json:"count,omitempty"`
	Type  string `json:"type,omitempty"`
}
func (c Cpus) String() string {
	return Stringify(c)
}

type Memory struct {
	Total string `json:"total,omitempty"`
}
func (m Memory) String() string {
	return Stringify(m)
}

type Drives struct {
	Count int    `json:"count,omitempty"`
	Size  string `json:"size,omitempty"`
	Type  string `json:"type,omitempty"`
}
func (d Drives) String() string {
	return Stringify(d)
}

type Nics struct {
	Count int    `json:"count,omitempty"`
	Type  string `json:"type,omitempty"`
}
func (n Nics) String() string {
	return Stringify(n)
}

type Features struct {
	Raid bool `json:"raid,omitempty"`
	Txt  bool `json:"txt,omitempty"`
}
func (f Features) String() string {
	return Stringify(f)
}

type Pricing struct {
	Hourly  float32 `json:"hourly,omitempty"`
	Monthly float32 `json:"monthly,omitempty"`
}
func (p Pricing) String() string {
	return Stringify(p)
}

type PlanServiceOp struct {
	client *Client
}

func (s *PlanServiceOp) List() ([]Plan, *Response, error) {
	path := "plans"

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(planRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Plans, resp, err
}
