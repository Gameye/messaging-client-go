package eventstream

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
)

// Create a stream to a given url and return a json.Decoder to that stream
func Create(context context.Context,
	url string,
	queryStringParams map[string]string,
	headers map[string]string,
) (decoder *json.Decoder, err error) {
	var reqBody io.Reader

	request, err := http.NewRequest(http.MethodGet, url, reqBody)
	if err != nil {
		return nil, errors.Wrap(err, "Error: Could not create http request")
	}

	if context != nil {
		request = request.WithContext(context)
	}

	urlQuery := request.URL.Query()
	for key := range queryStringParams {
		urlQuery.Set(key, queryStringParams[key])
	}
	request.URL.RawQuery = urlQuery.Encode()

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err := getHTTPClient().Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "Error: Could not execute http request")
	}

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		msg := fmt.Sprintf(
			"Error: Http request returned not ok status: %v - %v",
			response.StatusCode,
			response.Status,
		)
		return nil, errors.New(msg)
	}

	// Start a goroutine to close the body when the context is done
	go func() {
		<-context.Done()
		err = response.Body.Close()
		if err != nil {
			log.Printf("eventstream.Create: Error closing response body: %v\n", err)
		}
	}()

	decoder = json.NewDecoder(response.Body)

	return decoder, nil
}

func getHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   0,
		Transport: http.DefaultTransport,
	}
}
