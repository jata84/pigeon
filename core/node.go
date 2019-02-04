package core

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"strings"
	"time"

	"github.com/orcaman/concurrent-map"

	"sync"
)

const (
	PRIVATE   = "PRIVATE"
	NAMESPACE = "NAMESPACE"
	PROTECTED = "PROTECTED"
	SYSTEM    = "SYSTEM"
)

type NodeTree struct {
	root *Node
}

func NewNodeTree() *NodeTree {
	return &NodeTree{
		root: NewNameSpaceNode(ROOT, nil),
	}
}

type ElementList struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func toList(n *Node, ar []*ElementList, namespace string) []*ElementList {
	if n.NodeList != nil {
		if n.NodeName != ROOT {
			if namespace != "" {
				namespace = namespace + "." + n.NodeName
			} else {
				namespace = n.NodeName
			}
		}
		for item := range n.NodeList.nodes.IterBuffered() {
			node := item.Val
			ar = toList(node.(*Node), ar, namespace)
		}
		return ar

	} else {
		if n.NodeType != NAMESPACE {
			if namespace != "" {
				namespace = namespace + "." + n.NodeName
			} else {
				namespace = n.NodeName
			}
			ar = append(ar, &ElementList{
				Name: namespace,
				Type: n.NodeType,
			})
		}

		return ar

	}
}

func (nt *NodeTree) ToJson() string {

	var list = []*ElementList{}
	list = toList(nt.root, list, "")
	json_string, err := json.MarshalIndent(list, "", "\t")
	if err != nil {
		return ""
	} else {
		return string(json_string)
	}
}

func (nt *NodeTree) validateNamespace(name string) {

}

func NewNameSpaceNode(namespace_node string, previousNode *NodeList) *Node {
	return &Node{
		Time:         time.Now(),
		NodeType:     NAMESPACE,
		NodeName:     namespace_node,
		clientList:   nil,
		NodeList:     nil,
		PreviousNode: previousNode,
	}
}

func NewPrivateNode(namespace_node string, previousNode *NodeList) *Node {
	defer Log.Debugf("new private node created %v", namespace_node)
	return &Node{
		Time:         time.Now(),
		NodeType:     PRIVATE,
		NodeName:     namespace_node,
		clientList:   nil,
		NodeList:     nil,
		PreviousNode: previousNode,
	}
}

func NewPublicNode(namespace_node string, previousNode *NodeList) *Node {
	defer Log.Debugf("new public node created %v", namespace_node)
	return &Node{
		Time:         time.Now(),
		NodeType:     PUBLIC,
		NodeName:     namespace_node,
		clientList:   nil,
		NodeList:     nil,
		PreviousNode: previousNode,
	}
}

func (nt *NodeTree) communicationType(split []*Node) string {
	last_item := split[len(split)-1]

	if last_item.NodeName == "*" {
		return WILDCARD
	} else {
		return NAMESPACE
	}
}

func (nt *NodeTree) splitNamespace(name string) []string {
	return strings.Split(name, ".")
}

func (nt *NodeTree) getLastNamespace(splitted_namespace []*Node) *Node {

	if len(splitted_namespace) > 0 {
		splitted_namespace = splitted_namespace[:len(splitted_namespace)-1]
	} else {
		return nil
	}

	initial_node := nt.root
	return initial_node.GetNode(splitted_namespace)
}

func (nt *NodeTree) FindNamespace(namespace []*Node) *Node {
	return nt.findNamespace(namespace)
}
func (nt *NodeTree) findNamespace(splitted_namespace []*Node) *Node {
	initial_node := nt.root
	return initial_node.GetNode(splitted_namespace)
}

func (nt *NodeTree) findNamespaceNode(splitted_namespace []*Node) *Node {
	initial_node := nt.root
	return initial_node.GetNode(splitted_namespace)
}

func (nt *NodeTree) FindNamespaceNode(namespace_string string) *Node {
	splitted_namespace := stringToNamespaces(namespace_string)
	return nt.findNamespaceNode(splitted_namespace)
}

func (nt *NodeTree) getPreviousNamespace(splitted_namespace []*Node) *Node {
	if len(splitted_namespace) == 1 {
		return nt.root
	} else {
		var namespace_value = splitted_namespace[:len(splitted_namespace)-1]
		node_found := nt.findNamespaceNode(namespace_value)
		return node_found
	}
}

func (nt *NodeTree) getPreviousNamespaceNode(splitted_namespace []*Node) *Node {
	if len(splitted_namespace) == 1 {
		return nt.root
	} else {
		var namespace_value = splitted_namespace[:len(splitted_namespace)-1]
		return nt.findNamespace(namespace_value)
	}
}

func (nt *NodeTree) NewNamespace(splitted_namespace []*Node) *Node {
	var root *Node
	var err error
	root = nt.root
	//namespace_length := len(splitted_namespace)
	var node_response *Node
	for _, node_element := range splitted_namespace {
		//last_element := (key == namespace_length-1)
		if root.NodeList != nil {
			next_node := root.NodeList.GetByName(node_element.NodeName)
			if next_node != nil {
				root = next_node
				if root.NodeType != node_element.NodeType {
					break
				} else {
					if root.NodeType == PRIVATE {
						return root
					}
				}
			} else {
				node_response, err = root.NewNode(node_element)
				if err != nil {
					break
				} else {
					root = node_response
				}

			}

		} else {
			node_response, err = root.NewNode(node_element)

			if err != nil {
				break
			} else {
				root = node_response
			}
		}

	}
	return node_response

}

func (nt *NodeTree) NewNode(splitted_namespace []*Node, namespace_type string) (*Node, error) {
	var node_element *Node
	node_element = nt.getPreviousNamespace(splitted_namespace)
	if node_element == nil {
		return nil, errors.New("Previous Namespace not found")

	} else if node_element.NodeType == PRIVATE {
		return nil, errors.New("Previous Namespace is Private")
	}
	return node_element.NewNode(splitted_namespace[len(splitted_namespace)-1])

}

func (nt *NodeTree) NewNodeList(node_list []*Node, namespace_type string) (*Node, error) {
	var node_element *Node
	node_element = nt.getPreviousNamespace(node_list)
	if node_element == nil {
		return nil, errors.New("Previous Namespace not found")
	} else if node_element.NodeType == PRIVATE {
		return nil, errors.New("Previous Namespace is Private")
	}
	return node_element.NewNode(node_list[len(node_list)-1])
}

func (nt *NodeTree) SendMessage(m IMessage) []error {
	var error_list []error
	var node *Node
	for _, namespace := range m.GetNamespacesNames() {
		node_list := stringToNamespaces(namespace)
		if node_list != nil {
			communication_type := nt.communicationType(node_list)
			if communication_type == WILDCARD {
				node = nt.getLastNamespace(stringToNamespaces(namespace))
			} else {
				node = nt.FindNamespaceNode(namespace)
			}

			if node != nil {
				if communication_type == WILDCARD {
					for item := range node.NodeList.nodes.IterBuffered() {
						node := item.Val
						er := node.(*Node).SendMessage(m)
						error_list = append(error_list, er)
					}

				} else {
					err := node.SendMessage(m)
					error_list = append(error_list, err)
				}

			} else {
				err := errors.New("node not found")
				error_list = append(error_list, err)
			}
		}

	}
	if len(error_list) == 0 {
		return nil
	}
	return error_list
}

type Node struct {
	sync.Mutex
	NodeType     string    `json:"type"`
	NodeName     string    `json:"name"`
	Time         time.Time `json:"datetime"`
	clientList   *ClientList
	NodeList     *NodeList `json:"nodes"`
	PreviousNode *NodeList `json:"-"`
}

func (n *Node) ToJson() string {
	json_data, _ := json.MarshalIndent(n, "", "\t")
	return string(json_data)
}

func (n *Node) communicationType(split []string) string {
	last_item := split[len(split)-1]
	if last_item == "*" {
		return WILDCARD
	} else {
		return NAMESPACE
	}
}

func (n *Node) newNode(namespace *Node, node_type string) (*Node, error) {
	if n.NodeList == nil {
		n.NodeList = NewNodeList()
		return n.NodeList.AddNode(namespace, node_type)
	} else {
		return n.NodeList.AddNode(namespace, node_type)
	}

}

func (n *Node) NewNode(namespace *Node) (*Node, error) {
	return n.newNode(namespace, namespace.NodeType)
}

func (n *Node) AddClient(client *Client) {
	if n.clientList == nil {
		n.clientList = NewClientList()
		n.clientList.AddClient(client)
	} else {
		n.clientList.AddClient(client)
	}

}

func (n *Node) RemoveClient(client *Client) {
	n.clientList.RemoveClient(client)
	count := n.clientList.Count()
	if count == 0 {
		n.PreviousNode.nodes.Remove(n.NodeName)
	}
}

func (n *Node) GetNode(namespace_list []*Node) *Node {
	if n.NodeList == nil {
		return nil
	}
	node := namespace_list[0]
	node_result := n.NodeList.GetByName(node.NodeName)
	if node_result != nil && len(namespace_list) > 1 {
		namespace_list = namespace_list[1:]
		return node_result.GetNode(namespace_list)
	} else {
		return node_result
	}
}

func (n *Node) wildcardChannels(namespace string) []*Node {
	return nil
}

func (n *Node) GetNodeType() string {
	return n.NodeType
}

func (n *Node) Count() int {
	return n.clientList.Count()
}

func (n *Node) GetNodeName() string {
	return n.NodeName
}

func (n *Node) SendMessage(m IMessage) error {
	n.clientList.SendMessage(m)
	return nil
}

func (n *Node) ClientExists(client *Client) bool {
	return n.clientList.ClientExists(client)
}

func (n *Node) ConnectionExists(con *websocket.Conn) bool {
	return n.clientList.ConnectionExists(con)
}

func (n *Node) IsValid() bool {
	return true
}

func (n *Node) ListClients() []*Client {
	return n.clientList.ToArray()
}

func (n *Node) GetClientById(id uuid.UUID) *Client {
	return n.clientList.GetById(id)
}

type NodeList struct {
	nodes cmap.ConcurrentMap
}

func (n *NodeList) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.nodes)
}

func NewNodeList() *NodeList {
	return &NodeList{
		nodes: cmap.New(),
	}
}

func (n *NodeList) AddNode(node *Node, node_type string) (*Node, error) {
	if node_type == NAMESPACE {
		n.nodes.Set(node.NodeName, NewNameSpaceNode(node.NodeName, n))
	} else if node_type == PRIVATE {
		n.nodes.Set(node.NodeName, NewPrivateNode(node.NodeName, n))

	} else {
		n.nodes.Set(node.NodeName, NewPublicNode(node.NodeName, n))
	}
	node_element, result := n.nodes.Get(node.NodeName)
	if result {
		return node_element.(*Node), nil
	} else {
		return nil, errors.New("namespace not created")
	}
}

func (c *NodeList) GetByName(name string) *Node {
	if node, ok := c.nodes.Get(name); ok {
		return node.(*Node)
	} else {
		return nil
	}
}
