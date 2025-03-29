package tftp

import (
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

// Server represents a read-only TFTP server that supports a subset of
// RFC 1350.
type Server struct {
	// 모든 읽기 요청에 반환될 페이로드
	Payload []byte // the payload served for all read requests

	// 전송 실패 시 재시도 횟수
	Retries uint8 // the number of times to retry a failed transmission

	// 전송 승인을 기다릴 시간
	Timeout time.Duration // the duration to wait for an acknowledgment
}

func (s Server) ListenAndServe(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	log.Printf("Listening on %s ...\n", conn.LocalAddr())

	return s.Serve(conn)
}

// Serve 서버 데이터 수신
func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}

	if s.Payload == nil {
		return errors.New("payload is required")
	}

	if s.Retries == 0 {
		s.Retries = 10
	}

	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}

	var rrq ReadReq

	for {
		buf := make([]byte, DatagramSize)

		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("[%s] bad request: %v", addr, err)
			continue
		}

		go s.handle(addr.String(), rrq)
	}
}

func (s Server) handle(clientAddr string, rrq ReadReq) {
	log.Printf("[%s] requested file: %s", clientAddr, rrq.Filename)

	// 클라이언트와의 연결
	conn, err := net.Dial("udp", clientAddr)
	if err != nil {
		log.Printf("[%s] dial: %v", clientAddr, err)
		return
	}
	defer func() { _ = conn.Close() }()

	var (
		ackPkt Ack
		errPkt Err
		// 전송하고자 하는 데이터 객체
		dataPkt = Data{Payload: bytes.NewReader(s.Payload)}
		buf     = make([]byte, DatagramSize)
	)

	// 패킷을 전부 수신하기 위한 for 문
NEXTPACKET:
	// 512 Bytes
	for n := DatagramSize; n == DatagramSize; {
		// 데이터 전송을 위한
		data, err := dataPkt.MarshalBinary()
		if err != nil {
			log.Printf("[%s] preparing data packet: %v", clientAddr, err)
			return
		}

	RETRY:
		for i := s.Retries; i > 0; i-- {
			// 데이터 전송
			// 즉, 데이터 블록 패킷 1번 전송 후 확인 패킷 대기하고 다시 전송하는 루틴
			n, err = conn.Write(data) // send the data packet
			if err != nil {
				log.Printf("[%s] write: %v", clientAddr, err)
				return
			}

			// wait for the client's ACK packet
			// 클라이언트의 ack 패킷 대기
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))

			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					// 에러 발생 시 RETRY 라벨의 루틴만 continue
					continue RETRY
				}

				log.Printf("[%s] waiting for ACK: %v", clientAddr, err)
				return
			}

			switch {
			// 데이터 전송에 대한 응답 패킷 수신
			case ackPkt.UnmarshalBinary(buf) == nil:
				if uint16(ackPkt) == dataPkt.Block {
					// received ACK; send next data packet
					continue NEXTPACKET
				}
				// 에러 패킷 수신
			case errPkt.UnmarshalBinary(buf) == nil:
				log.Printf("[%s] received error: %v",
					clientAddr, errPkt.Message)
				return
			default:
				log.Printf("[%s] bad packet", clientAddr)
			}
		}

		log.Printf("[%s] exhausted retries", clientAddr)
		return
	}

	log.Printf("[%s] sent %d blocks", clientAddr, dataPkt.Block)
}
