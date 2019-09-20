package testutils

import (
	"fmt"
	"net/http"
)

// Port is the current port number
var Port int32 = 6677

// GetNextPort returns the next port number and increments the Port
func GetNextPort() int32 {
	Port++
	return Port
}

// StartServer starts a local http listener
func StartServer(port int32, handler func(writer http.ResponseWriter, request *http.Request)) *http.Server {
	portString := fmt.Sprintf(":%d", port)

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", handler)

	server := &http.Server{
		Addr:    portString,
		Handler: serveMux,
	}

	go func() {
		server.ListenAndServe()
	}()
	return server
}
