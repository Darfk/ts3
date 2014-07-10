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
	
	parseResponse(res)
	parseResponse(tserr)

	//t.Log(parsed, parsederr)

	client.Close()
}
