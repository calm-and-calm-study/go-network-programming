package ch03

import (
	"net"
	"syscall"
	"testing"
	"time"
)

/**
* Dialer로 Timeout 할 수 있는 function 생성
* 이유: 자체적으로 생성하지 않으면 운영체제 timeout 따라가게 되어 서비스 이용 시 불편하게 될 수 있음
* 자체적으로 timeout을 서비스에 맞게 조절하는 것이 필요요
 */
func DialTimeout(network, address string, timeout time.Duration,
) (net.Conn, error) {
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
		t.Fatal("error is not a timeout")
	}
}
