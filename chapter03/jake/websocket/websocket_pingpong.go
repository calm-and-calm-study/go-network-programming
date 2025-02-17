package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// WebSocket 서버 핸들러 (Ping/Pong 테스트)
func pingPongWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Ping 핸들러 설정 (자동으로 Pong 응답)
	conn.SetPingHandler(func(appData string) error {
		fmt.Println("Received Ping, sending Pong...")
		return conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	for {
		// 메시지 읽기 (Ping/Pong 확인)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Println("Received message:", string(msg))
	}
}

// WebSocket 서버 시작
func startPingPongServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", pingPongWsHandler)

	server := &http.Server{Addr: "localhost:8082", Handler: mux}

	go func() {
		_ = server.ListenAndServe()
	}()

	// 서버가 시작될 시간을 줌
	time.Sleep(500 * time.Millisecond)
	return server
}
