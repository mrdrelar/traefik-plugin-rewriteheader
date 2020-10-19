package rewriteheader

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Config the plugin configuration.
type Config struct {
	FromHead string `json:"fromhead,omitempty"` // target header
	Regex    string `json:"regex,omitempty"`    // variable for creating a new header that will store data from the target header
	Create   string `json:"create,omitempty"`   // creating a new header for store extracted data from the old
	Prefix   string `json:"prefix,omitempty"`   // add prefix for a new header
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates and returns a plugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.FromHead) == 0 {
		return nil, fmt.Errorf("FromHead can't be empty")
	}
	re, err := regexp.Compile(config.Regex)

	if err != nil {
		return nil, fmt.Errorf("error compiling regex %q: %w", config.Regex, err)
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		head := req.Header.Get(config.FromHead)
		result := re.FindString(head)
		if config.Prefix != "" {
			result = config.Prefix + result
		}
		rw.Header().Set(config.Create, result)
		req.Header.Set(config.Create, result)
		next.ServeHTTP(rw, req)
	}), nil
}
