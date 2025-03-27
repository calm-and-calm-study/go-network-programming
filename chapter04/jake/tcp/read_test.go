package tcp

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {
	// 1을 24bit 왼쪽으로 이동시킨다. == 16MB 의 슬라이스 채널을 생성
	// 2^24 = 16777216 (약 16MB)
	payload := make([]byte, 1<<24)
	// payload 에 랜덤으로 데이터를 생성
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// 송신 고루틴
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	// 수신고루틴
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// 512 KB
	buf := make([]byte, 1<<19)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("read %d bytes", n)
	}

	conn.Close()
}
