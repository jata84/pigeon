package core

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sacOO7/gowebsocket"
	"gopkg.in/resty.v1"
	"log"
)

type Connections struct {
	Uuid     uuid.UUID `json:uuid`
	Datetime string    `json:datetime`
}

type WS struct {
	subcription_channel chan *Client
	Id                  uuid.UUID
	socket              gowebsocket.Socket
	OnTextMessage       func(message string, socket gowebsocket.Socket, id *uuid.UUID)
}

func (ws *WS) SetId(id uuid.UUID) {
	ws.Id = id
}

func (ws *WS) GetId() uuid.UUID {
	return ws.Id
}

func (ws *WS) onDisconnected(err error, socket gowebsocket.Socket) {
	log.Println("Disconnected from server")
	return
}

func (ws *WS) onConnected(socket gowebsocket.Socket) {
	log.Println("Connected to server")
}

func (ws *WS) onConnectError(err error, socket gowebsocket.Socket) {
	log.Println("Connected Error")
}

func (ws *WS) Subscribe(name string) bool {
	client := <-ws.subcription_channel
	namespace_node := stringToNamespaces("-prueba.test.mensajes")
	message := NewSubscriptionMessage(namespace_node, client)
	json_message, _ := message.ToJson()
	resp, err := resty.R().SetHeader("Authorization", "test").SetBody(json_message).Post("http://localhost:8084/api/subscribe")
	if err != nil {
		return false
	} else {
		println(resp)
		return false
	}
}

func (ws *WS) SendMessage(message string) bool {
	namespace_node := stringToNamespaces("-prueba.test.mensajes")
	client := Client{Id: ws.Id}
	message_json := NewContentMessage(namespace_node, &client, message, "")
	json_message, _ := message_json.ToJson()
	resp, err := resty.R().SetHeader("Authorization", "test").SetBody(json_message).Post("http://localhost:8084/api/message")
	if err != nil {
		return false
	} else {
		println(resp)
		return false
	}
}

func (ws *WS) GetConnections() []Connections {
	resp, err := resty.R().SetHeader("Authorization", "test").Get("http://localhost:8084/api/clients")
	if err == nil {

		if resp.StatusCode() == 200 {
			json_response := resp.Body()
			keys := make([]Connections, 0)
			json.Unmarshal(json_response, &keys)
			return keys
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (ws *WS) KillConnection(connection string) bool {
	url := fmt.Sprintf("http://localhost:8084/api/clients/%s", connection)
	resp, err := resty.R().SetHeader("Authorization", "test").Delete(url)
	if err == nil {

		if resp.StatusCode() == 200 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (ws *WS) onBinaryMessage(data []byte, socket gowebsocket.Socket) {
	log.Println("Binary Message received")
}

func (ws *WS) Connect() {
	ws.socket.Connect()
	ws.subcription_channel = make(chan *Client)
	ws.socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {

		message_json, _ := NewMessageFromJson(message)
		if message_json.Comunication.ComunicationType == REGISTRATION {
			ws.Id = message_json.Comunication.GetClient().Id
			ws.subcription_channel <- &Client{Id: ws.Id}
			log.Println(message_json.Comunication.ComunicationType)
			log.Println(ws.Id)
		} else if message_json.Comunication.ComunicationType == SUBSCRIPTION {
			log.Println(message_json.Comunication.ComunicationType)
		} else if message_json.Comunication.ComunicationType == MESSAGE {
			log.Println(message_json.Comunication.ComunicationType)
			log.Println(ws.Id)
			log.Println(message_json.Data)
		}

	}

}

func (ws *WS) registration(text string) {
	ws.socket.SendText(text)
}

func (ws *WS) Auth() {
	message := NewRegistrationMessage(nil)
	message_json, _ := message.ToJson()
	ws.registration(string(message_json))
}

func NewWS(url string) WS {
	ws := WS{
		socket: gowebsocket.New(url),
		Id:     uuid.UUID{},
	}

	ws.socket.OnDisconnected = ws.onDisconnected
	ws.socket.OnConnected = ws.onConnected
	ws.socket.OnConnectError = ws.onConnectError
	ws.socket.OnBinaryMessage = ws.onBinaryMessage

	return ws

}
