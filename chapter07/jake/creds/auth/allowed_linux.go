package auth

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os/user"

	"golang.org/x/sys/unix"
)

// Allowed
// net.UnixConn 파일 객체를 변수로 저장합니다.
// 이 객체는 호스트상의 유닉스 도메인 소켓 연결 객체를 의미합니다.
func Allowed(conn *net.UnixConn, groups map[string]struct{}) bool {
	if conn == nil || groups == nil || len(groups) == 0 {
		return false
	}

	file, _ := conn.File()
	defer func() { _ = file.Close() }()

	var (
		err   error
		ucred *unix.Ucred
	)

	for {
		// 파일 객체의 디스크립터, 어느 프로토콜 계층에 속하였는지를 나타내는 상수인 unix.SOL_SOCKET, 옵션 값인 unix.SO_PEERCRED 를 전송
		// 반환 받은 ucred 에는 피어의 프로세스 정보, 사용자ID, 그룹ID 정보가 있음
		// 사용자는 하나 이상의 그룹에 속할 수 있으므로 각각의 그룹에 대한 권한을 확인해야 한다.
		// 허용된 그룹들에 대해 각각의 그룹 ID 를 비교해야 한다.
		ucred, err = unix.GetsockoptUcred(int(file.Fd()), unix.SOL_SOCKET, unix.SO_PEERCRED)
		if errors.Is(err, unix.EINTR) {
			continue // syscall interrupted, try again
		}
		if err != nil {
			log.Println(err)
			return false
		}

		break
	}

	u, err := user.LookupId(fmt.Sprint(ucred.Uid))
	if err != nil {
		log.Println(err)
		return false
	}

	gids, err := u.GroupIds()
	if err != nil {
		log.Println(err)
		return false
	}

	for _, gid := range gids {
		if _, ok := groups[gid]; ok {
			return true
		}
	}

	return false
}
