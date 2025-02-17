package main

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

// Read/Write 테스트
func TestReadIntoBuffer(t *testing.T) {

	payload := make([]byte, 1<<24) // 16 MB
	//rand.Read 로 payload 에 난수 생성 할당
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	//listener 생성
	listener, err := net.Listen("tcp", "127.0.0.1:7138")
	if err != nil {
		t.Fatal(err)
	}

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

	//Dial TCP 데이터 스트림생성
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	//512KB 버퍼 생성
	buf := make([]byte, 1<<19) // 512 KB

	for {
		// 16MB/512KB 이므로 32번 데이터 Read함
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("read %d bytes", n)       // 몇바이트 표시
		t.Logf("read %x bytes", buf[:n]) // 데이터 표시
	}

	conn.Close()
}
