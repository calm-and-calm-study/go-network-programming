package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

// 자체적으로 타임아웃나게 생성한거라 패킷캡처에 찍히지 않음
func TestDialContext(t *testing.T) {
	//데드라인을 현재 시간 +5초로 생성
	dl := time.Now().Add(5 * time.Second)
	// context.WithDeadline 생성성
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()

	var d net.Dialer
	// 다이얼 컨트롤을 데드라인 살짝 넘기도록 생성
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}
	// 다이얼 컨텍스트로 데드라인 ctx 할당당
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}

	// 네트워크 에러 감지
	nErr, ok := err.(net.Error)

	// 아니면 에러 발생
	if !ok {
		t.Error(err)
	} else {
		// 맞으면 타임아웃인지 체크
		if !nErr.Timeout() {
			// 아니면 에러 발생
			t.Errorf("error is not a timeout: %v", err)
		}
	}

	// ctx 에러가 데드라인 초과가 아니면 에러 발생
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}
