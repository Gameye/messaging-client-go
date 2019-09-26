package testutils

import (
	"fmt"
	"log"
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

	serverChan := make(chan *http.Server)
	go func() {
		serverChan <- server
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	return <-serverChan
}
