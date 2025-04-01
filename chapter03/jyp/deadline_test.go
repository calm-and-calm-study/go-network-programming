package ch03

import (
	"io"
	"net"
	"testing"
	"time"
)

// 패킷 캡처 안잡힘
func TestDeadline(t *testing.T) {
	sync := make(chan struct{})

	listener, err := net.Listen("tcp", "127.0.0.1:7138")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer func() {
			conn.Close()
			close(sync) // read from sync shouldn't block due to early return
		}()

		// connection 연결된 후 데드라인 설정 5초
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		//버퍼 생성 후 conn.Read 로 데이터 읽어들임임
		buf := make([]byte, 1)
		_, err = conn.Read(buf) // blocked until remote node sends data
		nErr, ok := err.(net.Error)
		if !ok || !nErr.Timeout() {
			t.Errorf("expected timeout error; actual: %v", err)
		}

		sync <- struct{}{}

		// conn.Read 상태일 떄떄 데드라인을 뒤로 미룸룸
		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)
		if err != nil {
			t.Error(err)
		}
	}()

	// Dial
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	//?
	<-sync
	//connection 에서 Write 1 함
	_, err = conn.Write([]byte("1"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf) // 원격 노드가 데이터를 보낼 때까지 블로킹됨
	if err != io.EOF {
		t.Errorf("expected server termination; actual: %v", err)
	}
}
