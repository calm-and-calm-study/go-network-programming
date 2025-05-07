## Caddy
- Go 동시성과 병렬성의 강점을 취하여 상당량의 웹 트래픽 서빙 가능
- Let's Encrypt 통합으로 인해 ACME(자동화된 인증서 관리 환경) 프로토콜 지원

### 설치 경로
```
git clone "https://github.com/caddyserver/caddy.git"
```

### Caddy 기능
- 관리자 포트 2019
- GET 으로 config 확인 가능
- POST 요청으로 JSON 데이터에 config 지정함

### Caddy Middleware
