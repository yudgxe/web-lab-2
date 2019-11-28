package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client struct {
	Addr string
}

func (cln *Client) StartClient() error {
	conn, err := net.Dial("tcp", cln.Addr)
	if err != nil {
		return err
	}
	fmt.Println("Connect to " + cln.Addr)

	go cln.handleConn(conn)
	cln.read(conn)

	return nil
}

func (cln *Client) read(conn net.Conn) error {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		writer.WriteString(str)
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	return nil
}

func (cln *Client) handleConn(conn net.Conn) error {
	defer func() {
		fmt.Println("Close connecting from server")
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(os.Stdout)
	scanr := bufio.NewScanner(reader)

	for {
		if !scanr.Scan() {
			if err := scanr.Err(); err != nil {
				fmt.Println(err)
				return err
			}
			break
		}
		fmt.Fprintln(writer, scanr.Text())
		writer.Flush()
	}
	return nil
}
