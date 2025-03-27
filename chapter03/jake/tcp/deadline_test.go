package tcp

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})
	// 랜덤 port 로 listen
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// 고루틴 생성
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer func() {
			conn.Close()
			close(sync)
		}()

		// 5초 타임아웃 지정
		if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1)
		_, err = conn.Read(buf)
		// error 타입 변환
		nErr, ok := err.(net.Error)
		// timeout 에러가 아닌 경우 에러로그 출력
		if !ok || !nErr.Timeout() {
			t.Errorf("expected timeout error; actual: %v", err)
		}
		// 완료 신호 채널로 전송
		// 첫번째 연결이 종료된 후 완료 신호를 보냄
		sync <- struct{}{}

		if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Error(err)
			return
		}

		// 두번째 conn 은 아래의 client conn 으로부터 데이터를 받음
		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
	}()

	// client 역할을 하는 conn 연결 구현체
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-sync
	// 데이터를 전송해서 두번째 conn.Read 에서는 데이터를 받게 됨
	_, err = conn.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != io.EOF {
		t.Errorf("expected server termination; actual: %v", err)
	}
}
