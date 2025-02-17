package tcp

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(10*time.Second),
	)

	listener, err := net.Listen("tcp", "127.0.0.1:8092")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		// 연결을 수락한 후 에러가 없으면 바로 연결을 종료시킴
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	dial := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup) {
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

	// 10개의 다른 연결을 생성하지만 이미 최초 연결 이후 아이디값이 존재하면 채널을 통해서 이벤트를 전달함
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	// 최초 연결의 아이디를 수신하면 cancel 을 호출하여 종료 신호를 호출함
	response := <-res
	// context 는 부모 context 에서 dial 함수로 전달했기 때문에 부모 context 의 cancel 신호는 10개의 dial 함수에 종료신호로 전파된 후 종료
	cancel()
	// 모든 고루틴이 정상적으로 종료되는 것을 기다림
	wg.Wait()
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual : %s", ctx.Err())
	}

	t.Logf("dialer %d retrieved the resource", response)
}
