package packngo

import "fmt"

const deviceBasePath = "/devices"

type DeviceService interface {
	List(ProjectId string) ([]Device, *Response, error)
	Get(string) (*Device, *Response, error)
	Create(*DeviceCreateRequest) (*Device, *Response, error)
	Delete(string) (*Response, error)
	Reboot(string) (*Response, error)
	PowerOff(string) (*Response, error)
	PowerOn(string) (*Response, error)
}

type DevicesRoot struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
  Href         string    `json:"href,omitempty"`
	Hostname     string    `json:"hostname,omitempty"`
	State        string    `json:"state,omitempty"`
	Created      string    `json:"created_at,omitempty"`
	Updated      string    `json:"updated_at,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	BillingCycle string    `json:"billing_cycle,omitempty"`
	Network      []*IP     `json:"ip_addresses"`
	OS           *OS       `json:"operating_system,omitempty"`
	Plan         *Plan     `json:"plan,omitempty"`
	Facility     *Facility `json:"facility,omitempty"`
	Project      *Project  `json:"project,omitempty"`
  ProvisionPer float32   `json:"provisioning_percentage,omitempty"`
}
func (d Device) String() string {
	return Stringify(d)
}

type DeviceCreateRequest struct {
	Name         string   `json:"name"`
	Plan         string   `json:"plan"`
	Facility     string   `json:"facility"`
	OS           string   `json:"operating_system"`
	BillingCycle string   `json:"billing_cycle"`
	ProjectId    string   `json:"project_id"`
	UserData     string   `json:"user_data"`
	Tags         []string `json:"tags"`
}
func (d DeviceCreateRequest) String() string {
	return Stringify(d)
}

type DeviceActionRequest struct {
	Type string `json:"type"`
}
func (d DeviceActionRequest) String() string {
	return Stringify(d)
}

type IP struct {
	Family  int    `json:"address_family"`
	Cidr    int    `json:"cidr"`
	Address string `json:"address"`
	Gateway string `json:"gateway"`
	Public  bool   `json:"public"`
}
func (n IP) String() string {
	return Stringify(n)
}

type DeviceServiceOp struct {
	client *Client
}

func (s *DeviceServiceOp) List(projectId string) ([]Device, *Response, error) {
	path := fmt.Sprintf("%s/%s/devices", projectBasePath, projectId)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(DevicesRoot)
	resp, err := s.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Devices, resp, err
}

func (s *DeviceServiceOp) Get(deviceId string) (*Device, *Response, error) {
	path := fmt.Sprintf("%s/%s", deviceBasePath, deviceId)

	req, err := s.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	device := new(Device)
	resp, err := s.client.Do(req, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, err
}

func (s *DeviceServiceOp) Create(createRequest *DeviceCreateRequest) (*Device, *Response, error) {
	path := fmt.Sprintf("%s/%s/devices", projectBasePath, createRequest.ProjectId)

	req, err := s.client.NewRequest("POST", path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	device := new(Device)
	resp, err := s.client.Do(req, device)
	if err != nil {
		return nil, resp, err
	}

	return device, resp, err
}

func (s *DeviceServiceOp) Delete(deviceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", deviceBasePath, deviceID)

	req, err := s.client.NewRequest("DELETE", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *DeviceServiceOp) Reboot(deviceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s/actions", deviceBasePath, deviceID)

	action := &DeviceActionRequest { Type: "reboot" }
	req, err := s.client.NewRequest("POST", path, action)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *DeviceServiceOp) PowerOff(deviceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s/actions", deviceBasePath, deviceID)

	action := &DeviceActionRequest { Type: "power_off" }
	req, err := s.client.NewRequest("POST", path, action)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

func (s *DeviceServiceOp) PowerOn(deviceID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s/actions", deviceBasePath, deviceID)

	action := &DeviceActionRequest { Type: "power_on" }
	req, err := s.client.NewRequest("POST", path, action)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}
