package ts3

import (
	"strings"
	"log"
)

func Login(user, pass string) (c Command) {
	return Command{
		Command: "login",
		Params: map[string][]string {
			"client_login_name":[]string{user},
			"client_login_password":[]string{pass}},
	}
}

func Use(index string) (c Command) {
	return Command{
		Command: "use",
		Params: map[string][]string {
			"sid":[]string{index}},
	}
}

func Kick(clients []string) (c Command) {
	return Command{
		Command: "clientkick",
		Params: map[string][]string {
			"clid":clients,
			"reasonid":[]string{"5"},
		},
	}
}

func (c *Command) String() (s string) {
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
		log.Print(subResponses)
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
