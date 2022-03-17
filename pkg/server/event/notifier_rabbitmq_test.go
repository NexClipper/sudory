package event_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/streadway/amqp"
)

func Dial(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn, nil
}

type OptQueueDeclare struct {
	//lint:ignore U1000 ignore this!
	reserved1 uint16
	Queue     string
	Passive   bool
	Durable   bool
	// Exclusive  bool
	AutoDelete bool
	NoWait     bool
	Arguments  map[string]interface{}
}

func Codinator(url string, opt *OptQueueDeclare) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				panic(err)
			}
		}
	}()
	panic := func(err error) {
		if err == nil {
			return
		}
		panic(err)
	}

	conn, err := amqp.Dial(url)
	panic(err)
	defer conn.Close()

	ch, err := conn.Channel()
	panic(err)

	_, err = ch.QueueDeclare(
		opt.Queue,
		opt.Durable,
		opt.AutoDelete,
		false,
		opt.NoWait,
		opt.Arguments,
	)
	panic(err)

	fmt.Fprintf(os.Stdout, "rabbitmq queue declare queue=%v durable=%v auto-delete=%v no-wait=%v arguments=%v",
		opt.Queue, opt.Durable, opt.AutoDelete, opt.NoWait, opt.Arguments)

	return
}

func Send(url, queue string, payload bytes.Buffer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				panic(err)
			}
		}
	}()
	panic := func(err error) {
		if err == nil {
			return
		}
		panic(err)
	}

	conn, err := amqp.Dial(url)
	panic(err)
	defer conn.Close()

	ch, err := conn.Channel()
	panic(err)

	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload.Bytes(),
		},
	)
	panic(err)

	// fmt.Fprintf(os.Stdout, "rabbitmq send payload queue=%s", queue)
	return
}

func TestCodinator(t *testing.T) {

	url := "amqp://localhost:5672/"
	qname := "T-Queue"

	err := Codinator(url, &OptQueueDeclare{
		Queue:      qname,
		Durable:    true,
		AutoDelete: false,
		NoWait:     false,
		Arguments:  nil,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSend(t *testing.T) {

	url := "amqp://localhost:5672/"
	qname := "T-Queue"

	buf := bytes.Buffer{}
	buf.Write([]byte("hello rabbitmq"))
	err := Send(url, qname, buf)
	if err != nil {
		t.Fatal(err)
	}

}
