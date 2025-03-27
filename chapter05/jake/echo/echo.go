package echo

import (
	"context"
	"fmt"
	"net"
)

// 전송한 데이터가 되돌아오는 서버
func echoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	go func() {
		go func() {
			<-ctx.Done()
			_ = s.Close()
		}()

		// conn, err := listener.Accept()
		// TCP 와 달리 UDP 에서는 Accept 메서드를 사용하지 않음
		// UDP 에서는 핸드쉐이크 과정이 없기 때문임

		// IP 패킷의 최대 크기 : MTU (Maximum Transmission Unit)
		// 이더넷의 기본 MTU : 1500 바이트
		// UDP/IP 헤더 크기 = 8바이트(UDP) + 20바이트(IPv4) = 28 바이트
		// UDP 데이터 최대 크기 = 1500 - 20 - 8 = 1472 바이트
		// 안전한 기본값 설정인 1024
		// MTU 보다 작은 값을 사용하는 이유는 패킷이 분할될 가능성이 낮기 때문
		buf := make([]byte, 1024)
		for {
			// clientAddr : 송신자 주소
			n, clientAddr, err := s.ReadFrom(buf) // client to server
			if err != nil {
				return
			}

			// 송신자에게 데이터를 다시 전송함
			_, err = s.WriteTo(buf[:n], clientAddr) // server to client
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
