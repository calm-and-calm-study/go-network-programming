package tcp

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestPayloads(t *testing.T) {
	b1 := Binary("Clear is better than clever.")
	b2 := Binary("Don't panic.")

	s1 := String("Errors are values.")

	// 테스트하고자 하는 데이터 목록
	payloads := []Payload{&b1, &s1, &b2}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for _, p := range payloads {
			// 각 데이터 타입에 맞는 데이터를 Write
			_, err := p.WriteTo(conn)
			if err != nil {
				t.Error(err)
				break
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < len(payloads); i++ {
		actual, err := decode(conn)
		if err != nil {
			t.Fatal(err)
		}

		if expected := payloads[i]; !reflect.DeepEqual(actual, expected) {
			t.Errorf("value mismatch: %v != %v", expected, actual)
			continue
		}

		t.Logf("[%T] %[1]q", actual)
	}
}

func TestMaxPayloadSize(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := buf.WriteByte(BinaryType); err != nil {
		t.Fatal(err)
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(1<<30)); err != nil {
		t.Fatal(err)
	}

	var b Binary
	_, err := b.ReadFrom(buf)
	if err != nil {
		t.Fatalf("expected ErrMaxPayloadSize; actual: %v", err)
	}
}
