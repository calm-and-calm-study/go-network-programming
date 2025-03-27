package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-ping/ping"
)

// sudo go run icmp_ping.go
// ICMP Ping 체크
func icmpPing(address string) bool {
	pinger, err := ping.NewPinger(address)
	if err != nil {
		log.Printf("[ICMP Ping 실패]: %v\n", err)
		return false
	}

	// MacOS에서는 반드시 필요!
	pinger.SetPrivileged(true)

	// 타임아웃 설정
	pinger.Count = 3
	pinger.Timeout = 2 * time.Second

	log.Printf("[ICMP Ping 요청]: %s\n", address)
	err = pinger.Run()
	if err != nil {
		log.Printf("[ICMP Ping 실패]: %v\n", err)
		return false
	}

	stats := pinger.Statistics()
	log.Printf("[ICMP Ping 성공]: %s, 평균 RTT: %v\n", address, stats.AvgRtt)
	return true
}

func main() {
	target := "8.8.8.8" // Google DNS
	if icmpPing(target) {
		fmt.Println("✅ 서버 연결 성공!")
	} else {
		fmt.Println("❌ 서버 연결 실패!")
	}
}
