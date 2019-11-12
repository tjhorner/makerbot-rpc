package makerbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tjhorner/makerbot-rpc/jsonrpc"
	"github.com/tjhorner/makerbot-rpc/reflector"
)

const printFileBlockSize = 50000 // 50 KB

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
	Connected bool
	IP        string
	Port      string
	Printer   *Printer
	Timeout   time.Duration
	verbose   bool
	stateCbs  []func(old, new *PrinterMetadata)
	cameraCh  *chan CameraFrame
	cameraCbs []func(*CameraFrame)
	discCb    *func()
	rpc       *jsonrpc.Client
	mux       sync.Mutex // special mutex for sending print parts
}

// SetVerbose will enable or disable verbose logging for both
// the client and JSON-RPC client.
func (c *Client) SetVerbose(verbose bool) {
	c.verbose = verbose
	if c.rpc != nil {
		c.rpc.Verbose = verbose
	}
}

func (c *Client) logVerbose(format string, a ...interface{}) {
	if !c.verbose {
		return
	}

	fmt.Printf("[makerbot.Client] %v\n", fmt.Sprintf(format, a...))
}

// HandleDisconnect calls `cb` when the printer has been
// disconnected for some reason.
//
// At this point, you should stop using this Client and create
// a new one. It is currently not safe to continue using Clients
// when they encounter a bad disconnected state. Hopefully the GC
// does its job.
func (c *Client) HandleDisconnect(cb func()) {
	c.discCb = &cb
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
// `id` and will authenticate using the MakerBot account token `accessToken`.
//
// Since authentication is already performed by Reflector, you do not need
// to perform any additional authentication after it is connected.
func (c *Client) ConnectRemote(id, accessToken string, useRefl ...*reflector.Client) error {
	var refl reflector.Client
	if len(useRefl) == 0 {
		refl = reflector.NewClient(accessToken)
	} else {
		refl = *useRefl[0]
	}

	call, err := refl.CallPrinter(id)
	if err != nil {
		return err
	}

	split := strings.Split(call.Call.Relay, ":")

	if len(split) < 2 {
		return fmt.Errorf("reflector relay address was malformed (%s)", call.Call.Relay)
	}

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
	c.rpc.Verbose = c.verbose

	err := c.rpc.Connect()
	if err != nil {
		return err
	}

	c.Connected = true

	return nil
}

func (c *Client) handshake() error {
	c.rpc.HandleReadError(func(err error) {
		c.Connected = false
		if c.discCb != nil {
			(*c.discCb)()
		}
	})

	printer, err := c.sendHandshake()
	if err != nil {
		return err
	}

	c.Printer = printer

	// Ping-pong!
	go func() {
		for {
			c.mux.Lock()

			resp := make(chan bool, 1)

			go func() {
				res, err := c.ping()
				if err != nil {
					resp <- false
				}

				resp <- *res
			}()

			select {
			case <-resp:
				// Do nothing
			case <-time.After(c.Timeout):
				c.Connected = false
				c.Close()

				if c.discCb != nil {
					(*c.discCb)()
				}
				return
			}

			c.mux.Unlock()

			time.Sleep(10 * time.Second)
		}
	}()

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
		metadata := unpackCameraFrameMetadata(c.rpc.GetRawData(16))
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
	if !c.Connected {
		return errors.New("client is not connected to printer")
	}

	return c.rpc.Call(method, args, &result)
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
