package testutils

import (
	"fmt"
	"net/http"
)

var Port int32 = 6677

func GetNextPort() int32 {
	Port++
	return Port
}

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
