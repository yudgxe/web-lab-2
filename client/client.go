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

func (cln *Client) StartClient() error {
	conn, err := net.Dial("tcp", cln.Addr)
	if err != nil {
		return err
	}
	fmt.Println("Connect to " + cln.Addr)

	hashStr := protector.Get_hash_str()
	fmt.Println("Start hash:", hashStr)
	cln.currentKey = protector.Get_session_key()
	fmt.Println("Start key:", cln.currentKey)
	cln.protectorClient = protector.NewSessionProtector(hashStr)

	conn.Write([]byte(hashStr + "\n" + cln.currentKey + "\n"))

	reader := bufio.NewReader(conn)
	currentKey, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err, conn.RemoteAddr())
		return err
	}
	fmt.Println("Key form server:", currentKey[:len(currentKey)-1])
	cln.currentKey = cln.protectorClient.Next_session_key(cln.currentKey)
	if currentKey[:len(currentKey)-1] != cln.currentKey {
		fmt.Println(":(")
		return nil
	} else {
		fmt.Println(":)")
	}
	go cln.send(conn)
	cln.get(conn)

	return nil
}

func (cln *Client) send(conn net.Conn) error {
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
		cln.currentKey = cln.protectorClient.Next_session_key(cln.currentKey)
		fmt.Println("Send key:", cln.currentKey)
		writer.WriteString(cln.currentKey + "\n")
		writer.WriteString(str)
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return nil
}

func (cln *Client) get(conn net.Conn) error {
	defer func() {
		fmt.Println("Close connecting from server")
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(os.Stdout)
	scanr := bufio.NewScanner(reader)

	for {
		if !scanr.Scan() {
			if err := scanr.Err(); err != nil {
				fmt.Println(err)
				return err
			}
			break
		}
		currentKey := scanr.Text()
		fmt.Println("Get key: ", currentKey)
		cln.currentKey = cln.protectorClient.Next_session_key(cln.currentKey)
		if currentKey != cln.currentKey {
			fmt.Println(":(")
			return nil
		} else {
			fmt.Println("Current key:", cln.currentKey)
		}
	}
	return nil
}
