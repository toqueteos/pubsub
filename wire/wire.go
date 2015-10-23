// wire contains the building blocks
package wire

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ikanor/pubsub"
)

func ReadString(buf *bytes.Buffer, s string) (string, error) {
	c, err := buf.ReadByte()
	if err != nil {
		return "", err
	}
	if c != pubsub.WireString {
		return "", fmt.Errorf("expecting %x, got %x", pubsub.WireString, c)
	}
	var length int
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return "", err
	}
	if length == 0 {
		return "", nil
	}
	p := make([]byte, length)
	n, err := buf.Read(p)
	if err != nil {
		return "", err
	}
	if n != length {
		return "", fmt.Errorf("expecting %d bytes read, got %d", length, n)
	}
	return string(p), nil
}

func ReadSlice(buf *bytes.Buffer, payload []byte) ([]byte, error) {
	c, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if c != pubsub.WireSlice {
		return nil, fmt.Errorf("expecting %x, got %x", pubsub.WireSlice, c)
	}
	var length int
	if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	}
	p := make([]byte, length)
	n, err := buf.Read(p)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, fmt.Errorf("expecting %d bytes read, got %d", length, n)
	}
	return p, nil
}

func WriteString(buf *bytes.Buffer, s string) (err error) {
	err = buf.WriteByte(pubsub.WireString)
	if err != nil {
		return err
	}
	length := len(s)
	err = binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return err
	}
	n, err := buf.WriteString(s)
	if n != length || err != nil {
		return err
	}
	return nil
}

func WriteSlice(buf *bytes.Buffer, payload []byte) (err error) {
	err = buf.WriteByte(pubsub.WireSlice)
	if err != nil {
		return err
	}
	length := len(payload)
	err = binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return err
	}
	n, err := buf.Write(payload)
	if n != length || err != nil {
		return err
	}
	return nil
}
