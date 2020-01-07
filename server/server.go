package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"web-lab-2/protector"
)

type Server struct {
	Addr            string
	IdleTimeout     time.Duration
	protectorServer *protector.SessionProtector
	MaxInit         int
	currentInit     int
}

func (srv *Server) StartServer() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":8080"
	}
	srv.currentInit = 0

	fmt.Println("Starting server on", addr)

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}
	defer listener.Close()

	for {

		newConn, err := listener.Accept()
		if srv.currentInit < srv.MaxInit {
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
			if srv.control(conn) != nil {
				fmt.Println("Error connection: ", err.Error())
			}
			srv.currentInit++
			go srv.handle(conn)
		} else {
			newConn.Close()
		}

	}
}

func (srv *Server) control(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	hashStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err, conn.RemoteAddr())
		return err
	}

	currentKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err, conn.RemoteAddr())
		return err
	}

	hashStr = hashStr[:len(hashStr)-1]
	currentKey = currentKey[:len(currentKey)-1]
	fmt.Println("from", conn.RemoteAddr(), hashStr)
	fmt.Println("from", conn.RemoteAddr(), currentKey)

	srv.protectorServer = protector.NewSessionProtector(hashStr)
	conn.Write([]byte(srv.protectorServer.Next_session_key(currentKey) + "\n"))
	return nil
}

func (srv *Server) handle(conn net.Conn) error {
	defer func() {
		fmt.Println("Close connecting from", conn.RemoteAddr())
		srv.currentInit--
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

		currentKey := scanr.Text()
		if !scanr.Scan() {
			if err := scanr.Err(); err != nil {
				fmt.Println(err, conn.RemoteAddr())
				return err
			}
			break
		}
		message := scanr.Text()
		fmt.Println("from", conn.RemoteAddr(), currentKey)
		fmt.Println("from", conn.RemoteAddr(), message)
		fmt.Println("Send to:", conn.RemoteAddr(), srv.protectorServer.Next_session_key(currentKey))
		writer.WriteString(srv.protectorServer.Next_session_key(currentKey) + "\n")
		writer.Flush()
	}
	return nil
}
