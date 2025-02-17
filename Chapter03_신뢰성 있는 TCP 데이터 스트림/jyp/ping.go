package ch03

import (
	"context"
	"io"
	"time"
)

// ping 인터벌 30초,
const defaultPingInterval = 30 * time.Second

// 설명 필요..
func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration

	select {
	case <-ctx.Done():
		return
	case interval = <-reset: // pulled initial interval off reset channel
	default:
	}
	// 인터벌이 0 이하면 초기화화
	if interval <= 0 {
		interval = defaultPingInterval
	}

	// 인터벌로 타이머 생성
	timer := time.NewTimer(interval)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	//
	for {
		select {
		case <-ctx.Done():
			return
		case newInterval := <-reset:
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			if _, err := w.Write([]byte("ping")); err != nil {
				// track and act on consecutive timeouts here
				return
			}
		}

		_ = timer.Reset(interval)
	}
}
