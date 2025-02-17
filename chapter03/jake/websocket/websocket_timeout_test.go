package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketWriteTimeout(t *testing.T) {
	// WebSocket 서버 시작
	server := startWriteTimeoutServer()
	defer server.Close()

	// 클라이언트 WebSocket 연결
	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second, // 연결 타임아웃 설정
	}

	conn, _, err := dialer.Dial("ws://localhost:8091/ws", nil)
	assert.NoError(t, err, "WebSocket 연결 실패")
	defer conn.Close()

	// 일정 시간 동안 메시지를 보내지 않고 대기 (타임아웃 유발)
	time.Sleep(7 * time.Second)

	// WebSocket 쓰기 (WriteMessage) 테스트 → 타임아웃 발생 예상
	err = conn.WriteMessage(websocket.TextMessage, []byte("test message"))
	assert.Error(t, err, "WebSocket write should fail due to timeout")
	fmt.Println("Test Passed: WebSocket write timeout occurred")
}
