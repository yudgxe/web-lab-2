package main

import (
	"fmt"
	"os"
	"time"

	"web-lab-2/client"
	"web-lab-2/server"
)

func main() {
	fmt.Println("Enter the :port or ip:port")
	var addr string
	fmt.Fscan(os.Stdin, &addr)
	if addr == "client" {
		cln := client.Client{
			Addr: "127.0.0.1:8000",
		}
		err := cln.StartClient()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		srv := server.Server{
			Addr:        addr,
			IdleTimeout: 3 * time.Minute,
		}
		err := srv.StartServer()
		if err != nil {
			fmt.Println(err)
		}
	}
}
