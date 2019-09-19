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
	server := &http.Server{Addr: portString}
	http.HandleFunc("/", handler)
	go func() {
		err := http.ListenAndServe(portString, nil)
		fmt.Printf("Http server running, listening on port %d", port)
		if err != nil {
			panic(err)
		}
	}()
	return server
}
