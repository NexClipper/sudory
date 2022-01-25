package events

import (
	"context"
	"encoding/json"
	"log"
	"regexp"
	"runtime"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type EventArgs struct {
	Sender string
	Args   interface{}
}

type EventManager struct {
	Stop         context.CancelFunc
	ErrorHandler ErrorHandler
	sender       *SafeChannel
}

type ErrorHandler func(ctx ListenerContext, value interface{}, err error)

func defaultErrorHandler(ctx ListenerContext, value interface{}, err error) {
	buf, _ := json.Marshal(value)
	log.Printf("event error handler name='%s' type='%s' dest='%s' value='%s' error='%v'\n", ctx.Name(), ctx.Type(), ctx.Dest(), string(buf), err)
}

func NewManager(handler ErrorHandler) *EventManager {
	if handler == nil {
		handler = defaultErrorHandler
	}

	return &EventManager{
		ErrorHandler: handler,
		sender:       NewSafeChannel(0), //new sender
	}
}

func Invoke(ctx echo.Context, req, rsp interface{}, err error) {

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
	}
	if err != nil {
		args["error"] = err.Error()
	}

	Manager.sender.SafeSend(args)
}

func (me *EventManager) Start() {
	wg := sync.WaitGroup{}

	//new reciver
	reciver := NewSafeChannel(0)

	closing := make(chan struct{})
	closed := make(chan struct{})

	//set stop
	me.Stop = func() {
		select {
		case closing <- struct{}{}:
			<-closed
		case <-closed:
		}
	}

	//sender
	go func() {
		defer func() {
			close(closed)
			reciver.SafeClose()
		}()

		for {
			select {
			case <-closing:
				return
			default:
			}

			select {
			case <-closing:
				return
			case v := <-me.sender.C:
				reciver.SafeSend(v)
			}
		}
	}()

	//reciver
	num_cpu := runtime.NumCPU()
	wg.Add(num_cpu) // wg add
	for n := 0; n < num_cpu; n++ {
		go func() {
			defer wg.Done() //wg done

			for v := range reciver.C {
				if args, ok := v.(map[string]interface{}); ok {
					notify(args, me.ErrorHandler) //notify
				}
			}
		}()
	}

	wg.Wait() //wg wait
}

func notify(args map[string]interface{}, err_handler ErrorHandler) {

	foreach := func(ctxs []ListenerContext, fn func(ctx ListenerContext)) {
		for _, listener := range ctxs {
			fn(listener)
		}
	}

	path := args["path"].(string)
	for pattern, ctxs := range Listeners {

		rx, _ := regexp.Compile(pattern)

		ok := rx.Match([]byte(path))
		if ok {
			foreach(ctxs, func(ctx ListenerContext) {
				err := ctx.Raise(map[string]interface{}{
					"name": ctx.Name(),
					"type": ctx.Type(),
					"time": time.Now(),
					"args": args,
				})
				if err != nil {
					err_handler(ctx, args, err)
				}
			})
		}

	}
}
