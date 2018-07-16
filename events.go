package packngo

// Event struct
type Event struct {
	ID            string    `json:"id,omitempty"`
	State         string    `json:"state,omitempty"`
	Type          string    `json:"type,omitempty"`
	Body          string    `json:"body,omitempty"`
	Relationships []Href    `json:"relationships,omitempty"`
	Interpolated  string    `json:"interpolated,omitempty"`
	CreatedAt     Timestamp `json:"created_at,omitempty"`
	Href          string    `json:"href,omitempty"`
}

type eventsRoot struct {
	Events []Event `json:"events,omitempty"`
	Meta   meta    `json:"meta,omitempty"`
}
