package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// WebSocket 서버 핸들러 (Write Timeout 발생)
func writeTimeoutWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 클라이언트가 일정 시간 안에 메시지를 보내지 않으면 타임아웃 발생
	writeTimeout := 5 * time.Second

	for {
		// 쓰기 타임아웃 설정
		conn.SetWriteDeadline(time.Now().Add(writeTimeout))

		// 더미 메시지 전송 시도 (타임아웃 발생 여부 확인)
		err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			fmt.Println("Write timeout occurred:", err)
			return
		}

		// 메시지를 일정 간격(6초)으로 보냄 (타임아웃 초과)
		time.Sleep(6 * time.Second)
	}
}

func startWriteTimeoutServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", writeTimeoutWsHandler)

	server := &http.Server{Addr: "localhost:8081", Handler: mux}

	go func() {
		_ = server.ListenAndServe()
	}()

	// 서버가 시작될 시간을 줌
	time.Sleep(500 * time.Millisecond)
	return server
}

func TestWebSocketWriteTimeout(t *testing.T) {
	// WebSocket 서버 시작
	server := startWriteTimeoutServer()
	defer server.Close()

	// 클라이언트 WebSocket 연결
	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second, // 연결 타임아웃 설정
	}

	conn, _, err := dialer.Dial("ws://localhost:8081/ws", nil)
	assert.NoError(t, err, "WebSocket 연결 실패")
	defer conn.Close()

	// 일정 시간 동안 메시지를 보내지 않고 대기 (타임아웃 유발)
	time.Sleep(7 * time.Second)

	// WebSocket 쓰기 (WriteMessage) 테스트 → 타임아웃 발생 예상
	err = conn.WriteMessage(websocket.TextMessage, []byte("test message"))
	assert.Error(t, err, "WebSocket write should fail due to timeout")
	fmt.Println("Test Passed: WebSocket write timeout occurred")
}
