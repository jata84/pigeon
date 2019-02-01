package core

/*
import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"encoding/json"
	"errors"
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

type Channel struct {
	sync.Mutex
	Channel_type string `json:"type"`
	Channel_name string `json:"name"`
	Time time.Time `json:"datetime"`
	client_list  ClientList
	node *ChannelList

}

func (c *Channel) GetChannelType() string {
	return c.Channel_type
}

func (c *Channel) Count() int {
	return c.client_list.Count()
}



func (c *Channel) communicationType(split []string)(string){
	last_item := split[len(split)-1]
	if (last_item == "*"){
		return WILDCARD
	}else{
		return CHANNEL
	}
}


func (c *Channel) newNode(namespace string)(*ChannelList){
	if c.node == nil{
		c.node = NewChannelList()
		c.node.AddNode(namespace)
	}else{
		c.node.AddNode(namespace)
	}

}

func (c *Channel) wildcardChannels(namespace string) ([]*Channel){
	return nil
}

func (c *Channel) searchNamespace(split []string,namespace_type string) *Channel {
	if (len(split) == 2 && namespace_type==WILDCARD){

	}else if(len(split) == 1){
		c.wildcardChannels()

	}else{
		first_elem := split[0]
		next_split :=split[1:]

		return c.searchNamespace(next_split,namespace_type)

	}

}



func (c *Channel) GetChannelName() string {
	return c.Channel_name
}

func (c *Channel) AddClient(client *Client) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if (c.GetChannelType() == NAMESPACE){
		return errors.New("channels cannot add clients, its a namespaced root")
	}else{
		Log.Debugf("client %v added to channel [%v]", client.Id.String(), c.Channel_name)
		client.AddChannelReference(c)
		c.client_list.AddClient(client)
	}
}


func (c *Channel) RemoveClient(client *Client) error {
	//c.Mutex.Lock()
	//defer c.Mutex.Unlock()
	Log.Debugf("client %v removed from channel [%v]", client.Id.String(), c.Channel_name)
	c.client_list.RemoveClient(client)
	client.RemoveChannelReference(c)

}

func (c *Channel) SendMessage(m IMessage) error {
	c.client_list.SendMessage(m)
}

func (c *Channel) ClientExists(client *Client) bool {
	return c.client_list.ClientExists(client)
}

func (c *Channel) ConnectionExists(con *websocket.Conn) bool {
	return c.client_list.ConnectionExists(con)
}

func (c *Channel) IsValid() bool {
	return true
}

func (c *Channel) ToJson() string {
	json_data, _ := json.Marshal(c)
	return string(json_data)
}

func (c *Channel) ListClients() []*Client {
	return c.client_list.ToArray()
}

func (c *Channel) GetClientById(id uuid.UUID) *Client {
	return c.client_list.GetById(id)
}

func NewProtectedChannel(channel_name string) *Channel {
	return &Channel{
		Time: time.Now(),
		Channel_type:PROTECTED,
		Channel_name:channel_name,
		client_list:*(NewClientList()),
		node:nil,
	}
}





func NewPrivateChannel(channel_name string) *Channel {
	defer Log.Debugf("new channel private created %v",channel_name)
	return &Channel{
		Time: time.Now(),
		Channel_type:PRIVATE,
		Channel_name:channel_name,
		client_list: *(NewClientList()),
		node:nil,
	}
}
func NewPublicChannel(channel_name string) *Channel {
	defer Log.Debugf("new channel private created %v",channel_name)
	return &Channel{
		Time: time.Now(),
		Channel_type:PUBLIC,
		Channel_name:channel_name,
		client_list:*(NewClientList()),
		node:nil,
	}
}

func NewSystemChannel(channel_name string) *Channel {
	defer Log.Debugf("new channel private created %v",channel_name)
	return &Channel{
		Time: time.Now(),
		Channel_type:SYSTEM,
		Channel_name:channel_name,
		client_list:*(NewClientList()),
		node:nil,
	}
}

func NewChannels(channel_list ...Channel) []Channel {
	channel_length := len(channel_list)
	channel_array := make([]Channel, channel_length)
	for _, channel := range channel_list {
		channel_array = append(channel_array, channel)
	}
	return channel_array
}

type ChannelList struct {
	channels cmap.ConcurrentMap
	node *ChannelList
}

func NewChannelList() *ChannelList {
	return &ChannelList{
		channels: cmap.New(),
	}
}

func (c *ChannelList) RemoveChannel(channel *Channel){
	c.channels.Remove(channel.Channel_name)
}

func (c *ChannelList) AddNode(name string,node_type string){

	if (node_type == NAMESPACE){
		c.channels.Set(name,NewNameSpaceNode(name))
	}else{
		c.channels.Set(name,NewPrivateNode(name))
	}
}

func (c *ChannelList) SearchByName(name string) *Channel {
	if channel, ok := c.channels.Get(name); ok {
		return channel.(*Channel)
	} else {
		return nil
	}
}

func (c *ChannelList) GetByName(name string) *Channel {
	if channel, ok := c.channels.Get(name); ok {
		return channel.(*Channel)
	} else {
		return nil
	}
}

func (c *ChannelList) GetClientChannels(client *Client) []*Channel{
	channel_list := make([]*Channel,0)
	for item := range c.channels.IterBuffered(){
		channel := item.Val.(*Channel)
		if channel.GetClientById(client.Id) != nil{
			channel_list = append(channel_list,channel)
		}
	}
	return channel_list
}

func (c *ChannelList) GetByChannelList(channels []*Channel) []*Channel {
	channel_array := make([]*Channel, 0)
	for _, v := range channels {
		channel_exists := c.GetByName(v.Channel_name)
		if channel_exists != nil {
			channel_array = append(channel_array, channel_exists)
		}

	}
	return channel_array

}

func (c *ChannelList) GetByNameList(names []string) []*Channel {
	channel_array := make([]*Channel, 0)
	for _, n := range names {
		channel_exists := c.GetByName(n)
		if channel_exists != nil {
			channel_array = append(channel_array, channel_exists)
		}

	}
	return channel_array
}

func (c *ChannelList) ClientExists(channel *Channel, client *Client) bool {
	channel_name := channel.Channel_name
	if c.ChannelExists(channel_name) {
		channel_found := c.SearchByName(channel.Channel_name)
		return channel_found.ClientExists(client)
	} else {
		Log.Debugf("client with uuid: %s not found in channel %s, channel don't exists", client.Id.String(), channel.Channel_name)
		return false
	}
}

func (c *ChannelList) ToArray() []*Channel {
	channel_list := make([]*Channel,0)
	for item := range c.channels.IterBuffered(){
		channel := item.Val.(*Channel)
		channel_list = append(channel_list,channel)
	}
	return channel_list
}

func (c *ChannelList) ToJson() string {
	json_data, _ := json.Marshal(c.ToArray())
	return string(json_data)
}

func (c *ChannelList) AddClientToChannels(channels []*Channel, client *Client) {
	for _, ch := range channels {
		ch.AddClient(client)
		client.AddChannelReference(ch)
	}
}

func (c *ChannelList) ChannelExists(channel_name string) bool {
	if _, ok := c.channels.Get(channel_name); ok {
		return true
	} else {
		return false
	}
}

func (c *ChannelList) AddChannel(channel *Channel) (*Channel,error) {
	if channel.IsValid(){
		channel_exists := c.GetByName(channel.Channel_name)
		if channel_exists != nil{
			return channel_exists,nil
		} else {
			channel_type := channel.Channel_type
			if channel_type == PRIVATE{
				new_channel := NewPrivateChannel(channel.Channel_name)
				c.channels.Set(new_channel.Channel_name,new_channel)

				return new_channel,nil
			}else if channel_type == PUBLIC{
				new_channel := NewPublicChannel(channel.Channel_name)
				c.channels.Set(new_channel.Channel_name,new_channel)
				return new_channel,nil
			}else if channel_type == PROTECTED{
				new_channel := NewProtectedChannel(channel.Channel_name)
				c.channels.Set(new_channel.Channel_name,new_channel)
				return new_channel,nil
			}else{
				return nil,errors.New("channel_type not defined")
			}


		}
	}else{

		return nil,errors.New("channel is not valid")
	}

}
*/
