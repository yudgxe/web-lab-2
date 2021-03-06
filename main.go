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
	if len(addr) > 5 {
		cln := client.Client{
			Addr: addr,
		}
		err := cln.StartClient()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		var max int
		fmt.Println("Enter max:")
		fmt.Fscan(os.Stdin, &max)
		srv := server.Server{
			Addr:        addr,
			IdleTimeout: 3 * time.Minute,
			MaxInit:     max,
		}
		err := srv.StartServer()
		if err != nil {
			fmt.Println(err)
		}
	}
}
