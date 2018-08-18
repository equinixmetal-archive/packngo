package packngo

import "fmt"

const notificationBasePath = "/notifications"

// Notification struct
type Notification struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Body      string `json:"body,omitempty"`
	Severity  string `json:"severity,omitempty"`
	Read      bool   `json:"read,omitempty"`
	Context   string `json:"context,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	User      Href   `json:"user,omitempty"`
	Href      string `json:"href,omitempty"`
}

type notificationsRoot struct {
	Notifications []Notification `json:"notifications,omitempty"`
	Meta          meta           `json:"meta,omitempty"`
}

// NotificationService interface defines available event functions
type NotificationService interface {
	List(*ListOptions) ([]Notification, *Response, error)
}

// NotificationServiceOp implements NotificationService
type NotificationServiceOp struct {
	client *Client
}

// List returns all notifications
func (s *NotificationServiceOp) List(listOpt *ListOptions) ([]Notification, *Response, error) {
	return listNotifications(s.client, notificationBasePath, listOpt)
}

// list helper function for all notification functions
func listNotifications(client *Client, path string, listOpt *ListOptions) ([]Notification, *Response, error) {
	var params string
	if listOpt != nil {
		params = listOpt.createURL()
	}

	root := new(notificationsRoot)

	path = fmt.Sprintf("%s?%s", path, params)

	resp, err := client.DoRequest("GET", path, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Notifications, resp, err
}
