package core

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
	"time"
)

//NAMESPACES
const (
	ROOT   = "ROOT"
	PUBLIC = "PUBLIC"
)

const (
	NODE     = "NODE"
	WILDCARD = "WILDCARD"
	CHANNEL  = "CHANNEL"
	TIMEOUT  = "TIMEOUT"
)

const (
	REDIS  = "REDIS"
	MEMORY = "MEMORY"
)

type Signal struct {
	signal_type    string
	signal_message string
}

func NewDisconnectSignal() *Signal {
	return &Signal{
		signal_type:    DISCONNECT,
		signal_message: "Disconnection signal",
	}
}

/**
func NewTimeoutSignal() *Signal {
	return &Signal{
		signal_type:    TIMEOUT,
		signal_message: "Timeout signal",
	}
}
*/
type IRouter interface {
}

type Router struct {
	sync.Mutex
	status               chan string
	active               bool
	signals              chan *Signal
	input                chan IMessage
	nodes                *NodeTree
	registered_clients   *ClientList
	unregistered_clients *ClientList
	engine               IEngine
}

func (r *Router) IsActive() bool {
	return r.active
}

func NewRouter(status chan string) (*Router, error) {

	engine_type := Configuration.Engine.Engine_type
	var engine_temp IEngine
	if engine_type == REDIS {
		engine_temp = NewRedisEngine(status)
	} else if engine_type == MEMORY {
		engine_temp = NewEngine(status)
	} else {
		return nil, errors.New("Engine not supported")
	}

	return &Router{
		status:               status,
		active:               true,
		registered_clients:   NewClientList(),
		unregistered_clients: NewClientList(),
		signals:              make(chan *Signal),
		input:                make(chan IMessage),
		nodes:                NewNodeTree(),
		engine:               engine_temp,
	}, nil
}

func (r *Router) GetClient(uuid uuid.UUID) *Client {
	return r.registered_clients.GetById(uuid)
}

func (r *Router) GetClientUnregistered(uuid uuid.UUID) *Client {
	return r.unregistered_clients.GetById(uuid)
}

func (r *Router) SendSignal(s *Signal) {
	go func() {
		r.engine.PublishSignal(s)
	}()
}

func (r *Router) Init() {
	//log.Printf("Init Router")
	r.engine.Init(r.input, r.signals)
	//root_namespace := (ROOT)
	//r.channels.AddChannel(root_namespace)
	//go r.unregisteredTimeout()
	go r.routing()
}

func (r *Router) routingCriteria(m IMessage) {
	message_type := m.GetMessageType()
	communication_type := m.GetCommunicationType()

	client := m.GetClient()
	var client_unregistered, client_registered *Client
	if client != nil {
		client_unregistered = r.GetClientUnregistered(client.Id)
		client_registered = r.GetClient(client.Id)
	}

	if (client_unregistered != nil) && (message_type == MESSAGE) && (communication_type == REGISTRATION) {
		if r.unregistered_clients.ClientExists(client) {
			client = client_unregistered
			r.registered_clients.AddClient(client)
			r.unregistered_clients.RemoveClient(client)
			client.SetRegistered(true)
			r.SendMessage(NewNotificationRegistration(client, WEBSOCKET))
			r.SendMessage(NewNotificationInfo(ADMIN_CHANNEL, nil, fmt.Sprintf("user registered %v", client.Id.String()), WEBSOCKET))

			Log.Debugf("user registered %v", client.Id.String())
		}

	} else if (client_registered != nil) && (message_type == MESSAGE) && (communication_type == SUBSCRIPTION) {

		var client *Client
		client = client_registered
		if client == nil {
			Log.Errorf("message does not contain client to subscribe")
			return
		}

		if client.registered {
			node_list := m.GetNamespacesNames()
			r.Subscribe(node_list, client)
			notification := NewNotificationSubscription(client, WEBSOCKET, node_list)
			r.SendMessage(notification)
			r.SendMessage(NewNotificationInfo(ADMIN_CHANNEL, nil, fmt.Sprintf("user subscribed %v", client.Id.String()), WEBSOCKET))
			Log.Infof("user %v subscribed to %v", client.Id.String(), m.GetNamespacesNames())
		} else {
			Log.Errorf("client not registered %v, cannot create subscription", client.Id.String())
		}
	} else if (client_registered != nil) && (message_type == MESSAGE) && (communication_type == UNSUBSCRIPTION) {
		client = client_registered
		if client.registered {
			node_list := m.GetNamespacesNames()
			r.Unsubscribe(node_list, client)
			notification := NewNotificationUnsubscription(client, WEBSOCKET, node_list)
			r.SendMessage(notification)
			Log.Infof("user %v unsubscribed to %v", client.Id.String(), m.GetNamespacesNames())
		} else {
			Log.Errorf("client not registered %v, cannot create subscription", client.Id.String())
		}
	} else if (message_type == MESSAGE) && (communication_type == MESSAGE) {
		for _, namespace := range m.GetNamespacesNames() {
			r.nodes.SendMessage(m)
			Log.Debugf("message send to %v", namespace)
		}

	} else if (client_registered != nil) && (message_type == NOTIFICATION) && (communication_type == REGISTRATION) {
		client = client_registered
		client.Write(m)
	} else if (client_registered != nil) && (message_type == NOTIFICATION) && (communication_type == SUBSCRIPTION) {
		client = client_registered
		client.Write(m)

	} else if (client_registered != nil) && (message_type == NOTIFICATION) && (communication_type == UNSUBSCRIPTION) {
		client = client_registered
		client.Write(m)

	} else if (client_registered != nil) && (message_type == NOTIFICATION) && (communication_type == INFO) {
		for _, namespace := range m.GetNamespacesNames() {
			r.nodes.SendMessage(m)
			Log.Debugf("notification send to %v", namespace)
		}

	} else if message_type == SIGNAL {
		for _, namespace := range m.GetNamespacesNames() {
			r.nodes.SendMessage(m)
			Log.Debugf("user %v subscribed to %v", client.Id.String(), namespace)

		}
	}

}

func (r *Router) routing() {

	r.active = true
	r.status <- ROUTER
	for {
		select {
		case m := <-r.input:
			r.routingCriteria(m)
			continue

		case s := <-r.signals:
			if s.signal_type == TIMEOUT {
				log.Println(s.signal_message)
			} else if s.signal_type == DISCONNECT {
				r.engine.PublishSignal(s)
				Log.Debugf(s.signal_message)
				break
			} else {
				Log.Debugf(s.signal_message)
			}
			break
		}
		break
	}
	r.active = false
	Log.Debugf("Closing Router")

}

func (r *Router) SendMessage(m IMessage) {
	if r.engine.GetStatus() {
		r.engine.PublishMessage(m)
	} else {
		Log.Errorf("message %v could not been send, engine is clossing", m.GetContent())
	}

}

func (r *Router) splitNamespace(name string) []string {
	return strings.Split(name, ".")
}

func (r *Router) Subscribe(nodes interface{}, client *Client) {
	node_list, node_list_ok := nodes.([]string)
	if node_list_ok {
		for _, n := range node_list {
			node_created := r.nodes.NewNamespace(stringToNamespaces(n))
			if node_created != nil {
				node_created.AddClient(client)
				client.node_list = append(client.node_list, n)
				Log.Debugf("User %v subscribed %v", client.Id.String(), node_created.NodeName)
			} else {
				Log.Errorf("Error subscribing to namespace")
			}

		}
	}

}

func (r *Router) Unsubscribe(nodes interface{}, client *Client) {
	node_list, node_list_ok := nodes.([]string)
	if node_list_ok {
		for _, n := range node_list {
			node_created := r.nodes.FindNamespaceNode(n)
			if node_created != nil {
				node_created.RemoveClient(client)
				Log.Debugf("User %v unsubscribed %v", client.Id.String(), node_created.NodeName)
			} else {
				Log.Errorf("Error unsubscribing to namespace")
			}

		}
	}
}

func (r *Router) NewClient(conn *websocket.Conn) *Client {
	client := NewClient(conn, r)
	r.unregistered_clients.AddClient(client)
	return client
}

func (r *Router) unregisteredTimeout() {
	ticker := time.NewTicker(60000 * time.Millisecond)
	defer ticker.Stop()
	for t := range ticker.C {
		for item := range r.unregistered_clients.clients.IterBuffered() {
			client := item.Val.(*Client)
			if client.Time.Sub(t)*(-1) > 10*time.Second {
				Log.Debugf("removing unregister user %v", client.Id.String())
				r.RemoveFromUnregisteredClients(client)
			}
		}
	}
}

func (r *Router) RemoveNode(c *Node) {
	r.RemoveNode(c)
}

func (r *Router) RemoveClient(c *Client) {
	//client_element := r.registered_clients.GetById(c.Id)
	//r.Unsubscribe(client_element.node_list, c)
	r.registered_clients.RemoveClient(c)
	c.conn.Close()
}

func (r *Router) RemoveFromUnregisteredClients(c *Client) {
	Log.Debugf("remove client %v unregistered for timeout", c.Id)
	r.unregistered_clients.RemoveClient(c)
	c.conn.Close()
}

func (r *Router) RemoveFromNodeClient(n *Node, c *Client) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	n.RemoveClient(c)
	c.conn.Close()
}
