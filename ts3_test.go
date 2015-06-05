package ts3

import (
	"testing"
)

type config struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Server   int    `json:"server"`
}

func TestTS3(t *testing.T) {

	var err error

	config := config{
		"teamspeak.darfk.net:10011",
		"test",
		"xWUkRRlM",
		1,
	}

	client, err := NewClient(config.Address)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("connected to %s", config.Address)

	var response Response

	response, err = client.Exec(Login(config.Username, config.Password))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("logged in as %s", config.Username)

	response, err = client.Exec(Use(config.Server))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("using server %d", config.Server)

	response, err = client.Exec(Version())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("version %q", response)

	t.Logf("doing something we're not allowed to do")
	response, err = client.Exec(Command{
		Command: "serverlist",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	t.Logf("%s", err)

	// Let's see if Nathan's online >:)

	var nathanIsOnline bool = false

	err = client.WalkClients(func(idx int, client map[string]string) {
		if nick, ok := client["client_nickname"]; ok && nick == "Nathan" {
			nathanIsOnline = true
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	if nathanIsOnline {
		t.Log("Nathan is online!")
	} else {
		t.Log("Nathan must be asleep!")
	}

}
