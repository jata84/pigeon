package core

import (
	"testing"
	"time"
)

func TestRedisEngine(t *testing.T) {
	json := `{
 "comunication": {
  "type": "CONTENT",
  "namespaces": ["namespace.public.test","namespace.public.test2"],
  "client": {
   "uuid": "cdc9eda1-808d-4235-9974-6b8fd2cde84d",
   "datetime": "2018-12-09T18:05:16.263640434+01:00"
  }
 },
 "type": "MESSAGE",
 "content": "test_message",
 "timestamp": "2018-12-09T18:05:16.263650532+01:00"
}`
	message, err := NewMessageFromJson(json)
	if err != nil {
		t.Errorf("Error encoding Json")
		return
	}

	status := make(chan string)
	engine := NewRedisEngine(status)
	o := make(chan IMessage)
	s := make(chan *Signal)
	err = engine.Init(o, s)
	if err != nil{
		t.Errorf("Error initialization engine")
		return
	}

	err = engine.PublishMessage(message)
	if err != nil{
		t.Errorf(" %s Error publishing message \n",err)
		return
	}
	engine.GetStatus()
	time.Sleep(5000 * time.Millisecond)
	engine.Close()

}

/*
func TestMain(m *testing.M) {
	//log.SetOutput(&buf)
	m.Run()
}

func TestEngineStructure(t *testing.T) {

	//Engine
	status := make(chan string)
	engine := NewEngine(status)
	o := make(chan IMessage)
	s := make(chan *Signal)
	go engine.Init(o,s)
	<-status

	//Channel Structure
	channel := NewPublicChannel("public")
	channel_list := NewChannelList()
	channel_list.AddChannel(channel)
	channel_array := channel_list.ToArray()

	test_content_1 := "Mensaje de prueba"
	test_content_2 := "Mensaje de prueba 2"

	client := NewClient(nil,nil)

	//Message Structure
	message := NewMessage(channel_array, client, test_content_1)

	//Init Channels and engine subscription

	//Publish first message and receive it
	engine.PublishMessage(message)

	message_out := <-o
	if message_out.GetContent() != test_content_1 {
		t.Errorf("Incorrect", "error")
	}

	//Publish second message and receive it
	message_2 := NewMessage(channel_array, client, test_content_2)
	engine.PublishMessage(message_2)
	message_out_2 := <-o
	if message_out_2.GetContent() != test_content_2 {
		t.Errorf("Incorrect", "error")
	}

	signal := NewDisconnectSignal()
	engine.PublishSignal(signal)
	for engine.IsActive() {
	}
	t.Log(buf.String())

}

func BenchmarkEngineMessages(b *testing.B) {

	status := make(chan string)
	engine := NewEngine(status)
	o := make(chan IMessage)
	s := make(chan *Signal)
	engine.Init(o,s)
	<-status

	//Channel Structure
	channel := NewPublicChannel("public")
	channel_list := NewChannelList()
	channel_list.AddChannel(channel)
	channel_array := channel_list.ToArray()

	var message_sended, message_received []IMessage
	message_sended = make([]IMessage, 1000)
	message_received = make([]IMessage, 1000)

	client := NewClient(nil,nil)

	//Message Structure

	for n := 0; n < 1000; n++ {

		message := NewMessage(channel_array, client, "Message_test")
		engine.PublishMessage(message)
		message_sended = append(message_sended, message)
		log.Printf("Send message number %v \n", n)
	}
	for n := 0; n < 1000; n++ {
		m := <-o
		message_received = append(message_received, m)
		log.Printf("Received message number %v %v \n", n, m)
	}
	if len(message_received) != len(message_sended) {
		b.Errorf("Not received all messages", "error")
	}
	signal := NewDisconnectSignal()
	engine.PublishSignal(signal)
	for engine.IsActive() {
	}
	//b.Log(buf.String())
}

func TestEngineBenchmarkConcurrentMessages(t *testing.T) {

}
*/
