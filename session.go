package packngo

import "os"

// Config ...
type Config interface {
	URL() string

	DefaultProject() string
	DefaultOrganization() string
	DefaultMetro() string
	DefaultDevicePlan() string
	DefaultDeviceOS() string

	Token() string
	DebugEnabled() bool
}

// MetalConfig ...
type MetalConfig struct {
	OrganizationID string `json:"organization_id,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
	Metro          string `json:"metro,omitempty"`
	Token          string `json:"token,omitempty"`
	URL            string `json:"url,omitempty"`
	Plan           string `json:"plan,omitempty"`
	OS             string `json:"os,omitempty"`
	Debug          bool   `json:"debug,omitempty"`
}

// DefaultConfig ...
type DefaultConfig struct {
	Config      *MetalConfig
	ConfigFiles []string
	ConfigPath  []string
}

// Service ...
type Service interface {
	Config

	NewClient(...ClientConfigurator) (*Client, error)
}

// SessionMaker ...
type SessionMaker interface {
	NewSession(...ServiceConfigurator)
}

// ServiceConfigurator ...
type ServiceConfigurator func(Service)

var (
	// ConfigFromEnv ...
	// TODO(displague) sessions should report the definitive Token, URL, etc
	// after consulting with each serviceconfigurator
	// services should report nil to a Session if the value is not set by that services configurator.
	//
	//

	ConfigFromEnv ServiceConfigurator = func(s Service) {
		// TODO instead of this pattern, ServiceWriter interface can implement all the setters like (s.SetDefaultTokenFn(withToken))
		// and "ConfigFromEnv" can Set that (withToken) configurator
		// its not safe to assume that a DefaultService was used here
		if s, ok := s.(*DefaultService); ok {
			s.Config.Debug = os.Getenv(debugEnvVar) != ""
			s.Config.Token = os.Getenv(authTokenEnvVar)
			s.Config.URL = baseURL
		}
	}

	// ConfigFromConfig ...
	ConfigFromConfig ServiceConfigurator = func(s Service) {

	}

	// ConfigFromMetadata ...
	ConfigFromMetadata ServiceConfigurator = func(s Service) {

	}
)

// ClientConfigurator ...
type ClientConfigurator func(*Client, Config) error

// DefaultService ...
type DefaultService struct {
	DefaultConfig
}

// NewClient ...
func (s *DefaultService) NewClient(configs ...ClientConfigurator) (*Client, error) {
	client := &Client{}
	for _, c := range configs {
		err := c(client, s)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

var _ Service = (*DefaultService)(nil)

// DebugEnabled ...
func (s *DefaultConfig) DebugEnabled() bool {
	return s.Config.Debug
}

// URL ...
func (s *DefaultConfig) URL() string {
	return s.Config.URL
}

// Token ...
func (s *DefaultConfig) Token() string {
	return s.Config.Token

}

// DefaultProject ...
func (s *DefaultConfig) DefaultProject() string {
	return ""
}

// DefaultOrganization ...
func (s *DefaultConfig) DefaultOrganization() string {
	return ""
}

// DefaultMetro ...
func (s *DefaultConfig) DefaultMetro() string {
	return ""
}

// DefaultDevicePlan ...
func (s *DefaultConfig) DefaultDevicePlan() string {
	return ""
}

// DefaultDeviceOS ...
func (s *DefaultConfig) DefaultDeviceOS() string {
	return ""
}

// NewSession ...
func (s *DefaultService) NewSession(configs ...ServiceConfigurator) {
	for _, c := range configs {
		c(s)
	}
}
