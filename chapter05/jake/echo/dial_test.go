package echo

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

func TestDialUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.Dial("udp", serverAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	// Send a message to the client from a rogue connection.
	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// UDP 는 비연결형 프로토콜 (연결을 유지하지 않음) 이므로 패킷이 도착할 때마다 NIC (네트워크 인터페이스 카드)에서 발생하는 인터럽트를 처리해야 한다.
	// 인터럽트의 연결로부터 클라이언트에게 메시지를 보냄
	interrupt := []byte("pardon me")
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_ = interloper.Close()

	if len(interrupt) != n {
		t.Fatalf("wrote %d bytes of %d", n, len(interrupt))
	}

	// Now write a message to the server that will prompt a reply.
	ping := []byte("ping")
	_, err = client.Write(ping)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err = client.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	// The first message the client reads should be the "ping" from the echo
	// server, not the queued up "test" message.
	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", ping, buf[:n])
	}

	// Verify no other incoming packets are waiting.
	err = client.SetDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Read(buf)
	if err == nil {
		t.Fatal("unexpected packet")
	}
}
