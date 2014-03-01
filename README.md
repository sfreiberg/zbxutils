zbxutils
======

Zbxutils is a simple library for interacting with Zabbix agents and servers. At the moment it marshals and unmarshals data according to the Zabbix protocol and includes the ability to query zabbix agents.

License
=======

Zagent is licensed under the MIT license.

Installation
============
`go get github.com/sfreiberg/zbxutils`

Documentation
=============
[GoDoc](http://godoc.org/github.com/sfreiberg/zbxutils)

[Zabbix Protocol](https://www.zabbix.com/documentation/2.2/manual/appendix/items/activepassive)

[Zabbix Items Supported by Platform](https://www.zabbix.com/documentation/2.2/manual/appendix/items/supported_by_platform)

Example
=======

```
package main

import (
	"github.com/sfreiberg/zbxutils"

	"fmt"
	"log"
)

func main() {
	// Talk to an agent on localhost port 10050
	agent := zbxutils.NewAgentHostPort("localhost", 10050)

	// Cconnect to agent and grab version
	ver, err := agent.Version()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to agent and grab hostname
	// https://www.zabbix.com/documentation/2.2/manual/appendix/items/supported_by_platform
	hostname, err := agent.Get("agent.hostname")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Version:", ver)
	fmt.Println("Hostname:", hostname)
}

```