package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	url   string
	token string
}

func NewHttpClient(uri, token string) *HttpClient {
	return &HttpClient{url: uri, token: token}
}

func (h *HttpClient) PutJson(data interface{}) ([]byte, error) {
	// Marshal data -> json data
	var bodyData []byte
	if data != nil {
		var err error
		bodyData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	// Create request
	req, err := http.NewRequest(http.MethodPut, h.url, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: req.Header.Add("token", h.token)

	// Request to update services's status to server.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("HTTP Error, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (h *HttpClient) PostJson(data interface{}) ([]byte, error) {
	// Marshal data -> json data
	var bodyData []byte
	if data != nil {
		var err error
		bodyData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, h.url, bytes.NewBuffer(bodyData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "json/application")
	// req.Header.Add("token", h.token)

	// Request to update services's status to server.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("HTTP Error, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
