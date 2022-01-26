package events

import (
	"regexp"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
)

type empty struct{}

type Waiter interface {
	Wait()
}

type WaitGroupHelper struct {
	mux sync.Mutex
	wg  sync.WaitGroup
}

var _ Waiter = (*EventContext)(nil)

func (wgh *WaitGroupHelper) Add(n int) {
	wgh.mux.Lock()
	defer wgh.mux.Unlock()

	wgh.wg.Add(n)
}

func (wgh *WaitGroupHelper) Done() {
	wgh.wg.Done()
}

func (wgh *WaitGroupHelper) Wait() {
	wgh.mux.Lock()
	defer wgh.mux.Unlock()

	wgh.wg.Wait()
}

type EventContexter interface {
	Name() string
	Pattern() string
	ListenerContexts() []ListenerContexter
	Raise(args EventArgs, error_handler ErrorHandler)
}

type EventContext struct {
	wg WaitGroupHelper

	opt       EventConfig
	listeners []ListenerContexter
}

var _ EventContexter = (*EventContext)(nil)
var _ Waiter = (*EventContext)(nil)

func NewEventContext(opt EventConfig, ctx ...ListenerContexter) *EventContext {
	return &EventContext{opt: opt, listeners: ctx}
}

func (ctx *EventContext) Name() string {
	return ctx.opt.Name
}

func (ctx *EventContext) Pattern() string {
	return ctx.opt.Pattern
}
func (ctx *EventContext) BuzyTimeout() time.Duration {
	const buzy_timeout int = 10

	if ctx.opt.BuzyTimeout == nil {
		return time.Duration(buzy_timeout) * time.Second
	}
	if *ctx.opt.BuzyTimeout <= 0 {
		return time.Duration(buzy_timeout) * time.Second
	}

	return time.Duration(*ctx.opt.BuzyTimeout) * time.Second
}

func (ctx *EventContext) ListenerContexts() []ListenerContexter {
	return ctx.listeners
}
func (ctx *EventContext) Raise(event_args EventArgs, error_handler ErrorHandler) {

	rx, _ := regexp.Compile(ctx.Pattern())

	ok := rx.Match([]byte(event_args.Sender)) //matching sender
	if !ok {
		return //not matched
	}

	wg := sync.WaitGroup{}
	wg.Add(len(ctx.listeners)) //add gorutine length

	ctx.wg.Add(len(ctx.listeners)) //WaitGroup Add

	for _, listener := range ctx.listeners {

		go func(listener_ ListenerContexter, event_args_ *EventArgs, error_handler_ ErrorHandler) {
			wg.Done() //gorutine actived

			defer func() {
				ctx.wg.Done() //WaitGroup Done
			}()
			buzy := onTime(func() {
				err := listener_.Raise(event_args_.Args)
				if err != nil {
					buf, _ := serialize_json(event_args_.Args)
					error_handler_("failed to notify name='%s' type='%s' value='%s' error='%v'\n", listener_.Name(), listener_.Type(), buf.String(), err)
				}
			}, ctx.BuzyTimeout())

			if buzy {
				log.Error("Buzy")
			}
		}(listener, &event_args, error_handler)
	}

	wg.Wait() //wait gorutine active
}

func (ctx *EventContext) Wait() {
	ctx.wg.Wait()
}

func onTime(fn func(), buzy_timeout time.Duration) bool {

	var buzy bool = false

	ch := make(chan empty, 1)
	go func() {
		fn()          //call fn
		ch <- empty{} // after channel active
	}()
	for {
		select {
		case <-ch:
			return buzy
		case <-time.After(buzy_timeout):
			// call timed out
			buzy = true
		}
	}
}
