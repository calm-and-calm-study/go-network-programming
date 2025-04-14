package middleware

import (
	"net/http"
	"path"
	"strings"
)

// RestrictPrefix prefix 를 검사하는 미들웨어 로직
// 접두어를 사용하여 허용한 파일만 제공하도록 검사
func RestrictPrefix(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for _, p := range strings.Split(path.Clean(r.URL.Path), "/") {
				if strings.HasPrefix(p, prefix) {
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}
