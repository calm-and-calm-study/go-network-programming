package main

import (
	"io"
	"net"
)

// 데이터 프록시
func proxyConn(source, destination string) error {
	// source IP로 다이얼
	connSource, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}
	defer connSource.Close()

	// destination으로 다이얼
	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()

	// 서버에서 클라이언트로 데이터 포워딩
	go func() { _, _ = io.Copy(connSource, connDestination) }()

	// 클라이언트에서 서버로 데이터 포워딩
	_, err = io.Copy(connDestination, connSource)

	return err
}

var _ = proxyConn
