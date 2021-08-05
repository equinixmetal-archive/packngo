package packngo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
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

func TestAccVLANAssignmentServiceOp(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()

	metro := testMetro()

	cr := DeviceCreateRequest{
		Hostname:     hn,
		Metro:        metro,
		Plan:         testPlan(),
		OS:           testOS,
		ProjectID:    projectID,
		BillingCycle: "hourly",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteDevice(t, c, d.ID, false)
	dID := d.ID

	d = waitDeviceActive(t, c, dID)
	bond0, err := d.GetPortByName("bond0")

	if err != nil {
		t.Fatal(err)
	}

	if len(bond0.AttachedVirtualNetworks) != 0 {
		t.Fatal("No vlans should be attached to a eth1 in the beginning of this test")
	}

	vlans := new([2]*VirtualNetwork)
	for i, vxlan := range []int{1222, 1234} {
		vncr := VirtualNetworkCreateRequest{
			ProjectID: projectID,
			Metro:     metro,
			VXLAN:     vxlan,
		}

		vlans[i], _, err = c.ProjectVirtualNetworks.Create(&vncr)
		if err != nil {
			t.Fatal(err)
		}

		defer deleteProjectVirtualNetwork(t, c, vlans[i].ID)
	}

	// test unassignment and assignment states with both supported formats, VLAN ID and VXLAN
	vabcr := &VLANAssignmentBatchCreateRequest{
		VLANAssignments: []VLANAssignmentCreateRequest{{
			VLAN:  vlans[0].ID,
			State: VLANAssignmentUnassigned,
		}, {
			VLAN:   strconv.Itoa(vlans[1].VXLAN),
			State:  VLANAssignmentAssigned,
			Native: func() *bool { b := false; return &b }(),
		}},
	}

	b, _, err := c.VLANAssignments.CreateBatch(bond0.ID, vabcr, nil)
	if err != nil {
		t.Fatal(err)
	}

	// We must unassign the port when done testing because we can not delete the VLAN when in use
	defer func() {
		_, _, err := c.Ports.Unassign(bond0.ID, strconv.Itoa(vlans[1].VXLAN))
		if err != nil {
			t.Log(err)
		}
	}()

	if b.Quantity != 2 {
		t.Fatal("Exactly two attachments should be batched")
	}

	if path.Base(b.Port.Href.Href) != bond0.ID {
		t.Fatal("mismatch in the UUID of the assigned Port")
	}

	// verify batch can be fetched
	b2 := waitVLANAssignmentBatch(t, c, bond0.ID, b.ID)
	if b2.CreatedAt != b.CreatedAt {
		t.Fatal("Reloaded VLANAssignment batch create time should match original response")
	}

	// verify port vlan assignments can be listed
	allAs, _, err := c.VLANAssignments.List(bond0.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We do not expect "unassigned" vlan assignments to persist
	if len(allAs) != 1 {
		t.Fatal("unexpected or missing VLANAssignments in List results")
	}

	// verify single assignment can be fetched
	a2, _, err := c.VLANAssignments.Get(bond0.ID, allAs[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if a2.CreatedAt != allAs[0].CreatedAt {
		t.Fatal("Reloaded VLANAssignment create time should match original response")
	}

	// verify port vlan batch assignments can be listed
	allBs, _, err := c.VLANAssignments.ListBatch(bond0.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	unseen := map[string]bool{
		b2.ID: true,
	}
	for _, b := range allBs {
		delete(unseen, b.ID)
	}
	if len(unseen) != 0 {
		t.Fatal("unexpected or missing Batch Assignments in ListBatch results")
	}
}

func waitVLANAssignmentBatch(t *testing.T, c *Client, portID, id string) *VLANAssignmentBatch {
	// 15 minutes = 180 * 5sec-retry
	for i := 0; i < 180; i++ {
		<-time.After(5 * time.Second)
		b, _, err := c.VLANAssignments.GetBatch(portID, id, nil)
		if err != nil {
			t.Fatal(err)
			return nil
		}
		if b.State == VLANAssignmentBatchCompleted {
			return b
		}
		if b.State == VLANAssignmentBatchFailed {
			t.Fatalf("vlan assignment batch %s provisioning failed: %s", id, strings.Join(b.ErrorMessages, "; "))
			return nil
		}
	}

	t.Fatal(fmt.Errorf("vlan assignment batch %s is still not complete after timeout", id))
	return nil
}
