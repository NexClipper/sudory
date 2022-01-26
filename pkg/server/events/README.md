# evnets

## event

- 서버에서 요청에 대해 응답 할 때 이벤트가 발생 한다
- 요청이 오면 http header path에서 리스너 패턴을 매칭하여 이밴트가 발생함
- 발생된 이벤트는 리스너 타입에 따라 이벤트 처리

## listener type

- file
- webhook

## config listener

- 설정 파일의 정의
  - 공통 설정
    - events: 이벤트 정의 시작
    - listener: 리스터 정의 시작
    - type: 리스터 타입
    - name: 리스너 이름
    - pattern: http header path 정규식 매칭 패턴
    - option: 리스터 옵션 정의 시작
  - type-webhook
    - method: http method
    - url: http url
    - content-type: http content-type
    - timeout: http timeout
  - type-file
    - path: 파일 경로

```yaml
events: 
  - listener:
    type: 'webhook'
    name: 'client --all'
    pattern: '/client/*'
    option:
      method: post
      url: 'http://localhost:8000/'
      content-type: 'application/json'
      timeout: 30
  - listener:
    type: 'file'
    name: 'server --all'
    pattern: '/server/*'
    option:
      path: ../log/server.log
```

## example code

- 리스너 등록 코드

```golang
import "github.com/NexClipper/sudory/pkg/server/events"

func main() {
    //events new
    ecfg, err := events.New(*configPath)
    if err != nil {
        panic(err)
    }
    err = ecfg.Regist() //events Regist
    if err != nil {
        panic(err)
    }
    defer events.Manager.Stop() //event stop
    
    //process loop here
}
```

- 이벤트 발생 코드

```golang
import "github.com/NexClipper/sudory/pkg/server/events"

func middleware(ctx echo.Context) error {

    var err error
    var req, rsp interface{}

    //event invoke
    defer func() {
        events.Invoke(ctx, req, rsp, err)
    }()

    //request & response here
}

```

## listener type auto-generate

- go-enum 을 사용한 코드 생성

1. install go-enum
    - go-enum.install.sh
2. go generate
    - listener_type.go 파일에 정의된 go:generate 실행
