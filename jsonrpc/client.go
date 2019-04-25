package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

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
	Data    struct {
		Args []string `json:"args"`
		Name string   `json:"name"`
	} `json:"data"`
}

func (e *rpcError) Error() string { return fmt.Sprintf("rpc error: %s: %s", e.Data.Name, e.Message) }

type rpcResponse struct {
	ID      *string          `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Version string           `json:"jsonrpc"`
	Error   *rpcError        `json:"error,omitempty"`
}

// Client is a JSON-RPC client
type Client struct {
	IP   string
	Port string
	rsps map[string]chan rpcResponse
	subs map[string]func(json.RawMessage)
	conn *net.TCPConn
}

// Connect connects to the remote JSON-RPC server
func (c *Client) Connect() error {
	addr, err := net.ResolveTCPAddr("tcp", c.IP+":"+c.Port)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	jd := make(chan []byte)
	jr := NewJSONReader(jd)

	go func() {
		for {
			b := make([]byte, 1)
			conn.Read(b)

			jr.FeedByte(b[0])
		}
	}()

	go func() {
		for {
			j := <-jd
			// log.Println(string(j))

			// need to determine if this is a request or a response
			var resp rpcResponse
			json.Unmarshal(j, &resp)

			// log.Printf("recv'd json: %+v\n", resp)

			if resp.Result == nil && resp.Error == nil && resp.ID == nil {
				// Request
				var req rpcServerRequest
				json.Unmarshal(j, &req)

				if sub, ok := c.subs[req.Method]; ok {
					// log.Printf("request: %+v\n", req)
					sub(req.Params)
				}
			} else {
				// Response
				if rsp, ok := c.rsps[*resp.ID]; ok {
					// log.Printf("response: %+v\n", resp)
					rsp <- resp
					delete(c.rsps, *resp.ID)
					// log.Printf("%+v\n", c.rsps)
				}
			}
		}
	}()

	c.conn = conn

	return nil
}

// Close closes the underlying TCP connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

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
	// log.Printf("ID: %+v\n", id)
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

	var msg chan rpcResponse
	if reply != nil {
		msg = make(chan rpcResponse)
		c.rsps[id] = msg
	}

	conn.Write(marshaledReq)

	if reply != nil {
		resp := <-msg
		// log.Println("this is good " + serviceMethod)
		if resp.Error != nil {
			return resp.Error
		}

		if resp.Result == nil {
			return nil
		}

		json.Unmarshal(*resp.Result, &reply)
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
