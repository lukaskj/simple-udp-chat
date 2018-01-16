package main

import (
	"net"
	"sync"
)

type Client struct {
	id       string
	username string
	address  *net.UDPAddr
}

// MAIN *********

var wg sync.WaitGroup

func main() {
	var port string = ":8080"

	var server Server
	server.Start(port)

	// common.CheckError(err)
	defer server.Disconnect()

	wg.Add(1)
	go server.handle()
	go server.handleCommands()

	wg.Wait()
}
