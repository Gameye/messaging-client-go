package eventstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	testutils "../../../testutils"
	"github.com/stretchr/testify/assert"
)

func TestStreamsJSONObjectsFromServer(test *testing.T) {

	type NestedType struct {
		NestedKey string `json:"nestedKey"`
	}

	type TestPayload struct {
		KeyOne   int64      `json:"keyOne"`
		KeyTwo   string     `json:"keyTwo"`
		KeyThree NestedType `json:"keyThree"`
	}

	const loops = 100

	commandHandler := func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/x-ndjson")
		writer.Header().Set("Connection", "Keep-Alive")
		writer.WriteHeader(200)

		time.Sleep(10)

		encoder := json.NewEncoder(writer)

		for i := 1; i <= loops; i++ {
			payload := TestPayload{
				KeyOne:   int64(i),
				KeyTwo:   "two",
				KeyThree: NestedType{"nestedValue"},
			}
			encoder.Encode(payload)
			time.Sleep(10)
		}
	}

	server := testutils.StartServer(testutils.GetNextPort(), commandHandler)

	url := fmt.Sprintf("http://127.0.0.1:%d/event/test", testutils.Port)
	queryParams := map[string]string{
		"paramOne": "valueOne",
	}
	context := context.Background()
	decoder, err := Create(context, url, queryParams, map[string]string{"X-HEADER-KEY": "HEADER_VALUE"})

	for {
		var payload TestPayload
		err := decoder.Decode(&payload)

		if err == io.EOF {
			break
		} else if err != nil {
			test.Error(err)
			assert.NotNil(test, err)
			break
		} else {
			assert.NotNil(test, payload)
			assert.Equal(test, "two", payload.KeyTwo)
			assert.GreaterOrEqual(test, payload.KeyOne, int64(1))
			assert.LessOrEqual(test, payload.KeyOne, int64(100))
			assert.IsType(test, NestedType{}, payload.KeyThree)
		}
	}

	assert.Nil(test, err)
	server.Close()
}
