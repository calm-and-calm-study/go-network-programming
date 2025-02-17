package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	// MaxPayloadSize 10 MB
	MaxPayloadSize uint32 = 10 << 20
)

var ErrMaxPayloadSize = errors.New("max payload size exceeded")

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

func (m Binary) Bytes() []byte {
	return m
}

func (m Binary) String() string {
	return string(m)
}

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.BigEndian, BinaryType); err != nil {
		return 0, err
	}
	var n int64 = 1
	if err := binary.Write(w, binary.BigEndian, uint32(len(m))); err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write(m)

	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8

	// 1Byte Type
	if err := binary.Read(r, binary.BigEndian, &typ); err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType {
		return n, errors.New("invalid Binary")
	}

	var size uint32
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return n, err
	}

	n += 4
	// 최대 Payload 크기를 부여해서 조건 검사
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size)
	o, err := r.Read(*m)
	return n + int64(o), err
}

type String string

func (s String) Bytes() []byte {
	return []byte(s)
}

func (s String) String() string {
	return string(s)
}

func (s String) WriteTo(w io.Writer) (int64, error) {
	if err := binary.Write(w, binary.BigEndian, StringType); err != nil {
		return 0, err
	}

	var n int64 = 1
	if err := binary.Write(w, binary.BigEndian, n); err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write([]byte(s))

	return n + int64(o), err
}

func (s *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	if err := binary.Read(r, binary.BigEndian, &typ); err != nil {
		return 0, err
	}
	var n int64 = 1
	if typ != StringType {
		return n, errors.New("invalid String")
	}

	var size uint32
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return n, err
	}

	n += 4
	buf := make([]byte, size)
	o, err := r.Read(buf)
	if err != nil {
		return n, err
	}
	*s = String(buf[:o])
	return n + int64(o), nil
}

func decode(r io.Reader) (Payload, error) {
	var typ uint8
	if err := binary.Read(r, binary.BigEndian, &typ); err != nil {
		return nil, err
	}

	var payload Payload

	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unknown type")
	}

	_, err := payload.ReadFrom(
		io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}
	return payload, nil
}
