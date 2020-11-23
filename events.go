package packngo

import (
	"path"
)

const eventBasePath = "/events"

// Event struct
type Event struct {
	ID            string     `json:"id,omitempty"`
	State         string     `json:"state,omitempty"`
	Type          string     `json:"type,omitempty"`
	Body          string     `json:"body,omitempty"`
	Relationships []Href     `json:"relationships,omitempty"`
	Interpolated  string     `json:"interpolated,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	Href          string     `json:"href,omitempty"`
}

type eventsRoot struct {
	Events []Event `json:"events,omitempty"`
	Meta   meta    `json:"meta,omitempty"`
}

// EventService interface defines available event functions
type EventService interface {
	List(*ListOptions) ([]Event, *Response, error)
	Get(string, *GetOptions) (*Event, *Response, error)
}

// EventServiceOp implements EventService
type EventServiceOp struct {
	client *Client
}

// List returns all events
func (s *EventServiceOp) List(listOpt *ListOptions) ([]Event, *Response, error) {
	return listEvents(s.client, eventBasePath, listOpt)
}

// Get returns an event by ID
func (s *EventServiceOp) Get(eventID string, getOpt *GetOptions) (*Event, *Response, error) {
	path := path.Join(eventBasePath, eventID)
	return get(s.client, path, getOpt)
}

// list helper function for all event functions
func listEvents(client requestDoer, endpointPath string, opts *ListOptions) (events []Event, resp *Response, err error) {
	path := opts.WithQuery(endpointPath)

	for {
		subset := new(eventsRoot)

		resp, err = client.DoRequest("GET", path, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		events = append(events, subset.Events...)

		if path = nextPage(subset.Meta, opts); path != "" {
			continue
		}
		return
	}

}

func get(client *Client, endpointPath string, opts *GetOptions) (*Event, *Response, error) {
	event := new(Event)

	path := opts.WithQuery(endpointPath)

	resp, err := client.DoRequest("GET", path, nil, event)
	if err != nil {
		return nil, resp, err
	}

	return event, resp, err
}
