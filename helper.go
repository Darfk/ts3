package ts3

import (
	"strings"
	"strconv"
	_"log"
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


func Use(index int) (c Command) {
	return Command{
		Command: "use",
		Params: map[string][]string {
			"sid":[]string{strconv.Itoa(index)}},
	}
}

func Poke(clients []string, reasonmsg string) (c Command) {
	return Command{
		Command: "clientkick",
		Params: map[string][]string {
			"clid":clients,
			"reasonid":[]string{"5"},
			"reasonmsg":[]string{reasonmsg},
		},
	}
}

func Kick(clients []string, reasonmsg string) (c Command) {
	return Command{
		Command: "clientkick",
		Params: map[string][]string {
			"clid":clients,
			"reasonid":[]string{"5"},
			"reasonmsg":[]string{reasonmsg},
		},
	}
}

func ClientList() (c Command) {
	return Command{
		Command: "clientlist",
	}
}

func (c Command) String() (s string) {
	var params, flags []string
	for k, v := range c.Params {
		if len(v) > 1 {
			var subParams []string
			for _, vv := range v {
				subParams = append(subParams, k + "=" + vv)
			}
			params = append(params, strings.Join(subParams, "|"))
		}else if len(v) == 1 {
			params = append(params, strings.Join([]string{k, v[0]}, "="))
		}else{
			params = append(params, k)
		}
	}
	s = strings.Join([]string{c.Command, strings.Join(params, " "), strings.Join(flags, " ")}, " ")
	return
}

func parseResponse(s string) (r Response) {
	r = Response{}
	subResponses := strings.Split(s, "|")
	for i := range subResponses {
		r.Params = append(r.Params, make(map[string]string))
		kvPairs := strings.Split(subResponses[i], " ")
		for ii := range kvPairs {
			kvPair := strings.SplitN(kvPairs[ii], "=", 2)
			if len(kvPair) > 1 {
				r.Params[i][kvPair[0]] = kvPair[1]
			}else{
				r.Params[i][kvPair[0]] = ""
			}
		}
	}
	return
}

func parseError(s string) (e TSError) {
	e = TSError{}
	kvPairs := strings.Split(s, " ")
	for i := range kvPairs {
		kvPair := strings.SplitN(kvPairs[i], "=", 2)
		if len(kvPair) > 1 {
			if kvPair[0] == "id" {
				id, err := strconv.ParseInt(kvPair[1], 10, 32)
				if err != nil {
					e.id = -1
				}else{
					e.id = int(id)
				}
			}else if kvPair[0] == "msg" {
				e.msg = kvPair[1]
			}
		}else{
			continue;
		}
	}
	return
}
