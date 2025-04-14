## HTTP 서비스 작성

#### Go HTTP 서버 
- 멀티플렉서는 클라이언트 요청을 수신 (Router 역할)
- 멀티플렉서는 요청의 목적지를 결정한 후 해당 요청을 처리 할 수 있는 능력이 있는 객체로 클라이언트 요청 전달
- 이 객체를 handler
- 이 중간에 middleware - 동작을 변경하거나 로깅, 인증 및 접근 제어

#### 응답 쓰는 순서
- 상태 코드 (Header) 부터 쓰고 body 를 씀
- Write 메서드 호출하면 Go는 암묵적으로 WriteHeader의 http.StatusOK를 호출
- pitfall_test.go 관련 예제 

#### 미들웨어
- 미들웨어는 http.handler를 매개변수로 받아서 http.handler를 한환하는 재사용할 수 있는 함수로 구성
- 민감한 파일 접근 제한 방법 - prefix가 (.)으로 시작하는 파일은 접근하지 못하도록함

#### 멀티플렉서
- 기본적으로 http.ServeMux는 모든 요청에 대해 404 Not Found 로 응답

#### HTTP/2 서버 푸시
- HTTP/1.1 과 HTTP/2 를 네트워크 탭에서 보면 동작방식이 다름
- HTTP/2는 push를 통해 추가적인 별도 요청에 대한 오버헤드가 줄어듬

#### 기타
- httptest 이용한 핸들러 테스트, httptest.NewRequest 함수는 에러를 반환하는 대신 패닉함
- 클라이언트가 요청-응답 생명주기의 기간을 정하게 해선 안됨 - 보안 (DDos 공격) 위협에 노출