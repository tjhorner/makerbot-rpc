package makerbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) httpGet(endpoint string, qs map[string]string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "http://"+c.IP+endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range qs {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var resp map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) getLocalAccessToken() (*string, error) {
	codeRes, err := c.httpGet("/auth", map[string]string{
		"response_type": "code",
		"client_id":     makerbotClientID,
		"client_secret": makerbotClientSecret,
	})
	if err != nil {
		return nil, err
	}

	var answerRes map[string]interface{}

	for {
		// Poll until knob is pressed
		answerRes, err = c.httpGet("/auth", map[string]string{
			"response_type": "answer",
			"client_id":     makerbotClientID,
			"client_secret": makerbotClientSecret,
			"answer_code":   codeRes["answer_code"].(string),
		})
		if err != nil {
			return nil, err
		}

		if answerRes["answer"].(string) == "accepted" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	tokenRes, err := c.httpGet("/auth", map[string]string{
		"response_type": "token",
		"client_id":     makerbotClientID,
		"client_secret": makerbotClientSecret,
		"context":       "jsonrpc",
		"auth_code":     answerRes["code"].(string),
	})
	if err != nil {
		return nil, err
	}

	accessToken, ok := tokenRes["access_token"].(string)

	if !ok {
		return nil, fmt.Errorf("could not authenticate with printer for some reason: %s", err)
	}

	return &accessToken, nil
}

func (c *Client) getThingiverseAccessToken(token, username string) (*string, error) {
	codeRes, err := c.httpGet("/auth", map[string]string{
		"response_type":     "code",
		"client_id":         makerbotClientID,
		"client_secret":     makerbotClientSecret,
		"thingiverse_token": token,
		"username":          username,
	})
	if err != nil {
		return nil, err
	}

	var answerRes map[string]interface{}

	for i := 0; i < 10; i++ {
		// Try 10 times for accepted auth
		answerRes, err = c.httpGet("/auth", map[string]string{
			"response_type": "answer",
			"client_id":     makerbotClientID,
			"client_secret": makerbotClientSecret,
			"answer_code":   codeRes["answer_code"].(string),
		})
		if err != nil {
			return nil, err
		}

		if answerRes["answer"].(string) == "accepted" {
			break
		}

		time.Sleep(2 * time.Second)
	}

	if answerRes["answer"].(string) != "accepted" {
		return nil, errors.New("could not authenticate with printer after 10 tries, please check that your Thingiverse account is authorized to access it")
	}

	tokenRes, err := c.httpGet("/auth", map[string]string{
		"response_type": "token",
		"client_id":     makerbotClientID,
		"client_secret": makerbotClientSecret,
		"context":       "jsonrpc",
		"auth_code":     answerRes["code"].(string),
	})
	if err != nil {
		return nil, err
	}

	accessToken, ok := tokenRes["access_token"].(string)

	if !ok {
		return nil, errors.New("could not authenticate with printer, please check that your Thingiverse account is authorized to access it")
	}

	return &accessToken, nil
}
