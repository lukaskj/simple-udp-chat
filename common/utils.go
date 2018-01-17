package common

import (
	"fmt"
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
	fmt.Println("*** LOG ", obj)
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