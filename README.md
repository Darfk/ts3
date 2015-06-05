Teamspeak 3 Server Query library for Golang
---

Package ts3 implements the Teamspeak 3 ServerQuery protocol described [here](http://media.teamspeak.com/ts3_literature/TeamSpeak%203%20Server%20Query%20Manual.pdf)

Features
---

- Notification Handling
- Some helper functions

Example Usage
---

See the test file for a more comprehensive example

    package main

    import (
        "darfk/ts3"
        "log"
    )

    func main() {

        client, err := ts3.NewClient("teamspeak.darfk.net:10011")

        username := "test"

        response, err := client.Exec(ts3.Login(username, "xWUkRRlM"))
        log.Printf("logged in as %s\n", username)

        response, err = client.Exec(ts3.Version())
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("version %q\n", response)

        client.Close()

    }

Licence
---

MIT
