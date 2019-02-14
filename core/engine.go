package core

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/op/go-logging"
	"sync"
	"time"
)

var LogEngine = logging.MustGetLogger("Engine")

type IEngine interface {
	PublishSignal(signal *Signal)
	PublishMessage(message IMessage) error
	Init(out chan IMessage, chan_signal chan *Signal) error
	GetStatus() bool
}

type RedisEngine struct {
	m   sync.Mutex
	c   redis.Conn
	psc redis.PubSubConn
	IEngine
	status chan string
	active bool
	buffer chan IMessage
	signal chan *Signal
}

func NewRedisEngine(status chan string) *RedisEngine {
	return &RedisEngine{
		status: status,
		active: true,
		buffer: make(chan IMessage),
		signal: make(chan *Signal),
	}
}

func (re *RedisEngine) Init(out chan IMessage, chan_signal chan *Signal) error {
	//healthCheckPeriod := time.Duration(5)
	redisServerAddr := "localhost:6379"
	var err error
	//, redis.DialReadTimeout(healthCheckPeriod+10*time.Second), redis.DialWriteTimeout(10*time.Second)
	re.c, err = redis.Dial("tcp", redisServerAddr)
	if err != nil {
		return err
	}
	re.psc = redis.PubSubConn{Conn: re.c}
	re.Subscribe(out, chan_signal)
	return nil
}

func (re *RedisEngine) Close() {
	re.active = false
}

func (re *RedisEngine) GetStatus() bool {
	return re.active
}

func (re *RedisEngine) IsActive() bool {
	return re.active
}

func (re *RedisEngine) unsubscribe() {
	re.active = false

}
func (re *RedisEngine) subscribe(out chan IMessage, chan_signal chan *Signal) {
	defer re.c.Close()
	re.active = true
	re.status <- ENGINE
	for {
		switch v := re.psc.Receive().(type) {
		case redis.Message:
			fmt.Println(v.Channel)
			if v.Channel == SIGNAL {
				fmt.Printf("%s: SIGNAL: %s\n", v.Channel, v.Data)
			} else if v.Channel == MESSAGE {
				message, json_error := NewMessageFromJson(string(v.Data))
				if json_error != nil {
					Log.Debugf("error decoding message %s", json_error)
				} else {

					Log.Debugf("receive message by engine %v", message)
					out <- message

				}

			} else if v.Channel == NOTIFICATION {
				message, json_error := NewMessageFromJson(string(v.Data))
				if json_error != nil {
					Log.Debugf("error decoding message %s", json_error)
				} else {
					Log.Debugf("receive message by engine %v", message)
					out <- message

				}

			}

		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)

		case error:
			fmt.Println(v)
		}
	}

	/*
		re.active = true
		re.status <- ENGINE
		defer re.unsubscribe()
		for {
			select {
			case message := <-re.buffer:
				Log.Debugf("receive message by engine %v",message.GetContent())
				out <- message
				continue
			case signal := <-re.signal:
				if signal.signal_type == DISCONNECT {
					re.active=false
					Log.Debugf("engine disconected by signal")
					break
				} else if signal.signal_type == TIMEOUT {
					Log.Debugf("engine disconected by timeout")
					break
				} else {
					break
				}
			}
			break
		}
		if (re.active){
			re.active = false
		}else{
			if (Configuration.Main.Gracetime > 0){
				Log.Info("Waiting grace period...")
				time.Sleep(10000 * time.Millisecond)
				chan_signal <- NewDisconnectSignal()

			}
		}

		Log.Info("Closing Engine")
	*/
}

func (re *RedisEngine) Subscribe(out chan IMessage, chan_signal chan *Signal) {
	err := re.psc.Subscribe(MESSAGE, NOTIFICATION, SIGNAL)

	if err != nil {
		Log.Debugf("error subscribing redis %s", err)
	}
	go re.subscribe(out, chan_signal)
}

func (re *RedisEngine) PublishSignal(s *Signal) {
	go func() {
		re.signal <- s
	}()

}

func (re *RedisEngine) PublishMessage(m IMessage) error {
	channel_type := m.GetMessageType()
	message_json, _ := m.ToJson()
	redisServerAddr := "localhost:6379"
	push_connection, err := redis.Dial("tcp", redisServerAddr)
	defer push_connection.Close()
	if re.c == nil {
		return errors.New("Connection closed")
	}
	_, err = push_connection.Do("PUBLISH", channel_type, string(message_json))
	//err := re.c.Send(, channel_type)
	if err != nil {
		return err
	}
	return nil
}

type Engine struct {
	m sync.Mutex
	IEngine
	status chan string
	active bool
	buffer chan IMessage
	signal chan *Signal
}

func NewEngine(status chan string) *Engine {
	return &Engine{
		status: status,
		active: true,
		buffer: make(chan IMessage),
		signal: make(chan *Signal),
	}
}

func (e *Engine) Init(out chan IMessage, chan_signal chan *Signal) error {
	e.m.Lock()
	defer e.m.Unlock()

	go e.subscribe(out, chan_signal)

	return nil
}

func (e *Engine) Close() {
	e.active = false
}

func (e *Engine) GetStatus() bool {
	//e.m.Lock()
	//defer e.m.Unlock()
	return e.active
}

func (e *Engine) IsActive() bool {
	return e.active
}

func (e *Engine) unsubscribe() {
	e.active = false

}
func (e *Engine) subscribe(out chan IMessage, chan_signal chan *Signal) {
	e.active = true
	e.status <- ENGINE
	defer e.unsubscribe()
	for {
		select {
		case message := <-e.buffer:
			Log.Debugf("receive message by engine %v", message.GetContent())
			out <- message
			continue
		case signal := <-e.signal:
			if signal.signal_type == DISCONNECT {
				e.active = false
				Log.Debugf("engine disconected by signal")
				break
			} else if signal.signal_type == TIMEOUT {
				Log.Debugf("engine disconected by timeout")
				break
			} else {
				break
			}
		}
		break
	}
	if e.active {
		e.active = false
	} else {
		if Configuration.Main.Gracetime > 0 {
			Log.Info("Waiting grace period...")
			time.Sleep(10000 * time.Millisecond)
			chan_signal <- NewDisconnectSignal()

		}
	}

	Log.Info("Closing Engine")

}

func (e *Engine) Subscribe() {

}

func (e *Engine) PublishSignal(s *Signal) {
	go func() {
		e.signal <- s
	}()

}

func (e *Engine) PublishMessage(m IMessage) error {
	go func() {
		//Log.Debugf("message sended")
		e.buffer <- m
	}()
	return nil
}
