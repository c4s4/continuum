/*
 * Provide a function that ensures that only one instance of this program may run
 * on a given machine.
 */

package main

import (
	"fmt"
	"net"
)

// listen keeps port open endlessly.
func listen(socket net.Listener) {
	for {
		socket.Accept()
	}
}

// IsAnotherInstanceRunning tells if another instance is running and keeps a lock
// (opening a port) if there is no other instance running.
func IsAnotherInstanceRunning(port int) bool {
	socket, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return true
	} else {
		go listen(socket)
		return false
	}
}
