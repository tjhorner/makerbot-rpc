// Package jsonrpc implements MakerBot's non-standard JSON-RPC 2.0 protocol.
package jsonrpc

import "encoding/json"

// NewClient creates a new JSON-RPC client
func NewClient(ip, port string) *Client {
	return &Client{
		IP:   ip,
		Port: port,
		rsps: make(map[string]chan rpcResponse),
		subs: make(map[string]func(json.RawMessage)),
	}
}
