## HTTP 클라이언트 작성

- HTTP 는 클라이언트 서버 기반의 세션을 갖지 않는 프로토콜
- HTTP는 TCP 위에서 동작하는 프로토콜
  
#### 통합 리소스 식별자(Uniform Resource Locator)
- 클라이언트가 웹 서버를 찾고 요청된 리소스를 식별하는데 사용
- 구분자
  - 스키마 : https, wss 등등
  - 권한정보 : github 접속 시 jyp-catenoid:token@host:port
  - 경로 : router path
  - 쿼리 파라미터 : query string
  - 정보 조각 : #으로 html id 접근자

#### 클라이언트 리소스 요청
- GET : 리소스 요청
- HEAD : GET 요청과 유사하지만 리소스 응답하지 않음
- POST : 새로운 리소스 생성
- PUT : 완전 교체
- PATCH : 부분 교체
- DELETE : 리소스 제거
- OPTIONS : 어떠한 메소드 지원하는지 확인용
- CONNECT : HTTP 터널링
- TRACE : 에코잉
  
#### 서버 응답
- 200 번대 : 성공적 응답
- 300 번대 : 클라이언트가 보낸 요청에 대해 추가로 무언가 해야되는 경우
- 400 번대 : 클라이언트 요청이 잘못되어 에러가 발생한 경우
- 500 번대 : 서버 측에서 에러가 발생한 경우
