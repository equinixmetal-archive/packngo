package packngo

type Rate struct {
	RequestLimit int `json:"request_limit"`
	RequestsRemaining int `json:"requests_remaining"`
	Reset Timestamp `json:"rate_reset"`
}
func (r Rate) String() string {
	return Stringify(r)
}
