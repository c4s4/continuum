package main

/*
 * Provide a function that ensures that only one instance of this program may run
 * on a given machine.
 */

import (
	"fmt"
	"net"
)

func listenPortOrExit(port int) {
	socket, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic("Another instance is already running")
	}
	for {
		socket.Accept()
	}
}

// Highlander returns silently if no other instance is running or it will cause a
// panic
func Highlander(port int) {
	go listenPortOrExit(port)
}
