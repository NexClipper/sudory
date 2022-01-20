package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClient struct {
	url   string
	token string
}

func NewHttpClient(uri, token string) *HttpClient {
	return &HttpClient{url: uri, token: token}
}

func (h *HttpClient) PutJson(param map[string]string, data interface{}) ([]byte, error) {
	var buffer *bytes.Buffer

	// Marshal data -> json data
	if data != nil {
		bodyData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewBuffer(bodyData)
	}

	// Create request
	req, err := http.NewRequest(http.MethodPut, h.url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// TODO: req.Header.Add("token", h.token)

	// Add param
	q := url.Values{}
	for k, v := range param {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

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

func (h *HttpClient) PostJson(param map[string]string, data interface{}) ([]byte, error) {
	var buffer *bytes.Buffer

	// Marshal data -> json data
	if data != nil {
		bodyData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buffer = bytes.NewBuffer(bodyData)
	}

	// Create request
	req, err := http.NewRequest(http.MethodPost, h.url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "json/application")
	// req.Header.Add("token", h.token)

	// Add param
	q := url.Values{}
	for k, v := range param {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

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
