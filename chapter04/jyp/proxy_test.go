package main

import (
	"io"
	"net"
	"sync"
	"testing"
)

func proxy(from io.Reader, to io.Writer) error {
	//from 이 io.Writer 인지 확인
	fromWriter, fromIsWriter := from.(io.Writer)
	//to 가 io.Reader 인지 확인
	toReader, toIsReader := to.(io.Reader)

	// (to -> from)
	if toIsReader && fromIsWriter {
		// to 데이터를 from 으로 복사하여 응답
		go func() { _, _ = io.Copy(fromWriter, toReader) }()
	}

	// (from -> to)
	_, err := io.Copy(to, from)
	return err
}

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	// 서버 리스너 생성
	server, err := net.Listen("tcp", "127.0.0.1:7139")
	if err != nil {
		t.Fatal(err)
	}

	//
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				return
			}

			go func(c net.Conn) {
				defer c.Close()

				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}

						return
					}

					// 버퍼에서 읽어들일 때 ping 이면 pong 송신, 아니면 온 그대로 송신
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

	//프록시 서버 생성
	proxyServer, err := net.Listen("tcp", "127.0.0.1:7138")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			conn, err := proxyServer.Accept()
			if err != nil {
				return
			}

			go func(from net.Conn) {
				defer from.Close()
				to, err := net.Dial("tcp",
					server.Addr().String())
				if err != nil {
					t.Error(err)
					return
				}
				defer to.Close()

				err = proxy(from, to)
				if err != nil && err != io.EOF {
					t.Error(err)
				}
			}(conn)
		}
	}()

	//클라이언트 다이얼
	conn, err := net.Dial("tcp", proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Reply string }{
		{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
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
			t.Errorf("%d: expected reply: %q; actual: %q",
				i, m.Reply, actual)
		}
	}

	_ = conn.Close()
	_ = proxyServer.Close()
	_ = server.Close()
	wg.Wait()
}
