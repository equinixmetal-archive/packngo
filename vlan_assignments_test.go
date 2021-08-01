package packngo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func mockAssignedPortBody(assignmentid, portid, vnid string) string {
	return fmt.Sprintf(`{
		"id":"%s",
		"created_at":"2021-05-28T16:02:33Z",
		"updated_at":"2021-05-28T16:02:33Z",
		"native":false,
		"virtual_network":
		   {"href":"/virtual-networks/%s"},
		"port":{"href":"/ports/%s"},
		"vlan":1234,
		"state":"assigned"
	}`, assignmentid, vnid, portid)
}

func TestVLANAssignmentServiceOp_Get(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		portID       string
		assignmentID string
		opts         *GetOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *VLANAssignment
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				client: (func() *MockClient {

					raw := mockAssignedPortBody("bar", "foo", "vlan")
					mockNR := mockNewRequest()
					mockDo := func(req *http.Request, obj interface{}) (*Response, error) {
						// baseURL is not needed here
						expectedPath := path.Join(portBasePath, "foo", portVLANAssignmentsPath, "bar")
						if expectedPath != req.URL.Path {
							return nil, fmt.Errorf("wrong url")
						}
						if err := json.NewDecoder(strings.NewReader(raw)).Decode(obj); err != nil {
							return nil, err
						}

						return mockResponse(200, raw, req), nil
					}

					return &MockClient{
						fnDoRequest: mockDoRequest(mockNR, mockDo),
					}
				})(),
			},
			args: args{portID: "foo", assignmentID: "bar"},
			want: &VLANAssignment{
				ID:             "bar",
				CreatedAt:      Timestamp{Time: func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-05-28T16:02:33Z"); return t }()},
				UpdatedAt:      Timestamp{Time: func() time.Time { t, _ := time.Parse(time.RFC3339, "2021-05-28T16:02:33Z"); return t }()},
				VirtualNetwork: &VirtualNetwork{Href: "/virtual-networks/vlan"},
				Port:           &Port{Href: &Href{Href: "/ports/foo"}},
				VLAN:           1234,
				State:          VLANAssignmentAssigned,
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				client: &MockClient{
					fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
						return nil, errBoom
					},
				},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &VLANAssignmentServiceOp{
				client: tt.fields.client,
			}
			got, _, err := s.Get(tt.args.portID, tt.args.assignmentID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("VLANAssignmentServiceOp.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("VLANAssignmentServiceOp.Get() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
