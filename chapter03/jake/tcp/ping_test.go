package tcp

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

const defaultPingInterval = 30 * time.Second

func TestPingerAdvanceDeadline(t *testing.T) {
	done := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:8090")
	if err != nil {
		t.Fatal(err)
	}

	begin := time.Now()
	go func() {
		defer func() {
			close(done)
		}()

		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			conn.Close()
		}()

		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second
		go Pinger(ctx, conn, resetTimer)

		if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])

			resetTimer <- 0
			if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ {
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	_, err = conn.Write([]byte("PONG!!"))
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
		t.Fatalf("expected EOF at 9 seconds; actual : %s", end)
	}
}

func Pinger(ctx context.Context, w io.Writer, resetTimer <-chan time.Duration) {
	var interval time.Duration
	select {
	case <-ctx.Done():
		return
	case interval = <-resetTimer:
	default:
	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case newInterval := <-resetTimer:
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}

		_ = timer.Reset(interval)
	}
}
