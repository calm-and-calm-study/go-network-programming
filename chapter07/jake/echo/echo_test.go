package echo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnix(t *testing.T) {
	// 하위 디렉터리 생성
	dir, err := ioutil.TempDir("", "echo_unix")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	// 소켓파일의 이름을 생성합니다.
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	rAddr, err := streamingEchoServer(ctx, "unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = os.Chmod(socket, os.ModeSocket|0666)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("unix", rAddr.String())
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("ping")
	for i := 0; i < 3; i++ { // write 3 "ping" messages
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf) // read once from the server
	if err != nil {
		t.Fatal(err)
	}

	// ping 메시지를 3번 받게 되면 세번의 문자열이 하나의 pingpingping 으로 수신하게 됩니다.
	// 스트림 기반의 연결에는 메시지의 구분자가 없기 때문입니다.
	// 서버의 바이트 스트림으로부터 하나의 메시지가 끝나는 지점을 읽고 구분하는 것은 코드 상에서 구현해야 합니다.
	expected := bytes.Repeat(msg, 3)
	if !bytes.Equal(expected, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", expected,
			buf[:n])
	}
}

func BenchmarkEchoServerUDP(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := datagramEchoServer(ctx, "udp", "127.0.0.1:")
	if err != nil {
		b.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	msg := []byte("ping")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			b.Fatal(err)
		}

		buf := make([]byte, 1024)
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			b.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			b.Fatalf("received reply from %q instead of %q", addr,
				serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			b.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}

func BenchmarkEchoServerTCP(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	rAddr, err := streamingEchoServer(ctx, "tcp", "127.0.0.1:")
	if err != nil {
		b.Fatal(err)
	}
	defer cancel()

	conn, err := net.Dial("tcp", rAddr.String())
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	msg := []byte("ping")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = conn.Write(msg)
		if err != nil {
			b.Fatal(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			b.Fatal(err)
		}

		if !bytes.Equal(msg, buf[:n]) {
			b.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}

func BenchmarkEchoServerUnix(b *testing.B) {
	dir, err := ioutil.TempDir("", "echo_unix_bench")
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			b.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	rAddr, err := streamingEchoServer(ctx, "unix", socket)
	if err != nil {
		b.Fatal(err)
	}
	defer cancel()

	conn, err := net.Dial("unix", rAddr.String())
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	msg := []byte("ping")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = conn.Write(msg)
		if err != nil {
			b.Fatal(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			b.Fatal(err)
		}

		if !bytes.Equal(msg, buf[:n]) {
			b.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}
	}
}

func TestEchoServerUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := datagramEchoServer(ctx, "udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	msg := []byte("ping")
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != serverAddr.String() {
		t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
	}

	if !bytes.Equal(msg, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
	}
}

func TestEchoServerTCP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rAddr, err := streamingEchoServer(ctx, "tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	conn, err := net.Dial("tcp", rAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	msg := []byte("ping")
	_, err = conn.Write(msg)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(msg, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
	}
}
