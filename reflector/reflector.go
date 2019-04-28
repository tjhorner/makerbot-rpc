/*
Package reflector is a Go client library for the MakerBot Reflector API.
*/
package reflector

import "net/http"

// NewClient returns a Client with the specified access
// token
func NewClient(accessToken string) Client {
	return Client{
		BaseURL:     "https://reflector.makerbot.com",
		accessToken: accessToken,
		http:        &http.Client{},
	}
}
