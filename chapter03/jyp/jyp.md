# 연결 수립
## TCP 쓰리핸드셰이크 
1. SYN
2. SYN/ACK
3. ACK

- 위의 세번의 과정을 통해 TCP 세션을 수립

## 3 hand shake 시 시퀀스 전달 방법 
- 이유 
    - SYN에 대한 ACK 가 어떤것인지 확인하기 위해
    - 몇번까지 보냈느지 확인하기 위해 
    - 중간에 끝긴 트래픽 SACK 패킷으로 확인 가능, 수신자가 어떤 패킷을 수신하였는지 송신자에게 알려줌
- SYN : sequence X
- SYN / ACK : SYN 은 sequence Y / ACK 는 sequence X+1
- ACK : sequnce Y+1

# 데이터 보내기
## 수신 버퍼와 슬라이스 윈도 크기
- 수신버퍼는 얼마만큼 받을수 있는지 예비해둔 공간 
- ACK 패킷에 특별하게 중요한 window size(수신 확인할 필요없이?에 의미가 애매)


# 연결 종료
- 일반 종료 과정
    - FIN (클라이언트)
        - 클라이언트 FIN_WAIT_1 상태
        - 서버는 CLOSE_WAIT 상태
    - ACK (서버)
        - FIN_WAIT_2
    - FIN (서버)
        - 서버는 LASK_ACK
        - 클라이언트는 TIME_WAIT
    - ACK (클라이언트)
        - 클라이언트는 ACK 보내고 CLOSED 상태
        - 서버도 받으면 CLOSED

- 갑자기 중간에 끝기면?
    - RST(reset-초기화) 패킷 보내서 수신 측에서 더이상 데이터 수신할 수 없다고 보냄냄



# TCP Timeout

- 일반적으로 Timeout 에 대한 설정이 없으면 운영체제의 Timeout을 따름
- 운영체제의 Timeout 이 2시간이 될 수 있어 서비스에 불필요한 Timeout 이 될 수 있음
- 코드에서 명시적으로 Timeout 설정


# Timeout 설정 방법
- DiatTimeout
- ContextDeadline
- ContextCancel - time.Sleep 으로 지연

# 패킷 캡처
https://zepplin86.tistory.com/11