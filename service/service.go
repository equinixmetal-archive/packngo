package service

// Service ...
type Service interface {
	URL() string

	DefaultProject() string
	DefaultOrganization() string
	DefaultMetro() string
	DefaultDeviceOS() string

	Token() string
	DebugEnabled() bool

	NewClient()
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
}

var _ Service = (*DefaultService)(nil)

// DebugEnabled ...
func (s *DefaultService) DebugEnabled() bool {
	return false
}

// URL ...
func (s *DefaultService) URL() string {
	return ""
}

// Token ...
func (s *DefaultService) Token() string {
	return ""
}

// DefaultProject ...
func (s *DefaultService) DefaultProject() string {
	return ""
}

// DefaultOrganization ...
func (s *DefaultService) DefaultOrganization() string {
	return ""
}

// DefaultMetro ...
func (s *DefaultService) DefaultMetro() string {
	return ""
}

// DefaultDeviceOS ...
func (s *DefaultService) DefaultDeviceOS() string {
	return ""
}

// NewSession ...
func (s *DefaultService) NewSession(configs ...Configurator) {
	for _, c := range configs {
		c(s)
	}
}
