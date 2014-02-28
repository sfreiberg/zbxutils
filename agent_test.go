package zbxutils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Defaults to testing localhost and port 10050. If you want to specify the
// the host and port export ZABBIX_HOST and ZABBIX_PORT respectively.
func TestAgentPing(t *testing.T) {
	agent := newAgent(t)

	pingable, err := agent.Ping()
	if err != nil {
		t.Fatal("Ping test failed with error:", err)
	}

	if !pingable {
		t.Fatal("Zabbix agent wasn't pingable")
	}

	fmt.Println("Ping test: PASS")
}

func TestAgentHostname(t *testing.T) {
	agent := newAgent(t)

	hostname, err := agent.Hostname()
	if err != nil {
		t.Fatal("Hostname test failed with error:", err)
	}

	if hostname == "" {
		t.Fatal("Zabbix hostname was empty")
	}

	fmt.Printf("Hostname test: PASS (%s)\n", hostname)
}

func TestAgentVersion(t *testing.T) {
	agent := newAgent(t)

	version, err := agent.Version()
	if err != nil {
		t.Fatal("Version test failed with error:", err)
	}

	if version == "" {
		t.Fatal("Zabbix version was empty")
	}

	fmt.Printf("Hostname test: PASS (%s)\n", version)
}

func TestAgentUnsupported(t *testing.T) {
	agent := newAgent(t)
	fakeKey := "Supercalifragilisticexpialidocious"

	payload, err := agent.Get(fakeKey)
	if err == nil {
		t.Fatal("An error isn't thrown when calling an unknown key")
	}

	if err != nil && !strings.HasSuffix(err.Error(), " is not supported") {
		t.Fatal(err)
	}

	if payload.Supported() {
		t.Fatal("Response.Supported() reports true and should be false")
	}

	fmt.Printf("Unsupported test: PASS (%s)\n", fakeKey)
}

func newAgent(t *testing.T) *Agent {
	var err error
	host := DefaultAgentHost
	port := DefaultAgentPort

	if os.Getenv("ZABBIX_HOST") != "" {
		host = os.Getenv("ZABBIX_HOST")
	}

	if os.Getenv("ZABBIX_PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("ZABBIX_PORT"))
		if err != nil {
			t.Fatal(err)
		}
	}

	return NewAgentHostPort(host, port)
}
