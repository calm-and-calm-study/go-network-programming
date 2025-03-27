package proxy

import (
	"io"
	"net"
	"sync"
	"testing"
)

func proxy(from io.Reader, to io.Writer) error {
	fromWriter, fromIsWriter := from.(io.Writer)
	toReader, toIsReader := to.(io.Reader)

	if toIsReader && fromIsWriter {
		go func() {
			_, _ = io.Copy(fromWriter, toReader)
		}()
	}

	_, err := io.Copy(to, from)
	return err
}

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	// 무작위 포트로 연결을 받음
	server, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				return
			}

			go func(c net.Conn) {
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
					}
					// 받은 msg 가 ping 이면 pong 을 응답하고, 아닌 경우 그대로 응답
					// :n 은 처음부터 n 의 길이까지 받은 byte 를 의미함
					switch msg := string(buf[:n]); msg {
					case "ping":
						_, err = c.Write([]byte("pong"))
					default:
						_, err = c.Write(buf[:n])
					}

					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
				}
			}(conn)
		}
	}()

	// proxy server 를 생성
	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			// 서버로부터 받은 listen connection
			conn, err := proxyServer.Accept()
			if err != nil {
				return
			}

			// proxyServer 로부터 들어오는 connection 당 고루틴 1개
			go func(from net.Conn) {
				defer from.Close()

				to, err := net.Dial("tcp", conn.LocalAddr().String())
				if err != nil {
					t.Error(err)
					return
				}
				defer to.Close()

				// to 로 전달함
				if err := proxy(from, to); err != nil {
					t.Error(err)
				}
			}(conn)
		}

	}()

	conn, err := net.Dial("tcp", proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Reply string }{
		{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "pong"},
		{"ping", "pong"},
	}

	for i, m := range msgs {
		_, err = conn.Write([]byte(m.Message))
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}

		actual := string(buf[:n])
		t.Logf("%q -> proxy -> %q", m.Message, actual)

		if actual != m.Reply {
			t.Errorf("%d: expected reploy: %q; actual: %q", i, m.Reply, actual)
		}
	}

	_ = conn.Close()
	_ = proxyServer.Close()
	_ = server.Close()

	wg.Wait()
}
