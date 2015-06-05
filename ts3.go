// Package ts3 implements the Teamspeak 3 ServerQuery protocol
// described here http://media.teamspeak.com/ts3_literature/TeamSpeak%203%20Server%20Query%20Manual.pdf
package ts3

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	conn                net.Conn
	scan                *bufio.Scanner
	line                chan string
	notify              chan string
	err                 chan string
	res                 string
	notifyHandler       func(Notification)
	notifyHandlerString func(string)
}

type Command struct {
	Command string
	Params  map[string][]string
	Flags   []string
}

type Response struct {
	Params []map[string]string
}

type Notification struct {
	Type   string
	Params []map[string]string
}

type TSError struct {
	id    int
	msg   string
	extra string
}

func (e TSError) Error() string {
	return fmt.Sprintf("ts3: %s (%d) %s", e.msg, e.id, e.extra)
}

func NewClient(address string) (client *Client, err error) {

	client = new(Client)
	client.conn, err = net.Dial("tcp", address)
	if err != nil {
		return
	}

	client.line = make(chan string)

	client.scan = bufio.NewScanner(client.conn)
	client.scan.Split(scanTS3Lines)
	go func() {
		for {
			client.scan.Scan()
			client.line <- client.scan.Text()
		}
	}()

	// Discard first 2 lines
	<-client.line
	<-client.line

	client.err = make(chan string)
	client.notify = make(chan string)

	// Spin up the notify handler
	go func() {
		for {
			notification := <-client.notify
			if client.notifyHandler != nil {
				client.notifyHandler(ParseNotification(notification))
			} else if client.notifyHandlerString != nil {
				client.notifyHandlerString(notification)
			} else {
				// discard
			}
		}
	}()

	// Spin up the line handler
	go func() {
		for {
			line := <-client.line
			if strings.Index(line, "error") == 0 {
				client.err <- line
			} else if strings.Index(line, "notify") == 0 {
				client.notify <- line
			} else {
				client.res = line
			}
		}
	}()

	return
}

func (client *Client) NotifyHandler(handler func(Notification)) {
	client.notifyHandlerString = nil
	client.notifyHandler = handler
}

func (client *Client) NotifyHandlerString(handler func(string)) {
	client.notifyHandlerString = handler
	client.notifyHandler = nil
}

func (client *Client) RemoveNotifyHandler() {
	client.notifyHandlerString = nil
	client.notifyHandler = nil
}

func (client *Client) Exec(command Command) (Response, error) {
	fmt.Fprintf(client.conn, "%s\n\r", command)
	err := <-client.err
	res := client.res
	client.res = ""
	return ParseResponse(res), ParseError(err)
}

func (client *Client) ExecString(command string) (string, string) {
	fmt.Fprintf(client.conn, "%s\n\r", command)
	err := <-client.err
	res := client.res
	client.res = ""
	return res, err
}

func (client *Client) Close() error {
	return client.conn.Close()
}

// This function is almost exactly like bufio.ScanLines except the \r\n are in opposite positions
func scanTS3Lines(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

func ParseResponse(s string) (r Response) {
	r = Response{}
	subResponses := strings.Split(s, "|")
	for i := range subResponses {
		r.Params = append(r.Params, make(map[string]string))
		kvPairs := strings.Split(subResponses[i], " ")
		for ii := range kvPairs {
			kvPair := strings.SplitN(kvPairs[ii], "=", 2)
			if len(kvPair) > 1 {
				r.Params[i][kvPair[0]] = Unescape(kvPair[1])
			} else {
				r.Params[i][kvPair[0]] = ""
			}
		}
	}
	return
}

func ParseNotification(s string) (n Notification) {
	n = Notification{}
	typeBody := strings.SplitN(s, " ", 2)
	n.Type = typeBody[0]
	subResponses := strings.Split(typeBody[1], "|")
	for i := range subResponses {
		n.Params = append(n.Params, make(map[string]string))
		kvPairs := strings.Split(subResponses[i], " ")
		for ii := range kvPairs {
			kvPair := strings.SplitN(kvPairs[ii], "=", 2)
			if len(kvPair) > 1 {
				n.Params[i][kvPair[0]] = Unescape(kvPair[1])
			} else {
				n.Params[i][kvPair[0]] = ""
			}
		}
	}
	return
}

func ParseError(s string) error {
	e := TSError{}
	kvPairs := strings.Split(s, " ")
	for i := range kvPairs {
		kvPair := strings.SplitN(kvPairs[i], "=", 2)
		if len(kvPair) > 1 {
			if kvPair[0] == "id" {
				id, err := strconv.ParseInt(kvPair[1], 10, 32)
				if err != nil {
					e.id = -1
				} else {
					e.id = int(id)
				}
			} else if kvPair[0] == "msg" {
				e.msg = Unescape(kvPair[1])
			} else if kvPair[0] == "extra_msg" {
				e.extra = Unescape(kvPair[1])
			}
		} else {
			continue
		}
	}

	if e.id != 0 {
		return e
	}

	return nil
}

// This makes the Command struct satisfy the fmt.Stringer interface
func (c Command) String() (s string) {
	var params, flags []string
	for k, v := range c.Params {
		if len(v) > 1 {
			var subParams []string
			for _, vv := range v {
				subParams = append(subParams, k+"="+Escape(vv))
			}
			params = append(params, strings.Join(subParams, "|"))
		} else if len(v) == 1 {
			params = append(params, strings.Join([]string{k, Escape(v[0])}, "="))
		} else {
			params = append(params, k)
		}
	}
	s = strings.Join([]string{c.Command, strings.Join(params, " "), strings.Join(flags, " ")}, " ")
	return
}
