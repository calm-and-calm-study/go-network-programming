package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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

	server := &http.Server{Addr: "localhost:8091", Handler: mux}

	go func() {
		_ = server.ListenAndServe()
	}()

	// 서버가 시작될 시간을 줌
	time.Sleep(500 * time.Millisecond)
	return server
}
