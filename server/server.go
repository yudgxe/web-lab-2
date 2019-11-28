package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type Server struct {
	Addr        string
	IdleTimeout time.Duration
}

func (srv *Server) StartServer() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":8080"
	}

	fmt.Println("Starting server on", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		newConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		fmt.Println("Accepted connection from: ", newConn.RemoteAddr())

		conn := &conn{
			Conn:        newConn,
			IdleTimeout: srv.IdleTimeout,
		}

		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))
		go srv.handle(conn)
	}
}

func (srv *Server) handle(conn net.Conn) error {
	defer func() {
		fmt.Println("Close connecting from", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	scanr := bufio.NewScanner(reader)

	for {
		if !scanr.Scan() {
			if err := scanr.Err(); err != nil {
				fmt.Println(err, conn.RemoteAddr())
				return err
			}
			break
		}
		fmt.Println(scanr.Text())
		writer.WriteString(strings.ToUpper(scanr.Text()) + "\n")
		writer.Flush()
	}
	return nil
}
