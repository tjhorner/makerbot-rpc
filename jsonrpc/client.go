package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
)

type rpcRequest struct {
	ID      string `json:"id"`
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
}

type rpcClientRequest struct {
	Params interface{} `json:"params"`
	rpcRequest
}

type rpcServerRequest struct {
	Params json.RawMessage `json:"params"`
	rpcRequest
}

type rpcEmptyParams struct{}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    json.RawMessage
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error (remote): %s: %s", string(e.Data), e.Message)
}

type rpcResponse struct {
	ID      *string          `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Version string           `json:"jsonrpc"`
	Error   *rpcError        `json:"error,omitempty"`
}

// Client is a JSON-RPC client
type Client struct {
	IP      string
	Port    string
	Verbose bool
	rsps    map[string]chan rpcResponse
	subs    map[string]func(json.RawMessage)
	jr      JSONReader
	errCb   *func(error)
	conn    *net.TCPConn
	mux     sync.Mutex
	rMux    sync.Mutex
}

func (c *Client) logVerbose(format string, a ...interface{}) {
	if !c.Verbose {
		return
	}

	fmt.Printf("[jsonrpc.Client] %v\n", fmt.Sprintf(format, a...))
}

// Connect connects to the remote JSON-RPC server
func (c *Client) Connect() error {
	c.logVerbose("resolving TCP address %s:%s", c.IP, c.Port)

	addr, err := net.ResolveTCPAddr("tcp", c.IP+":"+c.Port)
	if err != nil {
		return err
	}

	c.logVerbose("dialing resolved tcp address %s", addr.String())

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	c.logVerbose("TCP connection: %+v", *conn)

	conn.SetKeepAlive(true)

	done := func(j []byte) error {
		c.logVerbose("received JSON packet: %s", string(j))

		if !json.Valid(j) {
			c.logVerbose("it turns out the JSON packet was invalid")
			return errors.New("invalid JSON")
		}

		// need to determine if this is a request or a response
		var resp rpcResponse
		err := json.Unmarshal(j, &resp)
		if err != nil {
			c.logVerbose("error unmarshaling RPC response: %s", err.Error())
			return err
		}

		c.logVerbose("rpc response unmarshaled: %+v", resp)

		if resp.Result == nil && resp.Error == nil && resp.ID == nil {
			// Request
			var req rpcServerRequest
			json.Unmarshal(j, &req)

			if sub, ok := c.subs[req.Method]; ok {
				go sub(req.Params)
			}
		} else if resp.ID != nil {
			// Response
			if rsp, ok := c.rsps[*resp.ID]; ok {
				go func() { rsp <- resp }()
				c.rMux.Lock()
				delete(c.rsps, *resp.ID)
				c.rMux.Unlock()
			}
		}

		return nil
	}

	c.jr = NewJSONReader(done)

	go func() {
		// temporary array to pipe from the TCP connection to the
		// jsonreader
		// TODO: can we use io.Pipe here maybe?
		b := make([]byte, 1)

		for {
			_, err := conn.Read(b)

			if err != nil {
				c.conn.Close()
				c.conn = nil

				if c.errCb != nil {
					(*c.errCb)(err)
				}
				break
			}

			c.jr.Write(b)
		}
	}()

	c.conn = conn

	return nil
}

// HandleReadError calls `cb` when an error occurs while
// reading from the underlying TCP socket
func (c *Client) HandleReadError(cb func(error)) {
	c.errCb = &cb
}

// Close closes the underlying TCP connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	c.jr.Reset()
	return (*c.conn).Close()
}

// Call calls the remote JSON-RPC server with `serviceMethod`
func (c *Client) Call(serviceMethod string, args, reply interface{}) error {
	if c.conn == nil {
		return errors.New("Client is not connected (hint: call Connect())")
	}

	conn := *c.conn

	if args == nil {
		args = rpcEmptyParams{}
	}

	id := uuid.New().String()

	req := rpcClientRequest{
		Params: args,
	}

	req.ID = id
	req.Version = "2.0"
	req.Method = serviceMethod

	marshaledReq, err := json.Marshal(req)
	if err != nil {
		return err
	}

	c.mux.Lock()

	var msg chan rpcResponse
	if reply != nil {
		msg = make(chan rpcResponse)

		c.rMux.Lock()
		c.rsps[id] = msg
		c.rMux.Unlock()
	}

	conn.Write(marshaledReq)
	c.mux.Unlock()

	if reply != nil {
		resp := <-msg

		if resp.Error != nil {
			return resp.Error
		}

		if resp.Result == nil {
			return nil
		}

		json.Unmarshal(*resp.Result, &reply)
		close(msg)
	}

	return nil
}

// Subscribe subscribes to a notification channel that will be sent by the
// remote server. Every time something is received via that channel, `cb` will
// be called with the raw JSON the server sent in the `Params`. From there, you
// should unmarshal it yourself.
func (c *Client) Subscribe(namespace string, cb func(message json.RawMessage)) error {
	if _, ok := c.subs[namespace]; ok {
		return errors.New("already subscribed to " + namespace)
	}

	if c.conn == nil {
		return errors.New("Client is not connected (hint: call Connect())")
	}

	c.subs[namespace] = cb
	return nil
}

// Unsubscribe will unsubscribe any listeners from specified notification
// channel. It is safe to call this method even if there is nothing subscribed
// to the channel.
func (c *Client) Unsubscribe(namespace string) {
	delete(c.subs, namespace)
}

// GetRawData grabs raw data from the TCP connection until
// `length` is reached. The captured data is returned as an
// array of bytes.
func (c *Client) GetRawData(length int) []byte {
	return c.jr.GetRawData(length)
}

// Write writes bytes to the underlying TCP socket.
func (c *Client) Write(bs []byte) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.conn.Write(bs)
}
