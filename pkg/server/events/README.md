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
  - events: 이벤트 정의 배열
  - event: 이벤트 정의 시작
  - name: 이벤트 이름
  - pattern: EventArgs.Sender와 정규식 패턴 매칭
  - listeners: 리스너 정의 배열
  - listener: 리스너 정의 시작
  - type: 리스너 타입
  - config listener-type-webhook:
    - method: http method
    - url: http url
    - content-type: http content-type
    - timeout: http timeout
  - config listener-type-file:
    - path: 파일 경로

```yaml
events:
  - event:
    name: 'client poll'
    pattern: '/client/poll'
    listeners:
    - listener:
      type: 'webhook'
      method: post
      url: 'http://localhost:8000/'
      content-type: 'application/json'
      timeout: 30  
    - listener: 
      type: 'file'
      path: ./log/server.log
  - event:
    name: 'client --all'
    pattern: '/client/*'
    listeners:
    - listener:
      type: 'webhook'
      method: post
      url: 'http://localhost:8000/'
      content-type: 'application/json'
      timeout: 30
```

## example code

- 리스너 등록 코드

```golang
import "github.com/NexClipper/sudory/pkg/server/events"

func main() {
    //process init here
    
    //events
    var eventContexts []events.EventContexter
    var eventConfig *events.Configs
    //event config
    if eventConfig, err = events.NewConfig(*configPath); err != nil { //config file load
        panic(err)
    }
    //event config vaild
    if err = eventConfig.Vaild(); err != nil { //config vaild
        panic(err)
    }
    //event config make listener
    if eventContexts, err = eventConfig.MakeEventListener(); err != nil { //events regist listener
        panic(err)
    }
    //event manager
    eventInvoke := channels.NewSafeChannel(0)
    manager := events.NewManager(eventContexts, log.Printf)
    deactivate := manager.Activate(eventInvoke, len(eventContexts)) //manager activate
    defer func() {
        deactivate() //stop when closing
    }()
    events.Invoke = func(v *events.EventArgs) { eventInvoke.SafeSend(v) } //setting invoker

    //process loop here
}
```

- 이벤트 발생 코드

```golang
import "github.com/NexClipper/sudory/pkg/server/events"

func fn(ctx echo.Context) error {

    var err error
    var req, rsp interface{}

    //event invoke
    defer func() {

        path := ctx.Request().URL.Path
        method := ctx.Request().Method
        status := ctx.Response().Status
        query := ctx.QueryString()
  
        args := map[string]interface{}{
            "path":    path,
            "query":   query,
            "method":  method,
            "reqbody": req,
            "rspbody": rsp,
            "status":  status,
            "error":   err,
        }
  
        if err == nil {
            delete(args, "error")
        }
        events.Invoke(&events.EventArgs{Sender: path, Args: args})
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
