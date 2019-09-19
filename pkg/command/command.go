package command

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	httpTimeout = 30 * time.Second
)

// Invoke a command on a given url
func Invoke(url string, payload interface{}, headers map[string]string) (err error) {
	var requestBody io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return errors.Wrap(err, "Error: Could not encode json payload")
		}
		requestBody = bytes.NewBuffer(jsonData)
	}

	request, err := http.NewRequest(http.MethodPost, url, requestBody)
	if err != nil {
		return errors.Wrap(err, "Error: Could not create http request")
	}

	request.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err := getHTTPClient().Do(request)
	if err != nil {
		return errors.Wrap(err, "Error: Could not execute http request")
	}

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return errors.New("Error: Http request returned not ok status")
	}

	return nil
}

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   httpTimeout,
		Transport: http.DefaultTransport,
	}
}
