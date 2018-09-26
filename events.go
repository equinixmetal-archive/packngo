package packngo

import "fmt"

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
	Get(string, *ListOptions) (*Event, *Response, error)
}

// EventServiceOp implements EventService
type EventServiceOp struct {
	client *Client
}

// List returns all events
func (s *EventServiceOp) List(listOpt *ListOptions) ([]Event, *Response, error) {
	return list(s.client, eventBasePath, listOpt)
}

// Get returns an event by ID
func (s *EventServiceOp) Get(eventID string, listOpt *ListOptions) (*Event, *Response, error) {
	path := fmt.Sprintf("%s/%s", eventBasePath, eventID)
	return get(s.client, path, listOpt)
}

// list helper function for all event functions
func list(client *Client, path string, listOpt *ListOptions) ([]Event, *Response, error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}

	root := new(eventsRoot)

	path = fmt.Sprintf("%s?%s", path, params)

	resp, err := client.DoRequest("GET", path, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Events, resp, err
}

func get(client *Client, path string, listOpt *ListOptions) (*Event, *Response, error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}

	event := new(Event)

	path = fmt.Sprintf("%s?%s", path, params)

	resp, err := client.DoRequest("GET", path, nil, event)
	if err != nil {
		return nil, resp, err
	}

	return event, resp, err
}
