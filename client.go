package makerbot

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/tjhorner/makerbot-rpc/reflector"

	"github.com/tjhorner/makerbot-rpc/jsonrpc"
)

type rpcEmptyParams struct{}

type rpcSystemNotification struct {
	Info *PrinterMetadata `json:"info"`
}

// Client represents an RPC client that can connect to
// MakerBot 3D printers via the network.
//
// Calls to the printer (e.g. LoadFilament, Cancel, etc.)
// will block, so you may want to take this into consideration.
type Client struct {
	IP        string
	Port      string
	Printer   *Printer
	stateCbs  []func(old, new *PrinterMetadata)
	cameraCh  *chan CameraFrame
	cameraCbs []func(*CameraFrame)
	rpc       *jsonrpc.Client
}

// ConnectLocal connects to a local printer and performs the initial handshake.
// If it is successful, the Printer field will be populated with information
// about the machine this client is connected to.
//
// After using ConnectLocal, you must use one of the AuthenticateWith* methods
// to authenticate with the printer.
func (c *Client) ConnectLocal(ip, port string) error {
	c.IP = ip
	c.Port = port

	err := c.connectRPC()
	if err != nil {
		return err
	}

	return c.handshake()
}

// ConnectRemote uses MakerBot Reflector to remotely connect to a printer
// and performs the initial handshake. It will connect to printer with ID
// `id` and will authenticate using the Thingiverse token `accessToken`.
//
// Since authentication is already performed by Thingiverse, you do not need
// to perform any additional authentication after it is connected.
func (c *Client) ConnectRemote(id, accessToken string) error {
	refl := reflector.NewClient(accessToken)

	call, err := refl.CallPrinter(id)
	if err != nil {
		return err
	}

	split := strings.Split(call.Call.Relay, ":")
	c.IP = split[0]
	c.Port = split[1]

	err = c.connectRPC()
	if err != nil {
		return err
	}

	ok, err := c.sendAuthPacket(id, call)
	if err != nil {
		return err
	}

	if !*ok {
		return errors.New("could not authenticate with printer via Reflector call")
	}

	return c.handshake()
}

func (c *Client) connectRPC() error {
	c.rpc = jsonrpc.NewClient(c.IP, c.Port)
	return c.rpc.Connect()
}

func (c *Client) handshake() error {
	printer, err := c.sendHandshake()
	if err != nil {
		return err
	}

	c.Printer = printer

	onStateChange := func(message json.RawMessage) {
		var oldState *PrinterMetadata
		if c.Printer != nil {
			oldState = c.Printer.Metadata
		}

		var newState rpcSystemNotification
		json.Unmarshal(message, &newState)

		c.Printer.Metadata = newState.Info

		for _, cb := range c.stateCbs {
			go cb(oldState, newState.Info) // Async so we don't block other callbacks
		}
	}

	c.rpc.Subscribe("system_notification", onStateChange)
	c.rpc.Subscribe("state_notification", onStateChange)

	c.rpc.Subscribe("camera_frame", func(m json.RawMessage) {
		if len(c.cameraCbs) == 0 {
			go c.endCameraStream()
		}

		metadata := parseCameraFrameMetadata(c.rpc.GetRawData(16))

		data := c.rpc.GetRawData(int(metadata.FileSize))

		frame := CameraFrame{
			Data:     data,
			Metadata: &metadata,
		}

		if c.cameraCh != nil {
			*c.cameraCh <- frame
			c.cameraCh = nil
		}

		for _, cb := range c.cameraCbs {
			go cb(&frame) // Async so we don't block other callbacks
		}
	})

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

// HandleCameraFrame calls `cb` when the printer sends a camera frame.
func (c *Client) HandleCameraFrame(cb func(frame *CameraFrame)) {
	c.cameraCbs = append(c.cameraCbs, cb)
	go c.requestCameraStream()
}

func (c *Client) call(method string, args, result interface{}) error {
	if c.rpc == nil {
		return errors.New("client is not connected to printer")
	}

	return c.rpc.Call(method, args, &result)
}

func (c *Client) sendHandshake() (*Printer, error) {
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

type rpcAuthPacketParams struct {
	CallID     string `json:"call_id"`
	ClientCode string `json:"client_code"`
	PrinterID  string `json:"printer_id"`
}

func (c *Client) sendAuthPacket(id string, pc *reflector.CallPrinterResponse) (*bool, error) {
	params := rpcAuthPacketParams{
		CallID:     pc.Call.ID,
		ClientCode: pc.Call.ClientCode,
		PrinterID:  id,
	}

	var reply bool
	return &reply, c.call("auth_packet", params, &reply)
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

func (c *Client) requestCameraStream() error {
	return c.call("request_camera_stream", rpcEmptyParams{}, nil)
}

func (c *Client) endCameraStream() error {
	return c.call("end_camera_stream", rpcEmptyParams{}, nil)
}

// GetCameraFrame requests a single frame from the printer's camera
func (c *Client) GetCameraFrame() (*CameraFrame, error) {
	ch := make(chan CameraFrame)
	c.cameraCh = &ch

	err := c.requestCameraStream()
	if err != nil {
		return nil, err
	}

	data := <-ch
	close(ch)

	return &data, nil
}
