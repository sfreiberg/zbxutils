package zbxutils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	Header         = []byte("ZBXD\x01")
	NotSupported   = "ZBX_NOTSUPPORTED"
	InvalidHeader  = errors.New("Invalid Header")
	InvalidDataLen = errors.New("Invalid DataLen")
	InvalidData    = errors.New("Invalid Data")
)

// Payload is the structure that Zabbix uses to communicate.
type Payload struct {
	Header  []byte
	DataLen []byte
	Data    []byte
}

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

	LengthInBinary(p.DataLen, uint64(len(data)))

	return p
}

// Returns the payload as a slice of bytes
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

// Returns true if the header is valid
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

// Takes the size, converts it to binary and puts it in b
func LengthInBinary(b []byte, size uint64) {
	binary.LittleEndian.PutUint64(b, size)
}
