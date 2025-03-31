# net.Conn 인터페이스

- net.Conn 의 Read/Write 는 가각 io.Reader 와 io.Writer로 구현
- bufio.Scanner 로 공백으로 구분된 데이터를 버퍼로 읽어들임
- 동적 버퍼 사이즈 할당: TLV(Type-Length-Value)
- io.Copy를 통해 데이터 Proxy 가능
- io.MultiWriter 함수를 이용하여 단일 페이로드를 여러개의 네트워크 연결로 전송(다중 프록시?)
- io.TeeReader 함수를 사용하여 네트워크 연결로 부터 읽은 데이터를 로깅하는데 사용 가능
- ping.go 는 방화벽에서 ICMP 막혀있으면 대신해서 ping확인 가능 (go run ./ch04/ping.go 8.8.8.8:8080)
