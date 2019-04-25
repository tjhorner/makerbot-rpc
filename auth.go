package makerbot

import (
	"encoding/json"
	"errors"
	"net/http"
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

	r, err := c.http.Do(req)
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

	answerRes, err := c.httpGet("/auth", map[string]string{
		"response_type": "answer",
		"client_id":     makerbotClientID,
		"client_secret": makerbotClientSecret,
		"answer_code":   codeRes["answer_code"].(string),
	})
	if err != nil {
		return nil, err
	}

	if answerRes["answer"].(string) != "accepted" {
		return nil, errors.New("could not authenticate with printer, please check that your Thingiverse account is authorized to access it")
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
