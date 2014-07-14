package ts3

import (
	"testing"
)

func testConnection(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	client.Close()

}

func testRawCommand(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	res, tserr := client.rawCommand("version")

	t.Log(res, tserr)

	client.Close()
}

func TestParseResponse(t *testing.T) {

	client, err := NewClient("teamspeak.darfk.net:10011")
	if err != nil {
		t.Error("Connection failed")
	}

	res, tserr := client.rawCommand("version")

	t.Log(parseResponse(res))
	t.Log(parseResponse(tserr))

	//t.Log(parsed, parsederr)

	client.Close()
}

func testHelpers(t *testing.T) {

	var command Command

	command = Login("username", "password")
	t.Log(command.String())

	command = Kick([]string{"1", "2", "3"})
	t.Log(command.String())
}