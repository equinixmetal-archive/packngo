package packngo

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
)

func TestAccListHardwareReservations(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, projectID, tearDown := setupWithProject(t)
	defer tearDown()

	hardwareReservations, _, err := c.HardwareReservations.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(hardwareReservations) != 0 {
		t.Fatal("There should not be any hardware reservations.")
	}
}

func TestHardwareReservationServiceOp_List(t *testing.T) {
	type fields struct {
		client requestDoer
	}
	type args struct {
		projectID string
		opts      *ListOptions
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantReservations []HardwareReservation
		wantResp         *Response
		wantErr          bool
	}{{
		name: "Pagination",
		fields: fields{
			client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					// Return one reservation per page, tests pagination
					data := []HardwareReservation{{ID: "1"}, {ID: "2"}, {ID: "3"}}

					if v, ok := v.(*hardwareReservationRoot); ok && method == http.MethodGet {
						u, _ := url.Parse(path)
						q, _ := url.ParseQuery(u.RawQuery)
						page, _ := strconv.Atoi(q.Get("page"))
						if page == 0 {
							page = 1
						}
						v.HardwareReservations = []HardwareReservation{data[page-1]}
						v.Meta.Total = len(data)
						v.Meta.CurrentPageNum = page
						if page < v.Meta.Total {
							nextPage := page + 1
							nextHref := fmt.Sprintf("%s?page=%d", u.Path, nextPage)
							v.Meta.Next = &Href{Href: &nextHref}
						}
						return &Response{}, nil
					}

					return nil, errBoom
				},
			},
		},
		args:             args{projectID: testProjectId},
		wantReservations: []HardwareReservation{{ID: "1"}, {ID: "2"}, {ID: "3"}},
		wantResp:         &Response{},
		wantErr:          false,
	}, {
		name: "Error",
		fields: fields{
			client: &MockClient{
				fnDoRequest: func(method, path string, body, v interface{}) (*Response, error) {
					return nil, errBoom
				},
			},
		},
		args:             args{projectID: testProjectId},
		wantReservations: nil,
		wantResp:         nil,
		wantErr:          true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &HardwareReservationServiceOp{
				client: tt.fields.client,
			}
			gotReservations, gotResp, err := s.List(tt.args.projectID, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("HardwareReservationServiceOp.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReservations, tt.wantReservations) {
				t.Errorf("HardwareReservationServiceOp.List() gotReservations = %v, want %v", gotReservations, tt.wantReservations)
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("HardwareReservationServiceOp.List() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
