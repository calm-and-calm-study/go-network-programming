package echo

import (
	"context"
	"fmt"
	"net"
	"os"
)

func streamingEchoServer(ctx context.Context, network string, addr string) (net.Addr, error) {
	s, err := net.Listen(network, addr)
	if err != nil {
		return nil, fmt.Errorf("binding to %s %s: %w", network, addr, err)
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		for {
			conn, err := s.Accept()
			if err != nil {
				return
			}

			go func() {
				defer func() { _ = conn.Close() }()

				for {
					buf := make([]byte, 1024)
					n, err := conn.Read(buf)
					if err != nil {
						return
					}

					_, err = conn.Write(buf[:n])
					if err != nil {
						return
					}
				}
			}()
		}
	}()

	return s.Addr(), nil
}

func datagramEchoServer(ctx context.Context, network string,
	addr string) (net.Addr, error) {
	// 만일 netListen 함수나 net.ListenUnit 함수를 사용하지 않으면 golang 에서는 소켓 파일을 지우지 않음
	// 따라서 코드 상에서 반드시 소켓 파일을 지워야 함
	s, err := net.ListenPacket(network, addr)
	if err != nil {
		return nil, fmt.Errorf("binding to %s %s: %w", network, addr, err)
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
			if network == "unixgram" {
				_ = os.Remove(addr)
			}
		}()

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
