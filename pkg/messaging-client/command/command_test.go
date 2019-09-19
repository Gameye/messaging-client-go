package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	testutils "../../../testutils"
	"github.com/stretchr/testify/assert"
)

func TestSendJSONObject(test *testing.T) {

	type NestedType struct {
		NestedKey string `json:"nestedKey"`
	}

	type TestPayload struct {
		KeyOne   int64      `json:"keyOne"`
		KeyTwo   string     `json:"keyTwo"`
		KeyThree NestedType `json:"keyThree"`
	}

	payload := TestPayload{
		KeyOne:   1,
		KeyTwo:   "two",
		KeyThree: NestedType{"nestedValue"},
	}

	commandHandler := func(writer http.ResponseWriter, request *http.Request) {
		var requestPayload TestPayload
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&requestPayload)

		assert.Nil(test, err)
		assert.Equal(test, requestPayload.KeyOne, payload.KeyOne)
	}

	server := testutils.StartServer(testutils.GetNextPort(), commandHandler)

	url := fmt.Sprintf("http://127.0.0.1:%d/command/test", testutils.Port)
	err := Invoke(url, payload, map[string]string{"X-HEADER-KEY": "HEADER_VALUE"})
	assert.Nil(test, err)
	server.Close()
}

func TestErrorsWhenAServerErrorOccurs(test *testing.T) {

	commandHandler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	server := testutils.StartServer(testutils.GetNextPort(), commandHandler)

	url := fmt.Sprintf("http://127.0.0.1:%d/command/test", testutils.Port)
	err := Invoke(url, nil, nil)
	assert.NotNil(test, err)
	server.Close()
}

func TestErrorsWhenAJSONSerializationErrorOccurs(test *testing.T) {
	type TestPayload struct {
		KeyOne func(string)
		KeyTwo string `json:"keyTwo"`
	}

	payload := TestPayload{
		KeyOne: func(notValid string) {},
		KeyTwo: "two",
	}

	commandHandler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusInternalServerError)
	}

	server := testutils.StartServer(testutils.GetNextPort(), commandHandler)

	url := fmt.Sprintf("http://127.0.0.1:%d/command/test", testutils.Port)
	err := Invoke(url, payload, nil)
	assert.NotNil(test, err)
	server.Close()
}
