package packngo

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

type MetalConfig struct {
	OrganizationID string `json:"organization_id,omitempty"`
	ProjectID      string `json:"project_id,omitempty"`
	Metro          string `json:"metro,omitempty"`
	Token          string `json:"token,omitempty"`
	URL            string `json:"url,omitempty"`
	Plan           string `json:"plan,omitempty"`
	OS             string `json:"os,omitempty"`
	Debug          string `json:"debug,omitempty"`
}

type DefaultConfig struct {
	Config      *MetalConfig
	ConfigFiles []string
	ConfigPath  []string
}

// Service ...
type Service interface {
	Config

	NewClient() *Client
}

// SessionMaker ...
type SessionMaker interface {
	NewSession(...Configurator)
}

// Configurator ...
type Configurator func(Service)

var (
	// ConfigFromEnv ...
	ConfigFromEnv Configurator = func(s Service) {

	}

	// ConfigFromConfig ...
	ConfigFromConfig Configurator = func(s Service) {

	}

	// ConfigFromMetadata ...
	ConfigFromMetadata Configurator = func(s Service) {

	}
)

// DefaultService ...
type DefaultService struct {
	DefaultConfig
}

func (s *DefaultService) NewCient() *Client {
	return NewClientWithAuth()
}

var _ Service = (*DefaultService)(nil)

// DebugEnabled ...
func (s *DefaultConfig) DebugEnabled() bool {
	return false
}

// URL ...
func (s *DefaultConfig) URL() string {
	return ""
}

// Token ...
func (s *DefaultConfig) Token() string {
	return ""
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
func (s *DefaultService) NewSession(configs ...Configurator) {
	for _, c := range configs {
		c(s)
	}
}
