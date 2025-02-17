package main

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

const payload = "The bigger the interface, the weaker the abstraction."

// 스캐너 테스트
func TestScanner(t *testing.T) {
	// 리스너 생성
	listener, err := net.Listen("tcp", "127.0.0.1:7138")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		//연결 수립
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		//클라이언트에게 payload 보냄
		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	// 클라이언트가 연결 요청
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	//bufio.Scanner로 conn에 대해 스캐너 설정
	scanner := bufio.NewScanner(conn)
	//스캐너로 스캔 word 나눔
	scanner.Split(bufio.ScanWords)

	var words []string

	//scan 만큼 for문 실행
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}

	// 스캐너로 들어오는 기대값
	expected := []string{"The", "bigger", "the", "interface,", "the",
		"weaker", "the", "abstraction."}

	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
