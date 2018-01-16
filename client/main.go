package main

import (
	"fmt"
	"net"
	"../common"
	"../common/protocol"
	"encoding/binary"
)

type Client struct {
	connection *net.UDPConn
	id string
	username string
	receivedMessages chan string
	sentMessages chan string
	connected bool
}

func (c* Client) connect() {
	msg := protocol.Message{Action: protocol.ACTION_CONNECT}
	// create temporary client to send as message
	// tempClient := make(map[string]string)
	// tempClient[c.id] = c.username
	// tempClientBody, _ := json.Marshal(tempClient)
	// msg.Body = string(tempClientBody)
	msg.Body = c.username
	sendMsg, err := msg.Serialize()
	common.CheckError(err)
	_, err = c.connection.Write(sendMsg)
	common.CheckError(err)

	c.handleReturn(msg.Action)
}

func (c *Client) handleReturn(action uint16) {
	var buf [8192]byte
	n, err := c.connection.Read(buf[0:])
	common.CheckError(err)
	if n < 2 {
		return
	}
	msg := protocol.Message{Action: binary.BigEndian.Uint16(buf[0:2]), Body: string(buf[2:n])}
	fmt.Println("handleReturn", msg)
}






func main() {
	serverAddr := "localhost:8080"
	udpAddr, err := net.ResolveUDPAddr("udp4", serverAddr)
	common.CheckError(err)

	var client Client = Client {connected: false}
	client.receivedMessages = make(chan string)
	client.sentMessages = make(chan string)

	// *************
	fmt.Print("Digite seu nome: ")
	fmt.Scanln(&client.username)
	client.id = client.username
	// *************

	client.connection, err = net.DialUDP("udp", nil, udpAddr)
	common.CheckError(err)

	defer client.connection.Close()

	client.connect()
}

