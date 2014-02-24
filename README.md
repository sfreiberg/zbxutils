zbxutils
======

Zbxutils is a small/simple library for interacting with Zabbix agents and servers. At the moment it marshals and unmarshals data according to what Zabbix expects.

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

Example
=======

```
package main

import (
	"github.com/sfreiberg/zbxutils"
	"fmt"
)

func main() {
	zabbixCmd := []byte("agent.ping")
	payload := NewPayloadFromData(zabbixCmd)
}
```