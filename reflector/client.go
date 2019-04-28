package reflector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client is an HTTP client that talks to MakerBot Reflector.
type Client struct {
	BaseURL     string
	accessToken string
	http        *http.Client
}

func (c *Client) url(endpoint string) string {
	return fmt.Sprintf("%s%s", c.BaseURL, endpoint)
}

func (c *Client) httpGet(endpoint string) (*json.RawMessage, error) {
	req, err := http.NewRequest("GET", c.url(endpoint), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	r, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	jr := json.RawMessage(resp)

	return &jr, nil
}

func (c *Client) httpPost(endpoint string, params map[string]string) (*json.RawMessage, error) {
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", c.url(endpoint), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	r, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	jr := json.RawMessage(resp)

	return &jr, nil
}

// GetPrinters gets a list of printers connected to the Thingiverse account
func (c *Client) GetPrinters() (*json.RawMessage, error) {
	return c.httpGet("/printers")
}

// GetPrinter gets a printer with `id`
func (c *Client) GetPrinter(id string) (*json.RawMessage, error) {
	return c.httpGet(fmt.Sprintf("/printers/%s", id))
}

// CallPrinter returns a relay on which you can attach a makerbot.Client
func (c *Client) CallPrinter(id string) (*CallPrinterResponse, error) {
	resp, err := c.httpPost("/call", map[string]string{"printer_id": id})
	if err != nil {
		return nil, err
	}

	var res CallPrinterResponse
	json.Unmarshal(*resp, &res)

	return &res, nil
}
