package api

// API is a struct that defines the API actions for
// Smarthome app
type API struct {
	Logger
	NoSQLStorage
}

// Option is a function to apply settings to Scraper structure
type Option func(a *API) Option

// NewAPI returns a new instance of an API
func NewAPI(opts ...Option) *API {
	a := &API{
		Logger: &DefaultLogger{},
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// SetLogger sets the Logger for the API
func SetLogger(logger Logger) Option {
	return func(s *API) Option {
		prev := s.Logger
		if logger != nil {
			s.Logger = logger
		}
		return SetLogger(prev)
	}
}
