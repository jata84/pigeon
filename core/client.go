package core

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/orcaman/concurrent-map"
	"log"
	"time"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

/*
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
*/
type Client struct {
	Id         uuid.UUID `json:"uuid"`
	conn       *websocket.Conn
	send       chan []byte
	router     *Router
	registered bool
	Time       time.Time `json:"datetime"`
	//node_list *NodeList
	node_list []string
}

func (c *Client) GetClientNode() []string {
	return c.node_list
}

func (c *Client) AddNamespace(namespace string) {
	c.node_list = append(c.node_list, namespace)
}

func (c *Client) RemoveNamespace(namespace string) {
	for k, n := range c.node_list {
		if n == namespace {
			c.node_list = c.node_list[:k+copy(c.node_list[k:], c.node_list[k+1:])]
			break
		}
	}
}

func (c *Client) Read() {
	defer func() {

	}()

	c.conn.SetReadLimit(maxMessageSize)
	//c.conn.SetReadDeadline(time.Now().Add(pongWait))
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		message_json, err := NewMessageFromJson(string(message))
		c.SendMessage(message_json)
		if err != nil {
			log.Printf("error %v", err)
		}
	}
}
func (c *Client) GetRegistered() bool {
	return c.registered
}

func (c *Client) SetRegistered(reg bool) {
	c.registered = reg
}
func (c *Client) ToJson() string {
	json_data, _ := json.Marshal(c)
	return string(json_data)
}

func (c *Client) Write(m IMessage) {
	if c.conn != nil {
		m_json, _ := m.ToJson()
		error := c.conn.WriteMessage(websocket.TextMessage, m_json)

		if error != nil {
			Log.Debugf("Error writting communication to websocket %v", error)
		}
	} else {
		Log.Debugf("Error no connection available to send message")
	}
}

func (c *Client) SendMessage(m IMessage) {
	c.router.SendMessage(m)
}

func (c *Client) Close() {
	c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(time.Duration(1))
	c.conn.Close()

}

func NewClient(conn *websocket.Conn, r *Router) *Client {
	client_id := uuid.New()
	if conn != nil {
		conn.SetCloseHandler(func(code int, text string) error {
			client_registered := r.registered_clients.GetById(client_id)
			client_unregistered := r.unregistered_clients.GetById(client_id)
			var client *Client
			if client_registered == nil {
				client = client_unregistered
			} else {
				client = client_registered
			}
			if client == nil {
				return nil
			}
			for _, node := range client.node_list {
				namespace_node := r.nodes.FindNamespaceNode(node)
				if namespace_node != nil {
					namespace_node.RemoveClient(client)
					//TODO revisar esta parte
					if namespace_node.clientList.Count() == 0 {
						previous_node_list := namespace_node.PreviousNode
						previous_node_list.nodes.Remove(namespace_node.NodeName)
					}
				}
			}
			r.registered_clients.RemoveClient(client)
			r.unregistered_clients.RemoveClient(client)
			return nil
		})
	}

	return &Client{
		Id:         client_id,
		conn:       conn,
		registered: false,
		send:       make(chan []byte),
		router:     r,
		Time:       time.Now(),
		node_list:  []string{},
	}
}

type ClientList struct {
	clients cmap.ConcurrentMap
}

func NewClientList() *ClientList {
	return &ClientList{
		clients: cmap.New(),
	}
}

func (c *ClientList) Count() int {
	return c.clients.Count()
}

func (c *ClientList) ToJson() string {
	json_data, _ := json.Marshal(c.ToArray())
	return string(json_data)
}

func (c *ClientList) ClientExists(client *Client) bool {
	if _, ok := c.clients.Get(client.Id.String()); ok {
		return true
	} else {
		return false
	}
}

func (c *ClientList) ConnectionExists(conn *websocket.Conn) bool {
	for item := range c.clients.IterBuffered() {
		client := item.Val
		if conn == client.(*Client).conn {
			return true
		}
	}
	return false
}

func (c *ClientList) ToArray() []*Client {

	client_list := make([]*Client, 0)
	for item := range c.clients.IterBuffered() {
		client := item.Val
		client_list = append(client_list, client.(*Client))
	}
	return client_list
}

func (c *ClientList) SendMessage(m IMessage) {
	for item := range c.clients.IterBuffered() {
		client := item.Val
		client.(*Client).Write(m)
	}

}

func (c *ClientList) SearchById(id uuid.UUID) (client *Client) {

	if cli, ok := c.clients.Get(id.String()); ok {
		return cli.(*Client)
	} else {
		return nil
	}

}

func (c *ClientList) GetById(id uuid.UUID) (client *Client) {

	if cli, ok := c.clients.Get(id.String()); ok {
		return cli.(*Client)
	} else {
		return nil
	}

}

func (c *ClientList) RemoveClient(client *Client) {
	Log.Debugf("clients %v", c.ToJson())
	c.clients.Remove(client.Id.String())

}

func (c *ClientList) AddClient(client *Client) {

	if _, ok := c.clients.Get(client.Id.String()); ok {

	} else {
		c.clients.Set(client.Id.String(), client)
	}
}
