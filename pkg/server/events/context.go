package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

type ListenerContext interface {
	Type() string              //리스너 타입
	Name() string              //이름
	Pattern() string           //패턴
	Dest() string              //목적지
	Raise(v interface{}) error //이벤트 발생
}

func ErrorUndifinedEventListener(t string) error {
	return fmt.Errorf("undefined event listener type='%s'", t)
}

var Listeners map[string][]ListenerContext

var Manager *EventManager

var once sync.Once

func init() {

	once.Do(func() {
		Manager = NewManager(func(ctx ListenerContext, value interface{}, err error) {
			buf, _ := json.Marshal(value)
			log.Printf("event error handler name='%s' type='%s' dest='%s' value='%s' error='%v'\n", ctx.Name(), ctx.Type(), ctx.Dest(), string(buf), err)
		})
		go Manager.Start()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			Manager.Stop()
			log.Print("Interrupt")
		}()
	})

}
