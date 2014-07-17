package ts3

import (
	"testing"
	"fmt"
	"strings"
)

func testConnection(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	client.Close()

}

func testParseResponse(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	res, _ := client.ExecString("version")

	t.Log(parseResponse(res))

	client.Close()
}

func testParseError(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	_, tserr := client.ExecString("version")

	t.Log(parseError(tserr))

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

func TestLogin(t *testing.T) {

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