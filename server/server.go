package main

import "../common"
import "../common/protocol"
import (
	"fmt"
	"net"
	"strings"

	"github.com/satori/go.uuid"
	// "encoding/json"
)


func checkError(err error, fnc string) {
	if err != nil {
		fmt.Println("Error:", err, "in", fnc)
	}
}

type Server struct {
	running    bool
	connection *net.UDPConn
	clients    map[string]Client
	messages	  chan protocol.Message
}

func (s *Server) start(port string) {
	addr, err := net.ResolveUDPAddr("udp", port)
	checkError(err, "Server.start")
	s.clients = make(map[string]Client)
	s.messages = make(chan protocol.Message)
	s.running = true
	s.connection, err = net.ListenUDP("udp", addr)
	checkError(err, "Server.start")
	go s.handle() // wg.Done
   go s.handleReceivingMessages()
   go s.handleCommands()
}

func (s *Server) AddClient(c Client) {

}

func (s * Server) clientExists(id string) bool {
	if(s.clients[id] == Client{}) {
		return false
	}
	return true
}

func (s *Server) disconnect() {
	_ = s.connection.Close()
	close(s.messages)
}

func (s *Server) handle() {
	defer wg.Done()
	for s.running {
		var buf [8192]byte
		n, addr, err := s.connection.ReadFromUDP(buf[0:])
		if err != nil {
			if(!strings.Contains(err.Error(), "use of closed net")) {
				checkError(err, "Server.handle")
			}
		}
		
		rawMessage := buf[0:n]
		var msg protocol.Message
		err = msg.Parse(rawMessage) // Parse the result to type Message
		if err == nil {
			msg.Payload = addr.String()
			s.messages <- msg
		}
	}
}

func (s *Server) handleReceivingMessages() {
	for s.running {
		msg := <- s.messages
		switch msg.Action {
		case protocol.ACTION_CONNECT:
			c := Client{}
			addr, err := net.ResolveUDPAddr("", msg.Payload)
			if err != nil {
				fmt.Println("URL parse error", err)
			}
			c.address = addr
			c.id = uuid.Must(uuid.NewV4()).String()
			c.username = msg.Body
			s.clients[c.id] = c // Add client to array
			connectedMsgString := "User '" + c.username + "' connected"
			common.Log(connectedMsgString + " ('" + c.id + "', " + c.address.String() + ")")

			message := protocol.Message{}
			message.Action = protocol.ACTION_OK
			message.Token = c.id
			msgToClient := message.Serialize()
			_, err = s.connection.WriteToUDP(msgToClient, c.address)
			checkError(err, "Server.handleReceivingMessages")
			
			connectedMsg := protocol.Message{Action: protocol.ACTION_BROADCAST, Body: connectedMsgString}
			for i := range s.clients {
				if(s.clients[i].id != c.id) {
					_, err = s.connection.WriteToUDP(connectedMsg.Serialize(), s.clients[i].address)
					checkError(err, "Server.handleReceivingMessages " + s.clients[i].username)
				}
			}
		case protocol.ACTION_BROADCAST:
			c := s.clients[msg.Token]
			// Just send message if the sender exists
			if s.clientExists(c.id) {
				for i := range s.clients {
					// Don't send message to sending client
					if c.id != s.clients[i].id {
						msg.Body = c.username + ": " + msg.Body
						_, err := s.connection.WriteToUDP(msg.Serialize(), s.clients[i].address)
						checkError(err, "Server.handleReceivingMessages " + s.clients[i].username)
					}
				}
			}
		default:
			fmt.Println("Default", msg)
		}
	}
}

func (s *Server) handleCommands() {
	var command string
	for s.running {
		fmt.Scanf("%s\n", &command)
		// fmt.Println("Command: " + " " + command)
		switch command {
		case "q":
			fallthrough
		case "quit":
			common.Log("Server quit")
			s.running = false
			s.connection.Close()
			break
		case "users":
			common.Log("Connected users:")
			for i := range s.clients {
				fmt.Println(s.clients[i].username, "ID:", s.clients[i].id, "Address:", s.clients[i].address)
			}
			fmt.Println("Total", len(s.clients))
		default:
			fmt.Println("Command '" + command + "' not found")
		}
	}
}
