package main

import (
	"fmt"
	"net"
)

// HandleUDPData grabs the incoming UDP data
func HandleUDPData(ServerConn *net.UDPConn) {
	buf := make([]byte, 1024)
	for !exiting {
		n, addr, _ := ServerConn.ReadFromUDP(buf)
		fmt.Println("Connection from ", addr)
		if string(buf[0:n]) != "" {

			InsertValues(string(buf[0:n]))

		}
	}
	// main thread signaled us to shut down.
	fmt.Println("Shutting down UDP listener.")

	// notify the main thread we are shutting down.
	udpClosed = true
}
