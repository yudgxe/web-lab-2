package main

import (
	"fmt"
	"os"
	"time"

	"./server"
)

func main() {
	fmt.Println("Enter the :port or ip:port")
	var addr string
	fmt.Fscan(os.Stdin, &addr)
	if len(addr) > 5 {

	} else {
		srv := server.Server{
			Addr:        addr,
			IdleTimeout: 10 * time.Second,
		}
		srv.StartServer()
	}
}
