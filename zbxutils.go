/*
Zbxutils is a simple library for interacting with Zabbix agents and
servers. At the moment it marshals and unmarshals data according to the
Zabbix protocol and includes the ability to query zabbix agents.
*/
package zbxutils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	Header         = []byte("ZBXD\x01")            // Standard Zabbix header
	NotSupported   = []byte("ZBX_NOTSUPPORTED")    // Not Supported response
	InvalidHeader  = errors.New("Invalid Header")  // Invalid header error
	InvalidDataLen = errors.New("Invalid DataLen") // Invalid datalen error
	InvalidData    = errors.New("Invalid Data")    // Invalid data error
)

// Payload is the structure that Zabbix uses to communicate.
type Payload struct {
	Header  []byte
	DataLen []byte
	Data    []byte
}

// Create a payload from any io.Reader.
func NewPayloadFromReader(reader io.Reader) (*Payload, error) {
	payload := &Payload{
		Header:  make([]byte, 5),
		DataLen: make([]byte, 8),
	}

	// Read and validate the header
	n, err := reader.Read(payload.Header)
	if err != nil {
		return nil, err
	}
	if n != 5 {
		return nil, InvalidHeader
	}
	if !payload.ValidHeader() {
		return nil, InvalidHeader
	}

	// Read and validate the DataLen
	n, err = reader.Read(payload.DataLen)
	if err != nil {
		return nil, err
	}

	if n != 8 {
		return nil, InvalidDataLen
	}

	payload.Data = make([]byte, payload.DataLength())
	n, err = reader.Read(payload.Data)
	if err != nil {
		return nil, err
	}

	if uint64(n) != payload.DataLength() {
		return nil, InvalidDataLen
	}

	if !payload.ValidData() {
		return nil, InvalidData
	}

	return payload, nil
}

// Create a new Payload from a []byte.
func NewPayloadFromData(data []byte) *Payload {
	p := &Payload{
		Header:  Header,
		DataLen: make([]byte, 8),
		Data:    data,
	}

	lengthInBinary(p.DataLen, uint64(len(data)))

	return p
}

// Returns the payload as a slice of bytes.
func (p *Payload) Bytes() []byte {
	b := []byte{}
	b = append(b, p.Header...)
	b = append(b, p.DataLen...)
	b = append(b, p.Data...)

	return b
}

// Check to see if this Payload is valid. Basically validates
// that the header is correct and that the DataLen matches
// the actual length of the Data.
func (p *Payload) Valid() bool {
	if !p.ValidHeader() || !p.ValidData() {
		return false
	}

	return true
}

// Returns true if the header is valid.
func (p *Payload) ValidHeader() bool {
	return bytes.Equal(p.Header, Header)
}

// Returns true if Data is valid. In this case it's only
// checking the actual length versus the expected length.
func (p *Payload) ValidData() bool {
	return binary.LittleEndian.Uint64(p.DataLen) == uint64(len(p.Data))
}

// Returns the expected length of DataLen
func (p *Payload) DataLength() uint64 {
	return binary.LittleEndian.Uint64(p.DataLen)
}

// Returns true if the request is supported.
func (p *Payload) Supported() bool {
	return !p.NotSupported()
}

// Returns true if the request isn't supported.
func (p *Payload) NotSupported() bool {
	return bytes.Equal(p.Data, NotSupported)
}

// Takes the size, converts it to binary and puts it in b.
func lengthInBinary(b []byte, size uint64) {
	binary.LittleEndian.PutUint64(b, size)
}
