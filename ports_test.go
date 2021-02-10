package packngo

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

var (
	errBoom              = errors.New("boom")
	testPortID           = "12345"
	testVirtualNetworkID = "abcdef"
)

// reverse a string slice
// https://github.com/golang/go/wiki/SliceTricks#reversing
func reverse(a []string) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

// testPort returns a non-empty Port. Few parameters are configurable because
// the functions relying on this Port do not care about the contents. Any
// non-empty Port would be sufficient.
func testPort(id string) *Port {
	v := &Port{}
	v.ID = id
	v.Type = "NetworkPort"
	v.Name = "eth0"
	v.Data = PortData{
		MAC:    "aa:bb:cc:dd:ee:ff",
		Bonded: true,
	}
	v.NetworkType = NetworkTypeL3
	v.Bond = &BondData{
		ID:   "bond0-uuid",
		Name: "bond0",
	}
	v.AttachedVirtualNetworks = []VirtualNetwork{
		{
			ID:           "abcedf",
			Description:  "vlan-foo",
			VXLAN:        1234,
			FacilityCode: "facility-foo",
			CreatedAt:    "2020-07-02T15:24:28Z",
			Href:         "/virtual-networks/f343a677-c86a-48f4-9cc3-607afedc1ca2",
		},
	}
	return v
}

func TestAccPortServiceOp_Get(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()
	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	fac := testFacility()
	plan := testPlan()

	cr := DeviceCreateRequest{
		Hostname:     "test-portget",
		Facility:     []string{fac},
		Plan:         plan,
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}
	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	deviceID := d.ID

	d = waitDeviceActive(t, c, deviceID)
	defer deleteDevice(t, c, d.ID, false)

	for _, p := range d.NetworkPorts {
		port, resp, err := c.Ports.Get(p.ID, nil)
		if err != nil {
			t.Fatal(err)
		}

		if resp == nil || resp.StatusCode != 200 {
			t.Fatal("Expected a HTTP response with a 200 StatusCode")
		}

		if port.ID != p.ID {
			t.Fatal("Fetched port is not expected port")
		}
	}
}

func TestPortServiceOp_Assign(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
		vlanID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "assign" && method == http.MethodPost {
							*v = *(testPort(parts[1]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
				vlanID: testVirtualNetworkID,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
				vlanID: testVirtualNetworkID,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.Assign(tt.args.portID, tt.args.vlanID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.Assign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.Assign() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.Assign() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_AssignNative(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
		vlanID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "native-vlan" && method == http.MethodPost {
							*v = *(testPort(parts[1]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
				vlanID: testVirtualNetworkID,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
				vlanID: testVirtualNetworkID,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.AssignNative(tt.args.portID, tt.args.vlanID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.AssignNative() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.AssignNative() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.AssignNative() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_Unassign(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
		vlanID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "unassign" && method == http.MethodPost {
							*v = *(testPort(parts[1]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
				vlanID: testVirtualNetworkID,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
				vlanID: testVirtualNetworkID,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.Unassign(tt.args.portID, tt.args.vlanID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.Unassign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.Unassign() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.Unassign() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_UnassignNative(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "native-vlan" && method == http.MethodDelete {
							*v = *(testPort(parts[1]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.UnassignNative(tt.args.portID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.UnassignNative() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.UnassignNative() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.UnassignNative() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_Bond(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID     string
		bulkEnable bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "bond" && method == http.MethodPost {
							*v = *(testPort(parts[1]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID:     testPortID,
				bulkEnable: false,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID:     testPortID,
				bulkEnable: false,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.Bond(tt.args.portID, tt.args.bulkEnable)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.Bond() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.Bond() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.Bond() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_Disbond(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID     string
		bulkEnable bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "disbond" && method == http.MethodPost {
							*v = *(testPort(parts[1]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID:     testPortID,
				bulkEnable: false,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID:     testPortID,
				bulkEnable: false,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.Disbond(tt.args.portID, tt.args.bulkEnable)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.Disbond() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.Disbond() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.Disbond() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_ConvertToLayerTwo(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "layer-2" && parts[1] == "convert" && method == http.MethodPost {
							*v = *(testPort(parts[2]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.ConvertToLayerTwo(tt.args.portID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.ConvertToLayerTwo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.ConvertToLayerTwo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.ConvertToLayerTwo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_ConvertToLayerThree(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
		ips    []AddressRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[0] == "layer-3" && parts[1] == "convert" && method == http.MethodPost {
							*v = *(testPort(parts[2]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
				ips:    []AddressRequest{},
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
				ips:    []AddressRequest{},
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := i.ConvertToLayerThree(tt.args.portID, tt.args.ips)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.ConvertToLayerThree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.ConvertToLayerThree() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.ConvertToLayerThree() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPortServiceOp_Get(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID string
		opts   *GetOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Port
		want1   *Response
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						parts := strings.Split(path, "/")
						reverse(parts)

						if v, ok := v.(*Port); ok && parts[1] == "ports" && method == http.MethodGet {
							*v = *(testPort(parts[0]))
							return &Response{}, nil
						}

						return nil, errBoom
					},
				},
			},
			args: args{
				portID: testPortID,
				opts:   nil,
			},
			want:    testPort(testPortID),
			want1:   &Response{},
			wantErr: false,
		},
		{
			name: "ErrorIsHandled",
			fields: fields{client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, fmt.Errorf("boom")
				},
			}},
			args: args{
				portID: testPortID,
				opts:   nil,
			},
			want:    nil,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PortServiceOp{
				client: tt.fields.client,
			}
			got, got1, err := s.Get(tt.args.portID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("PortServiceOp.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PortServiceOp.Get() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PortServiceOp.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
