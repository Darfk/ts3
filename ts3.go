package ts3

import (
	"log"
	_"time"
	"bufio"
	"net"
	"fmt"
	"bytes"
	"strings"
)

type Client struct {
	conn net.Conn
	scan *bufio.Scanner
	line chan string
	notify chan string
	err chan string
	res chan string
}

type TSError struct {
	id int
	msg string
}

func NewClient(address string) (client *Client, err error) {

	client = new(Client) 
	client.conn, err = net.Dial("tcp", address)
	if err != nil {
		return
	}

	client.line = make(chan string)

	client.scan = bufio.NewScanner(client.conn)
	client.scan.Split(ScanLines)
	go func() {
		for {
			client.scan.Scan()
			client.line <- client.scan.Text()
		}
	}()

	// Discard first 2 lines
	<- client.line
	<- client.line

	client.err = make(chan string)
	client.res = make(chan string)
	client.notify = make(chan string)

	go func() {
		for {
			line := <- client.line
			if strings.Index(line, "error") == 0 {
				client.err <- line
			} else if strings.Index(line, "notify") == 0 {
				client.notify <- line
			} else {
				client.res <- line
			}
		}
	}()

	return
}

func (client *Client) rawCommand(command string) (string, string) {
	fmt.Fprintf(client.conn, "%s\n\r", command)
	return <- client.res, <- client.err 
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func parseResponse(res string) {
	for i := range strings.Split(res, "|") {
		
	}
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
 		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n\r")); i >= 0 {
		return i + 2, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}