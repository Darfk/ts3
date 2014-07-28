package ts3

import (
	"strconv"
)

func Login(user, pass string) (c Command) {
	return Command{
		Command: "login",
		Params: map[string][]string {
			"client_login_name":[]string{user},
			"client_login_password":[]string{pass}},
	}
}

func Version() (c Command) {
	return Command{
		Command: "version",
	}
}


func Use(index int) (Command) {
	return Command{
		Command: "use",
		Params: map[string][]string {
			"sid":[]string{strconv.Itoa(index)}},
	}
}

func Poke(clients []string, reasonmsg string) (Command) {
	return Command{
		Command: "clientkick",
		Params: map[string][]string {
			"clid":clients,
			"reasonid":[]string{"5"},
			"reasonmsg":[]string{reasonmsg},
		},
	}
}

func Kick(clients []string, reasonmsg string) (Command) {
	return Command{
		Command: "clientkick",
		Params: map[string][]string {
			"clid":clients,
			"reasonid":[]string{"5"},
			"reasonmsg":[]string{reasonmsg},
		},
	}
}

func ClientList() (Command) {
	return Command{
		Command: "clientlist",
	}
}

func (client *Client) WalkClients(step func(int, map[string]string)) (err TSError) {

	var res Response

	res, err = client.Exec(ClientList())
	if err.id != 0 {
		return err
	}

	for i := range res.Params {
		step(i, res.Params[i])
	}

	return err
}
