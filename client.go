package makerbot

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/tjhorner/makerbot-rpc/jsonrpc"
)

type rpcEmptyParams struct{}

type rpcSystemNotification struct {
	Info *PrinterMetadata `json:"info"`
}

// Client represents an RPC client that can connect to
// MakerBot 3D printers via the network
type Client struct {
	IP       string
	Port     string
	Printer  *Printer
	stateCbs []func(old, new *PrinterMetadata)
	http     *http.Client
	rpc      *jsonrpc.Client
}

// Connect connects to the printer and performs the initial handshake.
// If it is successful, the Printer field will be populated with information
// about the machine this client is connected to.
func (c *Client) Connect() error {
	if c.IP == "" || c.Port == "" {
		return errors.New("IP and Port are required fields for Client")
	}

	rpc := jsonrpc.NewClient(c.IP, c.Port)
	err := rpc.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	onStateChange := func(message json.RawMessage) {
		oldState := c.Printer.Metadata

		var newState rpcSystemNotification
		json.Unmarshal(message, &newState)

		c.Printer.Metadata = newState.Info

		for _, cb := range c.stateCbs {
			cb(oldState, newState.Info)
		}
	}

	rpc.Subscribe("system_notification", onStateChange)
	rpc.Subscribe("state_notification", onStateChange)

	c.rpc = rpc

	printer, err := c.handshake()
	if err != nil {
		return err
	}

	c.Printer = printer

	return nil
}

// Close closes the underlying TCP socket
// and should be called when the client is no
// longer needed
func (c *Client) Close() error {
	if c.rpc == nil {
		return nil // Nothing to do
	}

	return c.rpc.Close()
}

// HandleStateChange calls `cb` when the printer's state changes.
//
// The first parameter passed to `cb` is the previous state, and the
// second is the new state. You can use this to respond when e.g. a print
// fails for some reason, or when a print's progress changes.
func (c *Client) HandleStateChange(cb func(old, new *PrinterMetadata)) {
	c.stateCbs = append(c.stateCbs, cb)
}

func (c *Client) call(method string, args, result interface{}) error {
	if c.rpc == nil {
		return errors.New("client is not connected to printer")
	}

	return c.rpc.Call(method, args, &result)
}

func (c *Client) handshake() (*Printer, error) {
	var reply Printer
	return &reply, c.call("handshake", rpcEmptyParams{}, &reply)
}

type rpcAuthenticateParams struct {
	AccessToken string `json:"access_token"`
}

// authenticate performs authentication with the printer
// via an access token retrieved through the printer's
// HTTP server
func (c *Client) authenticate(accessToken string) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("authenticate", rpcAuthenticateParams{accessToken}, &reply)
}

// AuthenticateWithThingiverse performs authentication with the printers
// by using a Thingiverse token:username pair.
//
// Ensure that you have authenticated your Thingiverse account with this printer
// at least once in the past. You can do this logging into the MakerBot Print
// application and connecting to the printer.
func (c *Client) AuthenticateWithThingiverse(token, username string) error {
	accessToken, err := c.getThingiverseAccessToken(token, username)
	if err != nil {
		return err
	}

	_, err = c.authenticate(*accessToken)
	return err
}

type rpcLoadUnloadFilamentParams struct {
	ToolIndex int `json:"tool_index"`
}

// LoadFilament instructs the printer to begin loading filament into
// the extruder
func (c *Client) LoadFilament(toolIndex int) (*PrinterProcess, error) {
	var reply PrinterProcess
	return &reply, c.call("load_filament", rpcLoadUnloadFilamentParams{toolIndex}, &reply)
}

// UnloadFilament instructs the printer to begin unloading filament from
// the extruder
func (c *Client) UnloadFilament(toolIndex int) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("unload_filament", rpcLoadUnloadFilamentParams{toolIndex}, &reply)
}

// Cancel instructs the printer to cancel the current process, if any.
//
// It may result in a `ProcessNotCancellableException`, so you may want to
// check the `CurrentProcess` to ensure it is `Cancellable`. Or not, if you
// don't care if it fails.
func (c *Client) Cancel() (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("cancel", rpcEmptyParams{}, &reply)
}

type rpcChangeMachineNameParams struct {
	MachineName string `json:"machine_name"`
}

// ChangeMachineName instructs the printer to change its display name.
func (c *Client) ChangeMachineName(name string) (*json.RawMessage, error) {
	var reply json.RawMessage
	return &reply, c.call("cancel", rpcChangeMachineNameParams{name}, &reply)
}
