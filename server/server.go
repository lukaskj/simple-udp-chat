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

type Server struct {
	running    bool
	connection *net.UDPConn
	clients    map[string]Client
}

func (s *Server) Start(port string) {
	addr, err := net.ResolveUDPAddr("udp", port)
	common.CheckError(err)
	s.clients = make(map[string]Client)
	s.running = true
	s.connection, err = net.ListenUDP("udp", addr)
	common.CheckError(err)
}

func (s *Server) Disconnect() {
	s.connection.Close()
}

func (s *Server) handle() {
	for s.running {
		var buf [8192]byte
		n, addr, err := s.connection.ReadFromUDP(buf[0:])
		if err != nil {
			if(!strings.Contains(err.Error(), "use of closed net")) {
				common.CheckError(err)
			}
		}
		// common.CheckError(err)
		rawMessage := buf[0:n]
		var msg protocol.Message
		msg.Parse(rawMessage) // Parse the result to type Message
		switch msg.Action {
		case protocol.ACTION_CONNECT:
			c := Client{}
			c.address = addr
			c.id = uuid.Must(uuid.NewV4()).String()
			common.Log("'" + msg.Body + "' connected ('" + c.id + ", " + c.address.String() + ")")

			message := protocol.Message{}
			message.Action = protocol.ACTION_OK
			message.Body = c.id
			msgToClient, err := message.Serialize()
			common.CheckError(err)
			_, err = s.connection.WriteToUDP(msgToClient, c.address)
		}
	}
	wg.Done()
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
		default:
			fmt.Println("Command '" + command + "' not found")
		}
	}
}
