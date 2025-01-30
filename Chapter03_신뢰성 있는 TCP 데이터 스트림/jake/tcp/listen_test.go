package tcp

import (
	"errors"
	"io"
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	// ip 주소를 생략하면 리스너는 시스템상의 모든 유니캐스트와 애니캐스트 ip 주소에 바인딩함
	// port 를 0으로 설정하거나 비워두면 Go가 리스너에 무작위로 포트 번호를 할당함
	// IPv4 : tcp4 / IPv6 : tcp6
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	// 리스너를 종료하는데 실패하면 메모리 누수가 발생하거나 코드상에서 리스너의 Accept 메서드가 무한정 블로킹되며 데드락이 발생할 수 있음
	// 리스너를 종료하는 즉시 Accept 메서드의 블로킹이 해제됨
	defer func() {
		if err := listener.Close(); err != nil {
			t.Logf("unexpected close err : %v", err)
		}
	}()

	t.Logf("bound to %q", listener.Addr())
}

func TestDial(t *testing.T) {
	// 리스너를 시작해서 해당 포트로 들어오는 tcp 연결을 받는다.
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Log(err)
				return
			}

			// 아래 고루틴으로 로직을 처리하는 내용이 handler 이다.
			// 각각의 요청이 들어오면 커넥션 구조체를 handler 에 넘겨준다.
			go func(c net.Conn) {
				defer func() {
					// close 함수를 호출한다.
					if err := c.Close(); err != nil {
						t.Logf("unexpected close err : %v", err)
					}
					done <- struct{}{}
				}()

				buf := make([]byte, 1024)
				for {
					// 연결이 끊어지면 io.EOF 에러를 반환하면서 종료된다. 이는 반대편 연결이 종료된 것을 의미한다.
					n, err := c.Read(buf)
					if err != nil {
						if !errors.Is(err, io.EOF) {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}
			}(conn)
		}
	}()

	// Listener 의 연결 객체를 가져온다.
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	// 연결을 종료한다.
	if err := conn.Close(); err != nil {
		t.Logf("unexpected close err : %v", err)
	}
	<-done
	if err := listener.Close(); err != nil {
		t.Logf("unexpected close err : %v", err)
	}
	<-done
}
