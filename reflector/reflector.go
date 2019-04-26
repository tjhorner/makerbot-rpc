/*
Package reflector is a Go client library for the MakerBot Reflector API.
*/
package reflector

import "net/http"

func NewClient(accessToken string) Client {
	return Client{
		BaseURL:     "https://reflector.makerbot.com",
		accessToken: accessToken,
		http:        &http.Client{},
	}
}
