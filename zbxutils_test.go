package zbxutils

import (
	"bytes"
	"io"
	"testing"
)

var (
	// zbxBytes = Header + size (in binary) + "2.2.2"
	validBytes   = []byte{90, 66, 88, 68, 1, 5, 0, 0, 0, 0, 0, 0, 0, 50, 46, 50, 46, 50}
	validData    = "2.2.2"
	validDataLen = 5

	invalidHeader  = []byte{89, 66, 88, 68, 1, 5, 0, 0, 0, 0, 0, 0, 0, 50, 46, 50, 46, 50}
	invalidDataLen = []byte{90, 66, 88, 68, 1, 6, 0, 0, 0, 0, 0, 0, 0, 50, 46, 50, 46, 50}
)

func TestNewPayloadFromReader(t *testing.T) {
	buffer := bytes.NewBuffer(validBytes)

	payload, err := NewPayloadFromReader(buffer)
	if err != nil {
		t.Fatal("Unable to create Payload:", err)
	}
	if payload.Valid() == false {
		t.Fatal("Payload is invalid")
	}

	if !bytes.Equal(validBytes, payload.Bytes()) {
		t.Fatal("payload.Bytes != the raw data")
	}
}

func TestNewPayloadFromData(t *testing.T) {
	payload := NewPayloadFromData([]byte(validData))

	if !bytes.Equal(validBytes, payload.Bytes()) {
		t.Fatal("payload.Bytes != the raw data")
	}

	if !bytes.Equal(payload.Data, []byte(validData)) {
		t.Fatal("Converted data doesn't match input data")
	}

	if payload.Valid() == false {
		t.Fatal("Payload is invalid")
	}
}

func TestInvalidHeader(t *testing.T) {
	buffer := bytes.NewBuffer(invalidHeader)
	_, err := NewPayloadFromReader(buffer)
	if err != InvalidHeader {
		t.Fatal("Didn't catch the invalid header")
	}
}

func TestInvalidDataLen(t *testing.T) {
	buffer := bytes.NewBuffer(invalidDataLen)
	_, err := NewPayloadFromReader(buffer)
	if err != io.ErrUnexpectedEOF {
		t.Fatal("Didn't catch the invalid data length")
	}
}
