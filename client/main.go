package main

import (
	"time"
	"fmt"
	"net"
	"sync"

	"../common"
	"../common/protocol"
)

type Client struct {
	connection *net.UDPConn
	id string
	username string
	receivedMessages chan protocol.Message
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
	sendMsg := msg.Serialize()
	_, err := c.connection.Write(sendMsg)
	common.CheckError(err)

	var buf [8192]byte
	c.connection.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.connection.Read(buf[0:])
	common.CheckError(err)
	if n < 2 {
		return
	}
	msg.Parse(buf[0:n])
	switch msg.Action {
	case protocol.ACTION_OK:
		c.id = msg.Token
		c.connected = true
		common.Log("Connected")
	}
}

func (c *Client) disconnect() {
	c.connection.Close()
	c.connected = false
	close(c.receivedMessages)
	common.Log("Disconnected")
}

func (c *Client) handleServerMessages() {
	for c.connected {
		var buf [8192]byte
		c.connection.SetReadDeadline(time.Now().Add(5000 * time.Second))
		n, err := c.connection.Read(buf[0:])
		common.CheckError(err)
		if n < 2 {
			return
		}
		msg := protocol.Message{}
		msg.Parse(buf[0:n])
		c.receivedMessages <- msg
	}
}

func (c *Client) handleReceivingMessage() {
	for c.connected {
		msg := <- c.receivedMessages
		switch msg.Action {
		case protocol.ACTION_BROADCAST:
			fmt.Println(msg.Body)
		default:
			fmt.Println("Message not recognized", msg)
		}
	}
}

func (c *Client) sendMessage(msgString string) {
	msg := protocol.Message{}
	msg.Action = protocol.ACTION_BROADCAST
	msg.Body = msgString
	msg.Token = c.id
	serializedMsg := msg.Serialize()
	c.connection.Write(serializedMsg)
}

func (c *Client) handleSentMessages() {
	defer wg.Done()
	var command string
	for c.connected {
		fmt.Scanf("%s\n", &command)
		if len(command) > 0 {
			if command[0] == ':' {
				switch command {
				case ":q":
					fallthrough
				case ":quit":
					c.disconnect()
				}
			} else {
				c.sendMessage(command)
			}
		}
	}
}

var wg sync.WaitGroup

func main() {
	serverAddr := "localhost:8080"
	udpAddr, err := net.ResolveUDPAddr("udp4", serverAddr)
	common.CheckError(err)

	var client Client = Client {connected: false}
	client.receivedMessages = make(chan protocol.Message)

	// *************
	fmt.Print("Enter your name: ")
	fmt.Scanln(&client.username)
	// *************

	client.connection, err = net.DialUDP("udp", nil, udpAddr)
	common.CheckError(err)

	// defer client.disconnect()
	common.HandleExit(func() bool {
		client.disconnect()
		return true
	})

	client.connect()

	if client.connected {
		wg.Add(1)
		go client.handleSentMessages()
		go client.handleServerMessages()
		go client.handleReceivingMessage()
	}

	wg.Wait()
}

