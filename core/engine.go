package core

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/op/go-logging"
	"sync"
	"time"
)

var LogEngine = logging.MustGetLogger("Engine")

type IEngine interface {
	PublishSignal(signal *Signal)
	PublishMessage(message IMessage)
	Init(out chan IMessage, chan_signal chan *Signal)
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

func (re *RedisEngine) Init(out chan IMessage, chan_signal chan *Signal) {
	healthCheckPeriod := time.Duration(5)
	redisServerAddr := "localhost:6379"
	var err error
	re.c, err = redis.Dial("tcp", redisServerAddr, redis.DialReadTimeout(healthCheckPeriod+10*time.Second), redis.DialWriteTimeout(10*time.Second))
	if err != nil {
		fmt.Println(err)
	}

	re.psc = redis.PubSubConn{Conn: re.c}

	go re.subscribe()
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
func (re *RedisEngine) subscribe() {
	defer re.c.Close()
	re.active = true
	//re.status <- ENGINE

	err := re.psc.Subscribe(MESSAGE)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	} else {
		fmt.Println("Subscribed to Messages")
	}
	for {
		switch v := re.psc.Receive().(type) {
		case redis.Message:
			if v.Channel == SIGNAL {
				/*if signal.signal_type == DISCONNECT {
					re.active=false
					Log.Debugf("engine disconected by signal")
					break
				} else if signal.signal_type == TIMEOUT {
					Log.Debugf("engine disconected by timeout")
					break
				} else {
					break
				}*/
				fmt.Println(v.Data)
			} else if v.Channel == MESSAGE {
				fmt.Println(v.Data)
				//message := v.Data.(Message)
				//Log.Debugf("receive message by engine %v",.GetContent())
				//out <- message
			}
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
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

func (re *RedisEngine) Subscribe() {
	go re.subscribe()
}

func (re *RedisEngine) PublishSignal(s *Signal) {
	go func() {
		re.signal <- s
	}()

}

func (re *RedisEngine) PublishMessage(m IMessage) {
	channel_type := m.GetMessageType()
	message_json, _ := m.ToJson()
	err := re.c.Send(string(message_json), channel_type)
	if err != nil {
		fmt.Println(err)
	}
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

func (e *Engine) Init(out chan IMessage, chan_signal chan *Signal) {
	e.m.Lock()
	defer e.m.Unlock()
	go e.subscribe(out, chan_signal)
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

func (e *Engine) PublishMessage(m IMessage) {
	go func() {
		//Log.Debugf("message sended")
		e.buffer <- m
	}()

}
