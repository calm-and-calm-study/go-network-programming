package main

import (
	"flag"
	"fmt"
	"go-network-programming/chapter07/jake/creds/auth"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
)

func init() {
	// 명령줄 인자 파싱하는 로직
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(),
			"Usage:\n\t%s <group names>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func parseGroupNames(args []string) map[string]struct{} {
	groups := make(map[string]struct{})

	for _, arg := range args {
		// 그룹정보 조회
		grp, err := user.LookupGroup(arg)
		if err != nil {
			log.Println(err)
			continue
		}

		// 그룹 맵에 추가
		groups[grp.Gid] = struct{}{}
	}

	return groups
}

func main() {
	flag.Parse()

	groups := parseGroupNames(flag.Args())
	socket := filepath.Join(os.TempDir(), "creds.sock")
	addr, err := net.ResolveUnixAddr("unix", socket)
	if err != nil {
		log.Fatal(err)
	}

	s, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
	}

	// 애플리케이션 종료 신호
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		_ = s.Close()
	}()

	fmt.Printf("Listening on %s ...\n", socket)

	for {
		// AcceptUnix 를 이용하여 연결 수립 요청
		conn, err := s.AcceptUnix()
		if err != nil {
			break
		}
		if auth.Allowed(conn, groups) {
			_, err = conn.Write([]byte("Welcome\n"))
			if err == nil {
				// handle the connection in a goroutine here
				continue
			}
		}

		_, err = conn.Write([]byte("Access denied\n"))
		if err != nil {
			log.Println(err)
		}

		_ = conn.Close()
	}
}
