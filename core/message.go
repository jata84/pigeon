package core

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"
)

//COMUNICATION TYPE
const (
	MESSAGE      = "MESSAGE"
	NOTIFICATION = "NOTIFICATION"
	SIGNAL       = "SIGNAL"
	BLANK        = ""
)

//MESSAGE TYPE
const (
	CONTENT        = "CONTENT"
	SUBSCRIPTION   = "SUBSCRIPTION"
	UNSUBSCRIPTION = "UNSUBSCRIPTION"
	REGISTRATION   = "REGISTRATION"
)

//SIGNAL TYPE
const (
	INFO       = "INFO"
	ERROR      = "ERROR"
	WARNING    = "WARNING"
	DISCONNECT = "DISCONNECT"
)

//SIGNAL SUBTYPE
const (
	NODE_BEAT = "NODE_BEAT"
)

type IMessage interface {
	GetContent() interface{}
	GetSender() string
	GetClient() *Client
	GetNamespacesNames() []string
	GetNamespacesNodes() [][]*Node
	GetMessageType() string
	GetCommunicationType() string
	Validate() bool
	ToJson() ([]byte, error)
}

type Comunication struct {
	ComunicationType string    `json:"type"`
	Namespaces       [][]*Node `json:"namespaces" validate:"required"`
	Client           *Client   `json:"client"`
}

func (c *Comunication) NodeArrayToString() []string {
	namespace_list := []string{}

	for _, namespace := range c.Namespaces {
		namespace_len := len(namespace)
		var buffer bytes.Buffer
		for namespace_key, namespace_node := range namespace {
			buffer.WriteString(namespace_node.NodeName)
			if (namespace_len - 1) != namespace_key {
				buffer.WriteString(".")
			}
		}
		namespace_list = append(namespace_list, buffer.String())
	}
	return namespace_list
}

func (c *Comunication) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ComunicationType string   `json:"type"`
		Namespaces       []string `json:"namespaces"`
		Client           *Client  `json:"client"`
	}{
		ComunicationType: c.ComunicationType,
		Namespaces:       c.NodeArrayToString(),
		Client:           c.Client,
	})
}

func (c *Comunication) GetCommunicationType() string {
	return c.ComunicationType
}

func (c *Comunication) GetClient() *Client {
	return c.Client

}

func (c *Comunication) GetNamespaces() [][]*Node {
	return c.Namespaces
}

type Message struct {
	IMessage     `json:"-"`
	Sender       string       `json:"-"`
	Comunication Comunication `json:"comunication"`
	MessageType  string       `json:"type"`
	Data         interface{}  `json:"content"`
	Time         time.Time    `json:"timestamp"`
}

type Notification struct {
	IMessage     `json:"-"`
	Description  string       `json:"description"`
	Comunication Comunication `json:"comunication"`
	MessageType  string       `json:"type"`
	Data         interface{}  `json:"content"`
	Time         time.Time    `json:"timestamp"`
}

func interfaceToNamespaces(namespaces interface{}) [][]*Node {
	message_namespaces, slice_ok := namespaces.([][]*Node)
	message_namespace, namespace_ok := namespaces.([]*Node)

	if slice_ok {
		return message_namespaces
	} else if namespace_ok {
		return [][]*Node{message_namespace}
	} else {
		return nil
	}
}

func stringToNamespaces(namespaces string) []*Node {
	var node_type string
	if namespaces == "" {
		return nil
	}
	public_namespace := strings.HasPrefix(namespaces, "+")
	private_namespace := strings.HasPrefix(namespaces, "-")
	if public_namespace || private_namespace {
		namespaces = namespaces[1:]
	} else {
		private_namespace = true
	}

	splitted_namespace := strings.Split(namespaces, ".")
	namespace_length := len(splitted_namespace)
	var node_slice []*Node
	for k, element := range splitted_namespace {
		if k == (namespace_length - 1) {
			if node_type = PRIVATE; public_namespace {
				node_type = PUBLIC
			}
			node_slice = append(node_slice, &Node{
				NodeType: node_type,
				NodeName: element,
			})
		} else {
			node_slice = append(node_slice, &Node{
				NodeType: NAMESPACE,
				NodeName: element,
			})
		}

	}
	return node_slice
}

func sliceToNamespaces(namespaces []string) [][]*Node {
	var namespace_slice [][]*Node
	for _, namespace := range namespaces {
		splitted_namespace := strings.Split(namespace, ".")
		namespace_length := len(splitted_namespace)
		var node_slice []*Node
		for k, element := range splitted_namespace {
			if k == (namespace_length - 1) {
				node_slice = append(node_slice, &Node{
					NodeType: PRIVATE,
					NodeName: element,
				})
			} else {
				node_slice = append(node_slice, &Node{
					NodeType: NAMESPACE,
					NodeName: element,
				})
			}
		}
		namespace_slice = append(namespace_slice, node_slice)
	}
	return namespace_slice
}

func NewComumnication(communication_type string, namespaces interface{}, c *Client) *Comunication {
	return &Comunication{
		ComunicationType: communication_type,
		Namespaces:       interfaceToNamespaces(namespaces),
		Client:           c,
	}
}

func NewComumnicationNamespaceString(communication_type string, namespaces []string, c *Client) *Comunication {
	return &Comunication{
		ComunicationType: communication_type,
		Namespaces:       sliceToNamespaces(namespaces),
		Client:           c,
	}
}

func NewSubscriptionMessage(namespaces interface{}, c *Client) *Message {
	return &Message{
		Time:         time.Now(),
		Comunication: *NewComumnication(SUBSCRIPTION, namespaces, c),
		MessageType:  MESSAGE,
	}
}

func NewUnsubscriptionMessage(namespaces interface{}, c *Client) *Message {
	return &Message{
		Time:         time.Now(),
		Comunication: *NewComumnication(UNSUBSCRIPTION, namespaces, c),
		MessageType:  MESSAGE,
	}
}

func NewRegistrationMessage(c *Client) *Message {
	return &Message{
		Time:         time.Now(),
		Comunication: *NewComumnication(REGISTRATION, nil, c),
		MessageType:  MESSAGE,
	}
}

func NewContentMessage(namespaces interface{}, c *Client, message string, sender string) *Message {
	return &Message{
		Sender:       sender,
		Time:         time.Now(),
		Comunication: *NewComumnication(MESSAGE, namespaces, c),
		MessageType:  MESSAGE,
		Data:         message,
	}
}

func NewNotificationInfo(namespaces interface{}, c *Client, message string, sender string) *Message {
	return &Message{
		Sender:       sender,
		Time:         time.Now(),
		Comunication: *NewComumnication(INFO, namespaces, c),
		MessageType:  NOTIFICATION,
		Data:         message,
	}
}

func NewNotificationRegistration(c *Client, sender string) *Message {
	return &Message{
		Sender:       sender,
		Time:         time.Now(),
		Comunication: *NewComumnication(REGISTRATION, nil, c),
		MessageType:  NOTIFICATION,
		Data:         "",
	}
}

func NewNotificationSubscription(c *Client, sender string, namespace []string) *Message {
	return &Message{
		Sender:       sender,
		Time:         time.Now(),
		Comunication: *NewComumnicationNamespaceString(SUBSCRIPTION, namespace, c),
		MessageType:  NOTIFICATION,
		Data:         "",
	}
}

func NewNotificationUnsubscription(c *Client, sender string, namespace []string) *Message {
	return &Message{
		Sender:       sender,
		Time:         time.Now(),
		Comunication: *NewComumnicationNamespaceString(UNSUBSCRIPTION, namespace, c),
		MessageType:  NOTIFICATION,
		Data:         "",
	}
}

func NewNotificationError(namespaces interface{}, c *Client, message string) *Notification {
	return &Notification{
		Time:         time.Now(),
		Comunication: *NewComumnication(ERROR, namespaces, c),
		MessageType:  NOTIFICATION,
		Data:         message,
	}
}

func NewMessageFromJson(data interface{}) (*Message, error) {
	json_data_string, _ := data.(string)
	json_data_byte, byte_ok := data.([]byte)
	if !byte_ok {
		json_data_byte = []byte(json_data_string)
	}
	var message *Message

	if err := json.Unmarshal(json_data_byte, &message); err != nil {
		return nil, err
	} else {
		return message, nil
	}
}

func NewSignalFromJson(data interface{}) (*Signal, error) {
	json_data_string, _ := data.(string)
	json_data_byte, byte_ok := data.([]byte)
	if !byte_ok {
		json_data_byte = []byte(json_data_string)
	}
	var message *Signal

	if err := json.Unmarshal(json_data_byte, &message); err != nil {
		return nil, err
	} else {
		return message, nil
	}
}

func (m *Message) GetMessageType() string {
	return m.MessageType
}

func (m *Message) GetCommunicationType() string {
	return m.GetCommunication().ComunicationType
}

func (m *Message) GetCommunication() *Comunication {
	return &m.Comunication
}

func (m *Message) GetContent() interface{} {
	return m.Data
}

func (m *Message) GetClient() *Client {
	return m.GetCommunication().GetClient()
}

func (m *Message) GetChannels() [][]*Node {
	return m.GetCommunication().GetNamespaces()
}

func (m *Message) GetNamespaces() [][]*Node {
	return m.GetCommunication().GetNamespaces()
}

func (m *Message) GetSender() string {
	return m.Sender
}

func (m *Message) GetNamespacesNames() []string {
	namespaces_list := []string{}
	var buffer bytes.Buffer
	for _, namespace := range m.GetCommunication().GetNamespaces() {
		namespace_len := len(namespace)
		for namespace_key, namespace_node := range namespace {
			buffer.WriteString(namespace_node.NodeName)
			if (namespace_len - 1) != namespace_key {
				buffer.WriteString(".")
			}
		}
		namespaces_list = append(namespaces_list, buffer.String())
	}
	return namespaces_list
}

func (m *Message) GetNamespacesNodes() [][]*Node {
	return m.GetCommunication().GetNamespaces()
}

func (m *Message) Validate() bool {
	return true
}

func (m *Message) ToJson() ([]byte, error) {
	json, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return nil, err
	} else {
		return json, nil
	}
}

func (m *Notification) ToJson() ([]byte, error) {
	json, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return nil, err
	} else {
		return json, nil
	}
}

func (c *Comunication) UnmarshalJSON(b []byte) error {

	comunication := struct {
		ComunicationType string   `json:"type"`
		Namespaces       []string `json:"namespaces"`
		Client           *Client  `json:"client"`
	}{}
	err := json.Unmarshal(b, &comunication)
	c.Client = comunication.Client
	c.Namespaces = sliceToNamespaces(comunication.Namespaces)
	c.ComunicationType = comunication.ComunicationType

	if err != nil {
		return err
	}

	return nil
}
