package client

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
)

type FlagParam []string

func (f *FlagParam) String() string {
	return "string method"
}

func (f *FlagParam) Set(value string) error {
	*f = strings.Split(value, ",")
	return nil
}

var filters FlagParam

type WebsocketClient struct {
	url url.URL
	ws  *websocket.Conn
}

func NewWebsocketClient(host string, port string) *WebsocketClient {
	var addr = flag.String("addr", fmt.Sprintf("%s:%s", host, port), "ApiRouter entrypoint")
	var u = url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: ""}

	return &WebsocketClient{
		url: u,
	}
}

func (w *WebsocketClient) Send(json_message []byte) {

}

func (w *WebsocketClient) SendRegistration() {
	w.ws.WriteMessage(websocket.TextMessage, []byte(`{
    				"timestamp": "2018-10-17T19:51:54.126486939+02:00",
  					"type": "MESSAGE",
 					"comunication": {
    					"type": "REGISTRATION",
    					"channels": [],
    					"client": {}
  					},
  				"data": {}
				}`))
}

func (w *WebsocketClient) SendSubscription(namespace string) {
	w.ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{
    "timestamp":"2018-10-17T19:51:54.126486939+02:00",
    "type":"MESSAGE",
    "comunication":{
      "type":"SUBSCRIPTION",
      "namespaces":[{"type":"PRIVATE","name":"%s"
           }],
      "client":{}
    },
    "data":{}
  }`, namespace)))
}

func (w *WebsocketClient) Init() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	req, err := http.NewRequest("GET", w.url.String(), nil)
	req.URL = &w.url

	dialer := websocket.Dialer{}

	var resp *http.Response
	w.ws, resp, err = dialer.Dial(w.url.String(), req.Header)
	if err != nil {
		return
	} else if resp.StatusCode == http.StatusSwitchingProtocols {

		//log.Print("connected, watching event now:")

	}

	var data []byte
	if resp.ContentLength > 0 {
		defer resp.Body.Close()
		data, _ = ioutil.ReadAll(resp.Body)
	}

	if err != nil {
		log.Fatalf("Error Result:\n status    : %v\n return    : %v", resp.Status, string(data))
	}

	go func() {
		for {
			_, message, err := w.ws.ReadMessage()
			if err != nil {
				break
			}
			if string(message) == `{"IMessage":null,"message_type":"Message","channels":[{"type":"PROTECTED","name":"ALL"}],"client":null,"content":"Test","content_type":"content_message"}` {
				w.ws.Close()
				break
			}

		}
	}()
}

//format filter to json string
func formatFilter(filters *FlagParam) (*string, error) {
	result := map[string]map[string]bool{}
	for _, v := range *filters {
		//log.Printf("[Debug] original filter: %v", v)
		item := strings.SplitN(v, "=", 2)
		lenItem := len(item)
		switch {
		case item[0] == "container" || item[0] == "image" || item[0] == "event" || item[0] == "label":
			if lenItem == 1 {
				return nil, errors.New(fmt.Sprintf("Wrong filter format for [%v]", item[0]))
			} else {
				mm, ok := result[item[0]]
				if !ok {
					mm = make(map[string]bool)
					result[item[0]] = mm
				}
				mm[item[1]] = true
			}
		case item[0] == "":
			return nil, errors.New("filter name can not be empty")
		default:
			return nil, errors.New(fmt.Sprintf("filter only support container,image,label,event"))
		}
	}
	b, _ := json.Marshal(result)
	strResult := string(b)
	return &strResult, nil
}
