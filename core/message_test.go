package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistrationMessage_ToJson(t *testing.T) {
	client := NewClient(nil, nil)

	message := NewRegistrationMessage(client)
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	assert.Contains(t, string(message_json), REGISTRATION)

}
func TestSimpleSubscriptionMessage_ToJson(t *testing.T) {
	namespace := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test",
			NodeType: PUBLIC,
		},
	}
	client := NewClient(nil, nil)

	message := NewSubscriptionMessage(namespace, client)
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	//fmt.Println(string(message_json))
	assert.Contains(t, string(message_json), "namespace.public.test")
	assert.Contains(t, string(message_json), SUBSCRIPTION)
}
func TestMultipleSubscriptionMessage_ToJson(t *testing.T) {
	namespace1 := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test",
			NodeType: PUBLIC,
		},
	}

	namespace2 := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test2",
			NodeType: PUBLIC,
		},
	}

	namespace_list := [][]*Node{namespace1, namespace2}
	client := NewClient(nil, nil)

	message := NewSubscriptionMessage(namespace_list, client)
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	assert.Contains(t, string(message_json), "namespace.public.test")
	assert.Contains(t, string(message_json), "namespace.public.test2")
	assert.Contains(t, string(message_json), SUBSCRIPTION)

}
func TestSimpleContentMessage_ToJson(t *testing.T) {
	namespace := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test",
			NodeType: PUBLIC,
		},
	}
	client := NewClient(nil, nil)

	message := NewContentMessage(namespace, client, "test", "")
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	assert.Contains(t, string(message_json), "namespace.public.test")
	assert.Contains(t, string(message_json), CONTENT)
}
func TestMultipleContentMessage_ToJson(t *testing.T) {
	namespace1 := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test",
			NodeType: PUBLIC,
		},
	}
	namespace2 := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test2",
			NodeType: PUBLIC,
		},
	}

	namespace_list := [][]*Node{namespace1, namespace2}
	client := NewClient(nil, nil)

	message := NewContentMessage(namespace_list, client, "message test", "")
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	assert.Contains(t, string(message_json), "namespace.public.test")
	assert.Contains(t, string(message_json), "namespace.public.test2")
	assert.Contains(t, string(message_json), CONTENT)
}
func TestMultipleNotificationMessage_ToJson(t *testing.T) {
	namespace1 := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test",
			NodeType: PUBLIC,
		},
	}
	namespace2 := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test2",
			NodeType: PUBLIC,
		},
	}

	namespace_list := [][]*Node{namespace1, namespace2}
	client := NewClient(nil, nil)

	message := NewNotificationInfo(namespace_list, client, "test_notification", "")
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}

	assert.Contains(t, string(message_json), "namespace.public.test")
	assert.Contains(t, string(message_json), "namespace.public.test2")
	assert.Contains(t, string(message_json), NOTIFICATION)
}
func TestSimpleNotificationMessage_ToJson(t *testing.T) {
	namespace := []*Node{
		&Node{
			NodeName: "namespace",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "public",
			NodeType: NAMESPACE,
		},
		&Node{
			NodeName: "test",
			NodeType: PUBLIC,
		},
	}
	client := NewClient(nil, nil)

	message := NewNotificationInfo(namespace, client, "test_notification", "")
	message_json, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	assert.Contains(t, string(message_json), "namespace.public.test")
	assert.Contains(t, string(message_json), NOTIFICATION)
}
func TestRegistrationMessage_FromJson(t *testing.T) {
	json := `{
 "comunication": {
  "type": "REGISTRATION",
  "namespaces": [],
  "client": {
   "uuid": "cdc9eda1-808d-4235-9974-6b8fd2cde84d",
   "datetime": "2018-12-09T18:05:16.263640434+01:00"
  }
 },
 "type": "MESSAGE",
 "content": "",
 "timestamp": "2018-12-09T18:05:16.263650532+01:00"
}`
	message, err := NewMessageFromJson(json)
	if err != nil {
		t.Errorf("Error encoding Json")

	}
	json_message, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}

	assert.Contains(t, string(json_message), REGISTRATION)

}
func TestSimpleSubscriptionMessage_FromJson(t *testing.T) {
	json := `{
 "comunication": {
  "type": "SUBSCRIPTION",
  "namespaces": ["namespace.public.test"],
  "client": {
   "uuid": "cdc9eda1-808d-4235-9974-6b8fd2cde84d",
   "datetime": "2018-12-09T18:05:16.263640434+01:00"
  }
 },
 "type": "MESSAGE",
 "content": "",
 "timestamp": "2018-12-09T18:05:16.263650532+01:00"
}`
	message, err := NewMessageFromJson(json)
	if err != nil {
		t.Errorf("Error encoding Json")

	}
	json_message, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}

	assert.Contains(t, string(json_message), SUBSCRIPTION)
	assert.Contains(t, string(json_message), "namespace.public.test")
}
func TestMultipleSubscriptionMessage_FromJson(t *testing.T) {
	json := `{
 "comunication": {
  "type": "SUBSCRIPTION",
  "namespaces": ["namespace.public.test","namespace.public.test2"],
  "client": {
   "uuid": "cdc9eda1-808d-4235-9974-6b8fd2cde84d",
   "datetime": "2018-12-09T18:05:16.263640434+01:00"
  }
 },
 "type": "MESSAGE",
 "content": "",
 "timestamp": "2018-12-09T18:05:16.263650532+01:00"
}`
	message, err := NewMessageFromJson(json)
	if err != nil {
		t.Errorf("Error encoding Json")
	}
	json_message, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}

	assert.Contains(t, string(json_message), SUBSCRIPTION)
	assert.Contains(t, string(json_message), "namespace.public.test")
	assert.Contains(t, string(json_message), "namespace.public.test2")

}
func TestSimpleContentMessage_FromJson(t *testing.T) {
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
	}
	json_message, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}

	assert.Contains(t, string(json_message), CONTENT)
	assert.Contains(t, string(json_message), "namespace.public.test")
	assert.Contains(t, string(json_message), "namespace.public.test2")
}
func TestMultipleContentMessage_FromJson(t *testing.T) {
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
	}
	json_message, err := message.ToJson()
	if err != nil {
		t.Errorf("Error encoding Json")
	}

	assert.Contains(t, string(json_message), CONTENT)
	assert.Contains(t, string(json_message), "namespace.public.test")
	assert.Contains(t, string(json_message), "namespace.public.test2")
}
