package packngo

const userBasePath = "/users"

type UserService interface {
	Get(string) (*User, *Response, error)
}

type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name,omitempty"`
  LastName     string    `json:"last_name,omitempty"`
  FullName     string    `json:"full_name,omitempty"`
  Email        string    `json:"email,omitempty"`
  TwoFactor    string    `json:"two_factor_auth,omitempty"`
	AvatarUrl    string    `json:"avatar_url,omitempty"`
	Facebook     string    `json:"twitter,omitempty"`
	Twitter      string    `json:"facebook,omitempty"`
	LinkedIn     string    `json:"linkedin,omitempty"`
	Created      string    `json:"created_at,omitempty"`
	Updated      string    `json:"updated_at,omitempty"`
	TimeZone     string    `json:"timezone,omitempty"`
	Emails       []Email   `json:"email,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	Url          string    `json:"href,omitempty"`
}
func (u User) String() string {
	return Stringify(u)
}

type UserServiceOp struct {
	client *Client
}

func (s *UserServiceOp) Get(userID string) (*User, *Response, error) {
	req, err := s.client.NewRequest("GET", userBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}
