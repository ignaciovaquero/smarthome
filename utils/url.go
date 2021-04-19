package utils

import (
	"fmt"
	"net/url"
)

// ValidateURLsFromArray takes an array of strings and returns an
// array of URL. If any of the URLs can't be parsed, it returns
// an error.
func ValidateURLsFromArray(urls []string) error {
	if len(urls) == 1 {
		if urls[0] == "*" {
			return nil
		}
	}
	for _, u := range urls {
		if _, err := url.Parse(u); err != nil {
			return fmt.Errorf("Invalid URL %s: %w", u, err)
		}
	}
	return nil
}
