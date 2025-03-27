package tcp

import (
	"context"
	"errors"
	"io"
	"net"
	"syscall"
	"testing"
	"time"
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

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	c, err := DialTimeout("tcp", "10.0.0.1:http", 5*time.Second)
	if err == nil {
		c.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Fatal(err)
	}
	if !nErr.Timeout() {
		t.Fatal("connection is not a timed out")
	}
}

func TestDialContext(t *testing.T) {
	// deadline 지정
	dl := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()

	var d net.Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// 5초 이상의 timeout 을 강제로 발생시킴
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		// timeout 이 발생해야하기 때문에 err 가 없으면 테스트 실패
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		// err 가 net.Error 유형이 아닌 경우 실패처리
		t.Error(err)
	} else {
		// err 가 타임아웃 발생이 안된거면 테스트 실패
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout:%v", err)
		}
	}
	// context 의 에러가 DeadlineExceeded 가 아닌 경우 실패처리
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
	t.Logf("ctx error occur : %s", ctx.Err().Error())
}
