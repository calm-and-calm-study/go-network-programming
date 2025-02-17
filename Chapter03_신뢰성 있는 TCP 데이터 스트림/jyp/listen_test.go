package ch03

import (
	"net"
	"testing"
)

// 단순히 listener 생성 후 종료라 패킷 캡처 안됨
func TestListener(t *testing.T) {
	// listener 만들고 0 or 안적으면 random, 적으면 해당 포트로 listener 생성
	listener, err := net.Listen("tcp", "127.0.0.1:7138")
	if err != nil {
		t.Fatal(err)
	}
	// listener 생성후 close
	defer func() { _ = listener.Close() }()

	t.Logf("bound to %q", listener.Addr())
}
