/*
Package makerbot is a Go client library for MakerBot printers.
*/
package makerbot

import "net/http"

// These constants are used to communicate with the printer
// and are apparently hard-coded

const makerbotClientID = "MakerWare"
const makerbotClientSecret = "secret"

// NewClient creates a new client ready to connect to the printer
// located at `ip`. (Hint: see `Connect()`)
//
// If, for some reason, your printer does not listen on port 9999,
// you can use `NewClientWithPort` instead.
func NewClient(ip string) Client {
	return Client{
		IP:   ip,
		Port: "9999",
		http: &http.Client{},
	}
}

// NewClientWithPort does the same thing as `NewClient` except
// you can provide a custom port to use when connecting to the
// printer.
func NewClientWithPort(ip, port string) Client {
	return Client{
		IP:   ip,
		Port: port,
		http: &http.Client{},
	}
}
