package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"web-lab-2/protector"
)

type Client struct {
	Addr            string
	protectorClient *protector.SessionProtector
	currentKey      string
}

func (cln *Client) handSheak(conn net.Conn, hash string) bool {
	cln.protectorClient = protector.NewSessionProtector(hash)
	cln.currentKey = protector.Get_session_key()
	clientKey := cln.protectorClient.Next_session_key(cln.currentKey)

	fmt.Fprintf(conn, hash+"\n")
	fmt.Fprintf(conn, cln.currentKey+"\n")

	serverKey, _ := bufio.NewReader(conn).ReadString('\n')

	if clientKey == serverKey {
		cln.currentKey = clientKey
		return true
	}

	return false
}

func (cln *Client) StartClient() error {
	conn, err := net.Dial("tcp", cln.Addr)
	if err != nil {
		return err
	}

	fmt.Println("Connect to " + cln.Addr)

	println(cln.handSheak(conn, protector.Get_hash_str()))

	go cln.handleConn(conn)
	cln.readConsole(conn)

	return nil
}

func (cln *Client) readConsole(conn net.Conn) error {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	var str string
	var err error

	for {
		str, err = reader.ReadString('\n')
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
