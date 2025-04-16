package handlers

import (
	"html/template"
	"io"
	"net/http"
)

var t = template.Must(template.New("hello").Parse("Hello, {{.}}!"))

func DefaultHandler() http.Handler {
	// Go 의 http.Handler 인터페이스 구현한 부분
	return http.HandlerFunc(
		// 응답 부분을 정의
		func(w http.ResponseWriter, r *http.Request) {
			// 요청 body 에 대해서 close 처리
			defer func(r io.ReadCloser) {
				_, _ = io.Copy(io.Discard, r)
				_ = r.Close()
			}(r.Body)

			var b []byte

			switch r.Method {
			case http.MethodGet:
				b = []byte("friend")
			case http.MethodPost:
				var err error
				b, err = io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal server error",
						http.StatusInternalServerError)
					return
				}
			default:
				// not RFC-compliant due to lack of "Allow" header
				http.Error(w, "Method not allowed",
					http.StatusMethodNotAllowed)
				return
			}

			_ = t.Execute(w, string(b))
		},
	)
}
