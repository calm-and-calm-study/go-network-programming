package ch03

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

// 코드 이해가 안됨, 어떻게 하나만 되지?
func TestDialContextCancelFanOut(t *testing.T) {
	// context.WithDeadline
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)

	listener, err := net.Listen("tcp", "127.0.0.1:7138")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	// 다이얼마다 id 부여
	dial := func(ctx context.Context, address string, response chan int,
		id int, wg *sync.WaitGroup) {
		defer wg.Done()

		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			return
		}
		c.Close()

		select {
		case <-ctx.Done():
		case response <- id:
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	// 각 고루틴마다 Add(1) 호출
	for i := 0; i < 10; i++ {
		wg.Add(1)
		// 다이얼 생성하면서 id 부여하고 다이얼 끝나면 Done 처리리
		t.Log(listener.Addr().String())
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res
	cancel()
	// 고루틴 전부끝나는거 기다림
	wg.Wait()
	//끝나면 close
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s",
			ctx.Err(),
		)
	}

	t.Logf("dialer %d retrieved the resource", response)
}
