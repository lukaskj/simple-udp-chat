package main

import (
   "net"
   "sync"
   "time"
   "../common"
)

type Client struct {
   id              string
   username        string
   lastMessageTime time.Time
   address         *net.UDPAddr
}

// MAIN *********

var wg sync.WaitGroup

func main() {
   var port string = ":8080"

   var server Server
   defer server.disconnect()
   common.HandleExit(func () bool {
      server.disconnect()
      return true
   })
   wg.Add(1)
   server.start(port)

   // common.CheckError(err)

   

   wg.Wait()
}


