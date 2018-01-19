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
   debug      bool
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
   s.debug = false
   s.connection, err = net.ListenUDP("udp", addr)
   checkError(err, "Server.start")
   go s.handle() // wg.Done
   go s.handleReceivingMessages()
	go s.handleCommands()
	common.Log("Server running on port " + port)
}

func (s *Server) AddClient(c Client) {

}

func (s * Server) clientExists(id string) bool {
   _, ok := s.clients[id]
   return ok
}

func (s *Server) disconnect() {
	defer wg.Done()
	tempMsg := protocol.Message{Action: protocol.ACTION_DISCONNECT, Body: "Server disconnect"}
	msg := tempMsg.Serialize()
	for i := range s.clients {
		_, err := s.connection.WriteToUDP(msg, s.clients[i].address)
		checkError(err, "Server.disconnect")
	}
   _ = s.connection.Close()
	common.Log("Server disconnected")
   close(s.messages)
}

func (s *Server) handle() {
   // defer wg.Done()
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
      if s.debug {
         fmt.Println(msg)
      }
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
      case protocol.ACTION_DISCONNECT:
         if s.clientExists(msg.Token) {
            msg.Action = protocol.ACTION_BROADCAST
            c := s.clients[msg.Token]
            // Just send message if the sender exists
            msg.Body = " -> " + c.username + " disconnected"
            fmt.Println(msg.Body)
            if s.clientExists(c.id) {
               for i := range s.clients {
                  // Don't send message to sending client
                  if c.id != s.clients[i].id {
                     _, err := s.connection.WriteToUDP(msg.Serialize(), s.clients[i].address)
                     checkError(err, "Server.handleReceivingMessages " + s.clients[i].username)
                  }
               }
            }
            delete(s.clients, msg.Token)
         }
      case protocol.ACTION_BROADCAST:
         c := s.clients[msg.Token]
			msg.Body = c.username + ": " + msg.Body
         // Just send message if the sender exists
         if s.clientExists(c.id) {
            for i := range s.clients {
               // Don't send message to sending client
               if c.id != s.clients[i].id {
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
      command = common.Readln()
      
      switch command {
      case "q":
         fallthrough
      case "quit":
         s.disconnect()
         break
      case "debug":
         s.debug = !s.debug
         if s.debug {
            common.Log("Debug is now on")
         } else {
            common.Log("Debug is now off")
         }
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
