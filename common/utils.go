package common

import (
   "fmt"
   "bufio"
   "os"
   "os/signal"
)

func CheckError(err error) {
   if err != nil {
       fmt.Fprintf(os.Stderr, "Fatal error: %s \n", err.Error())
       os.Exit(1)
   }
}

func ByteArrayCompare(b1 *[]byte, b2 *[]byte) bool {
   if(len(*b1) != len(*b2)) {
      return false
   }

   for i := range *b1 {
      if (*b1)[i] != (*b2)[i] {
         return false
      }
   }

   return true
}

func Log(obj interface{}) {
   fmt.Println("*** ", obj)
}

func Readln() string {
   scanner := bufio.NewScanner(os.Stdin)
   if scanner.Scan() {
     return scanner.Text()
   }
   return ""
   // in := bufio.NewReader(os.Stdin)
   
   // line, _ := in.ReadString('\n')
   // fmt.Println(len(line))
   // if len(line) > 0 {
   //    return string(line[0:len(line) - 1])
   // }
   // return ""
}

type HandleExitFunc func() bool

func HandleExit(f HandleExitFunc) {
   c := make(chan os.Signal, 1)
   signal.Notify(c, os.Interrupt)
   go func(){
      for sig := range c {
         if sig == os.Interrupt {
            if f() {
               os.Exit(2)
            }
         }
      }
   }()
}