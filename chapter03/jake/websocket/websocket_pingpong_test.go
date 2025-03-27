package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// Ping/Pong 테스트
func TestWebSocketPingPong(t *testing.T) {
	// WebSocket 서버 시작
	server := startPingPongServer()
	defer server.Close()

	// 클라이언트 WebSocket 연결
	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second, // 연결 타임아웃 설정
	}

	conn, _, err := dialer.Dial("ws://localhost:8082/ws", nil)
	assert.NoError(t, err, "WebSocket 연결 실패")
	defer conn.Close()

	// Ping 메시지 전송 후 Pong 응답 확인
	pingData := "ping_payload"
	err = conn.WriteMessage(websocket.PingMessage, []byte(pingData))
	assert.NoError(t, err, "Ping 메시지 전송 실패")

	// Pong 응답 받기
	msgType, msg, err := conn.ReadMessage()
	assert.NoError(t, err, "Pong 응답 수신 실패")
	assert.Equal(t, websocket.PongMessage, msgType, "응답이 Pong이 아님")
	assert.Equal(t, pingData, string(msg), "Pong 응답 데이터 불일치")

	fmt.Println("Test Passed: WebSocket Ping/Pong 정상 동작")
}
