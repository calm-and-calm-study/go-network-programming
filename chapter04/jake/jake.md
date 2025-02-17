
### net.Conn 인터페이스 

- 네트워크 코드의 대부분은 net.Conn 인터페이스에서 프로그래밍을 할 수 있음

- io.Reader, io.Writer 인터페이스를 구현

- SetReadDeadline, SetWriteDeadline 메서드는 매개변수로 입력받은 시간을 데드라인으로 설정하여 
  각각 읽기및 쓰기 동시에 대해 매개변수로 입력받은 시간을 데드라인으로 설정합니다.

- 네트워크 연결로부터 데이터를 읽고 쓰는 것은 파일 객체에 데이터를 읽고 쓰는 것과 같음
  이유는 net.Conn 인터페이스가 파일 객체의 io 를 구현한 io.ReadWriteCloser 인터페이스를 구현했기 때문

- bufio.Scanner 메서드를 이용하여 특정 구분자를 만날 때까지 데이터를 읽는 방법을 알아봄

- 다변하는 페이로드 크기로부터 동적으로 버퍼를 할당하는 기본 프로토콜을 정의할 수 있도록 해주는 인코딩 메서드인 TLV 

- 네트워크 연결로부터 데이터를 읽고 쓸 때 발생하는 에러를 처리하는 방법을 알아보는 시간

<br />

### reflect.DeepEqual

- 두 개의 값을 깊은 비교(Deep Comparision)

- 두 값을 재귀저긍로 비교, 데이터의 순서까지 비교, 값이 모두 동일해야 True 반환

```go
package main

import (
	"fmt"
	"reflect"
)

func main() {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 2, "a": 1} // 키 순서 다름

	fmt.Println(reflect.DeepEqual(m1, m2)) // true (맵은 키 순서 무관)
}
```

- 신기한 부분은 map 의 키 순서가 달라도 true 반환함 (맵은 순서가 달라도 상관 없어서 같다고 판단함)


<br />

### TLV (type-Length-Value)

- TLV 인코딩 체계는 가변 길이의 데이터를 처리하기 위한 방법

- 데이터 유형을 나타내는 정해진 길이의 바이트, 값의 크기를 나타내는 정해진 길이의 바이트, 값 자체를 나타내는 가변 길이의 바이트로 표현됨

<br />

## 데이터를 읽고 쓰는 도중 에러 처리

- 데이터를 네트워크 연결에 쓰는 것은 항상 보장될 수 없음, 이유는 네트워크 연결이 항상 원활하다는 것을 보장하지 못하기 때문

- ```nErr.Temporary()``` 는 일시적인 네트워크 오류에 따른 에러 타입임

- 일시적인 네트워크 연결 오류로 데이터 전송이 실패하면 안되기 때문에 golang for 루프로 데이터 쓰기가 이루어짐

- 또한 안정적인 네트워크 애플리케이션을 만들기 위해 io.Copy, io.MultiWriter, io.TeeReader 함수를 제공하고 있음

<br />

### io.Copy

- 이 함수는 io.Reader 에서 데이터를 읽어서 io.Writer 로 데이터를 쓸 수 있는 함수

- io.Reader 와 io.Writer 인터페이스를 둘다 포함하고 있음으로 두 노드 중간에서 데이터를 전송하는 proxy 를 생성하는데 유용함

- 이 함수는 출발 노드에서 목적지 노드로 데이터를 복제하는 역할을 함