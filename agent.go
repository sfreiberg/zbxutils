package zbxutils

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	DefaultAgentPort    = 10050                           // Default port for contacting the Zabbix Agent
	DefaultAgentHost    = "localhost"                     // Default host for contacting the Zabbix Agent
	DefaultAgentTimeout = time.Duration(30 * time.Second) // Default timeout when contacting the Zabbix Agent

	PingKey     = "agent.ping"     // The key to ping the remote zabbix agent
	HostnameKey = "agent.hostname" // The key to retrieve the hostname of the remote zabbix agent
	VersionKey  = "agent.version"  // The key to retrieve the version of the remote zabbix agent
)

// Agent represents a remote zabbix agent.
type Agent struct {
	host string
	port int
}

// Creates a new Agent on localhost with a default port of 10050.
func NewAgent() *Agent {
	return &Agent{host: DefaultAgentHost, port: DefaultAgentPort}
}

// Creates a new Agent with a custom host and default port of 10050.
func NewAgentHost(host string) *Agent {
	return &Agent{host: host, port: DefaultAgentPort}
}

// Creates a new Agent with a custom host and port.
func NewAgentHostPort(host string, port int) *Agent {
	return &Agent{host: host, port: port}
}

// Run the check (key) against the Zabbix agent with the default timeout
func (a *Agent) Get(key string) (*Payload, error) {
	return a.GetWithTimeout(key, DefaultAgentTimeout)
}

// Run the check (key) against the Zabbix agent with the specified timeout
func (a *Agent) GetWithTimeout(key string, timeout time.Duration) (*Payload, error) {
	conn, err := net.DialTimeout("tcp", a.hostPort(), timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, key)
	if err != nil {
		return nil, err
	}

	payload, err := NewPayloadFromReader(conn)
	if err != nil {
		return nil, err
	}

	if payload.NotSupported() {
		return payload, fmt.Errorf("%s is not supported", key)
	}

	return payload, nil
}

// Call agent.ping and verifies it returns "1" for success.
func (a *Agent) Ping() (bool, error) {
	payload, err := a.Get(PingKey)
	if err != nil {
		return false, err
	}

	success := bytes.Equal([]byte{'1'}, payload.Data)
	return success, nil
}

// Calls agent.hostname on the zabbix agent and returns the result.
func (a *Agent) Hostname() (string, error) {
	payload, err := a.Get(HostnameKey)
	return string(payload.Data), err
}

// Calls agent.version on the zabbix host and returns the result.
func (a *Agent) Version() (string, error) {
	payload, err := a.Get(VersionKey)
	return string(payload.Data), err
}

// Join host and port
func (a *Agent) hostPort() string {
	return net.JoinHostPort(a.host, strconv.Itoa(a.port))
}
