package core

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"testing"
)

/*
func TestConsumnws(t *testing.T){
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket := gowebsocket.New("ws://localhost:4444");

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server")
		interrupt <- os.Interrupt
		return
	};

	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");
	};

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Println("Recieved connect error ", err)
	};

	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		log.Println("Recieved message " + message)
	};

	socket.OnBinaryMessage = func(data [] byte, socket gowebsocket.Socket) {
		log.Println("Recieved binary data ", data)
	};

	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved ping " + data)
	};

	socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Recieved pong " + data)
	};

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			socket.Close()
			return
		}
	}

}*/

func TestConsumnWs(t *testing.T) {
	runtime.GOMAXPROCS(4)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var channel_int chan int
	channel_int = make(chan int)

	var i = 0
	for i < 1 {
		log.Println(i)
		go func() {
			ws := NewWS("ws://localhost:4444")
			ws.Connect()

			go ws.Auth()

			ws.Subscribe("prueba")
			connections := ws.GetConnections()
			fmt.Println(connections)

		}()

		i++
	}
	<-channel_int
	//ws.SendMessage("texto")

}
