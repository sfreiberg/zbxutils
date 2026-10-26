package zbxutils

import (
	"bytes"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeAgent is a minimal passive Zabbix agent listening on localhost. It
// answers each key from responses with a framed payload.
type fakeAgent struct {
	ln        net.Listener
	responses map[string]string

	mu   sync.Mutex
	keys []string
}

func newFakeAgent(t *testing.T, responses map[string]string) *fakeAgent {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	fa := &fakeAgent{ln: ln, responses: responses}
	t.Cleanup(func() { ln.Close() })
	go fa.serve()
	return fa
}

func (fa *fakeAgent) serve() {
	for {
		conn, err := fa.ln.Accept()
		if err != nil {
			return
		}
		go fa.handle(conn)
	}
}

func (fa *fakeAgent) handle(conn net.Conn) {
	defer conn.Close()
	// The client sends the bare key without a terminator and waits for the
	// response, so read whatever arrives in the first packet.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	key := strings.TrimSpace(string(buf[:n]))

	fa.mu.Lock()
	fa.keys = append(fa.keys, key)
	fa.mu.Unlock()

	resp, ok := fa.responses[key]
	if !ok {
		resp = string(NotSupported)
	}
	conn.Write(NewPayloadFromData([]byte(resp)).Bytes())
}

func (fa *fakeAgent) agent() *Agent {
	addr := fa.ln.Addr().(*net.TCPAddr)
	return NewAgentHostPort(addr.IP.String(), addr.Port)
}

func TestFakeAgentQueries(t *testing.T) {
	fa := newFakeAgent(t, map[string]string{
		PingKey:     "1",
		HostnameKey: "web01",
		VersionKey:  "7.0.0",
	})
	a := fa.agent()

	ok, err := a.Ping()
	if err != nil || !ok {
		t.Fatalf("Ping() = %v, %v; want true, nil", ok, err)
	}

	host, err := a.Hostname()
	if err != nil || host != "web01" {
		t.Fatalf("Hostname() = %q, %v; want web01, nil", host, err)
	}

	version, err := a.Version()
	if err != nil || version != "7.0.0" {
		t.Fatalf("Version() = %q, %v; want 7.0.0, nil", version, err)
	}

	fa.mu.Lock()
	defer fa.mu.Unlock()
	want := []string{PingKey, HostnameKey, VersionKey}
	if strings.Join(fa.keys, ",") != strings.Join(want, ",") {
		t.Fatalf("agent received keys %v, want %v", fa.keys, want)
	}
}

func TestFakeAgentPingFailure(t *testing.T) {
	fa := newFakeAgent(t, map[string]string{PingKey: "0"})
	ok, err := fa.agent().Ping()
	if err != nil || ok {
		t.Fatalf("Ping() = %v, %v; want false, nil", ok, err)
	}
}

func TestFakeAgentNotSupported(t *testing.T) {
	fa := newFakeAgent(t, nil)
	payload, err := fa.agent().Get("no.such.key")
	if err == nil || !strings.HasSuffix(err.Error(), " is not supported") {
		t.Fatalf("err = %v, want 'is not supported'", err)
	}
	if payload == nil || payload.Supported() {
		t.Fatal("payload should be returned and report NotSupported")
	}
}

func TestFakeAgentLargePayload(t *testing.T) {
	large := strings.Repeat(`{"{#FSNAME}":"/","{#FSTYPE}":"ext4"},`, 5000)
	fa := newFakeAgent(t, map[string]string{"vfs.fs.discovery": large})
	payload, err := fa.agent().Get("vfs.fs.discovery")
	if err != nil {
		t.Fatal(err)
	}
	if !payload.Valid() || string(payload.Data) != large {
		t.Fatalf("large payload was not read completely (%d of %d bytes)", len(payload.Data), len(large))
	}
}

func TestGetConnectionRefused(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	ln.Close()

	if _, err := NewAgentHostPort("127.0.0.1", addr.Port).GetWithTimeout(PingKey, time.Second); err == nil {
		t.Fatal("expected a connection error")
	}
}

func TestPayloadRoundTrip(t *testing.T) {
	for _, data := range []string{"", "1", strings.Repeat("x", 300), strings.Repeat("y", 70000)} {
		p := NewPayloadFromData([]byte(data))
		got, err := NewPayloadFromReader(bytes.NewReader(p.Bytes()))
		if err != nil {
			t.Fatalf("len %d: %v", len(data), err)
		}
		if !got.Valid() || string(got.Data) != data || got.DataLength() != uint64(len(data)) {
			t.Fatalf("len %d: round trip mismatch", len(data))
		}
	}
}

func TestShortHeader(t *testing.T) {
	if _, err := NewPayloadFromReader(bytes.NewReader([]byte("ZBX"))); err == nil {
		t.Fatal("expected an error for a truncated header")
	}
}
