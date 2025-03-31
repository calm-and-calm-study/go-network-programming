package ch03

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

// pi
func TestPingerAdvanceDeadline(t *testing.T) {
	done := make(chan struct{})

	//리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	// 현재 시간 받음
	begin := time.Now()

	// 고루틴
	go func() {
		//고루틴 끝나면 done 채널 종료
		defer func() { close(done) }()

		//리스너에서 Acccep()
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		// 취소 가능한 컨텍스트 생성
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			conn.Close()
		}()

		// Pinger 간격을 조절하는 채널
		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second        // Ping 간격 1초
		go Pinger(ctx, conn, resetTimer) // ping 실행 1초간격

		// 데드라인 현재에서 5초후로 설정
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			t.Logf("[%s] %s",
				time.Since(begin).Truncate(time.Second), buf[:n])

			resetTimer <- 0
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}
	_, err = conn.Write([]byte("PONG!!!"))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}
	<-done
	end := time.Since(begin).Truncate(time.Second)
	t.Logf("[%s] done", end)
	if end != 9*time.Second {
		t.Fatalf("expected EOF at 9 seconds; actual %s", end)
	}
}
