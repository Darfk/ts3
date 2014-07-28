package ts3

import (
	"testing"
	"fmt"
	"strings"
)

func TestWalkClients(t *testing.T) {
	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}
	defer client.Close();

	client.Exec(Login(username, password))
	client.Exec(Use(1))
	client.WalkClients(func(idx int, client map[string]string) {
		t.Log(idx, client["client_nickname"])
	})
}

func testNotify(t *testing.T) {
	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}
	defer client.Close();

	notifyChan := make(chan string, 1)

	client.NotifyHandlerString(func(s string){
		notifyChan <- s
	})

	client.Exec(Login(username, password))
	client.Exec(Use(1))
	client.ExecString("servernotifyregister event=server")

	notification := ParseNotification(<-notifyChan)

	t.Log(notification)
	
}

func testConnection(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	client.Close()

}

func testParseResponse(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	defer client.Close()
	if err != nil {
		t.Error("Connection failed")
	}

	res, _ := client.ExecString("version")

	t.Log(ParseResponse(res))
}

func testParseError(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	_, tserr := client.ExecString("version")

	t.Log(ParseError(tserr))

	client.Close()
}

func testHelpers(t *testing.T) {

	var command Command
	var s fmt.Stringer

	command = Login("username", "password")
	s = &command
	t.Logf("%s", s)

	command = Kick([]string{"1", "2", "3"}, "")
	s = &command
	t.Logf("%s", s)
}

func testLogin(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	var res Response
	var tserr TSError

	client.Exec(Login(username, password))
	client.Exec(Use(1))
	res, tserr = client.Exec(ClientList())
	if tserr.id == 0 {
		for i := range res.Params {
			if nick, k := res.Params[i]["client_nickname"]; k && strings.Contains(nick, "Nathan") {
				client.Exec(Kick([]string{res.Params[i]["clid"]}, "GOLANG!"))
			}
		}
	}

	t.Log(res)
	t.Log(tserr)

	client.Close()
}