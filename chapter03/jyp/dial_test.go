package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	// 리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:7138")
	// 에러 발생 시 페이탈 에러
	if err != nil {
		t.Fatal(err)
	}

	// 구조체 만듬듬
	done := make(chan struct{})

	// 각각의 리스너 각각마다 ESTABLISHED 되면 conn 에
	go func() {
		// connection close 후 done 초기화
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) {
				defer func() {
					c.Close()
					done <- struct{}{}
				}()
				// 수신 버퍼 만듬듬
				buf := make([]byte, 1024)
				for {
					// connection 에 read 할 버퍼를 할당
					n, err := c.Read(buf)
					if err != nil {
						// EOF = 아무것도 데이터 없음이 아니면 에러 발생
						if err != io.EOF {
							t.Error(err)
						}
						return
					}

					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	// listener 로 connection 연결
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	// connection close
	conn.Close()
	<-done
	// 리스너 close
	listener.Close()
	<-done
}
