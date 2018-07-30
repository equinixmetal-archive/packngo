package packngo

const twoFactorAuthAppPath = "/user/otp/app"
const twoFactorAuthSmsPath = "/user/otp/sms"

// TwoFactorAuthService interface defines available two factor authentication functions
type TwoFactorAuthService interface {
	EnableApp() (*Response, error)
	DisableApp() (*Response, error)
	EnableSms() (*Response, error)
	DisableSms() (*Response, error)
}

// TwoFactorAuthServiceOp implements TwoFactorAuthService
type TwoFactorAuthServiceOp struct {
	client *Client
}

// EnableApp function enables two factor auth using authenticatior app
func (s *TwoFactorAuthServiceOp) EnableApp() (resp *Response, err error) {
	return s.client.DoRequest("POST", twoFactorAuthAppPath, nil, nil)
}

// EnableSms function enables two factor auth using sms
func (s *TwoFactorAuthServiceOp) EnableSms() (resp *Response, err error) {
	return s.client.DoRequest("POST", twoFactorAuthSmsPath, nil, nil)
}

// DisableApp function disables two factor auth using
func (s *TwoFactorAuthServiceOp) DisableApp() (resp *Response, err error) {
	return s.client.DoRequest("DELETE", twoFactorAuthAppPath, nil, nil)
}

// DisableSms function disables two factor auth using
func (s *TwoFactorAuthServiceOp) DisableSms() (resp *Response, err error) {
	return s.client.DoRequest("DELETE", twoFactorAuthSmsPath, nil, nil)
}
