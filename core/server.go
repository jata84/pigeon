package core

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var cwd, _ = os.Getwd()

var (
	addr = flag.String("localhost", ":8080", "http service address")
)

type Response struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}

func (r *Response) ToJson() string {
	response_json, _ := json.Marshal(r)
	return string(response_json)
}

func NewResponse(status int, description string) *Response {
	return &Response{
		Status:      status,
		Description: description,
	}
}

type Server struct {
	status        chan string
	server_router *mux.Router
}

func NewServer(status chan string) *Server {
	return &Server{
		status:        status,
		server_router: mux.NewRouter(),
	}
}

func (s *Server) server_websocket(router *Router, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Log.Error("")
		return
	}
	Log.Debug("Nueva conexión")
	client := router.NewClient(conn)
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		data_message, err := NewMessageFromJson(message)
		if err != nil {
			Log.Errorf("incorrect json format %v", err)
			//error_message := NewErrorMessage(client.GetClientChannel(),client,fmt.Sprintf("error in json format %v",string(message)))
			//router.SendMessage(error_message)
		} else {
			data_message.Comunication.Client = client
			router.SendMessage(data_message)
			Log.Debugf("message sended %v", data_message.Data)
		}
	}
	router.RemoveClient(client)
}

func (s *Server) list_clients(router *Router, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, router.registered_clients.ToJson())
}

func (s *Server) list_node_clients(router *Router, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace_name := vars["namespace_key"]
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, router.nodes.FindNamespace(stringToNamespaces(namespace_name)).clientList.ToJson())
}

func (s *Server) list_nodes(router *Router, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, router.nodes.ToJson())
}

func (s *Server) get_nodes(router *Router, w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	namespace_key := vars["namespace_key"]
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, router.nodes.FindNamespace(stringToNamespaces(namespace_key)).ToJson())
}

func (s *Server) send_message(router *Router, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message_json *Message
	err := decoder.Decode(&message_json)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, (NewResponse(http.StatusBadRequest, "error decoding json")).ToJson())
	} else {
		content_message := message_json
		if content_message.Comunication.GetCommunicationType() == MESSAGE {
			content_message.Sender = API
			router.SendMessage(message_json)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, (NewResponse(http.StatusOK, "message sended")).ToJson())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, (NewResponse(http.StatusBadRequest, "message type not allowed")).ToJson())
		}
	}

}

func (s *Server) close_client(router *Router, w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//client_id := []byte(vars["client_key"])
	//signal := NewDisconnectSignal()
}

func (s *Server) subscribe_channel(router *Router, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message_json *Message
	err := decoder.Decode(&message_json)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, (NewResponse(http.StatusBadRequest, "error decoding json")).ToJson())
	} else {
		content_message := message_json
		if content_message.Comunication.GetCommunicationType() == SUBSCRIPTION {
			message_json.Comunication.Client = router.registered_clients.GetById(content_message.GetClient().Id)
			message_json.Sender = API
			router.SendMessage(message_json)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, (NewResponse(http.StatusOK, "message sended")).ToJson())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, (NewResponse(http.StatusBadRequest, "message type not allowed")).ToJson())
		}
	}

}

func (s *Server) unsubscribe_channel(router *Router, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message_json *Message
	err := decoder.Decode(&message_json)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, (NewResponse(http.StatusBadRequest, "error decoding json")).ToJson())
	} else {
		content_message := message_json
		if content_message.Comunication.GetCommunicationType() == UNSUBSCRIPTION {
			message_json.Comunication.Client = router.registered_clients.GetById(content_message.GetClient().Id)
			message_json.Sender = API
			router.SendMessage(message_json)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, (NewResponse(http.StatusOK, "message sended")).ToJson())
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, (NewResponse(http.StatusBadRequest, "message type not allowed")).ToJson())
		}
	}

}

func (s *Server) get_clients(router *Router, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client_id := []byte(vars["client_key"])

	var uuid_param uuid.UUID

	err := uuid_param.UnmarshalText(client_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

	} else {
		client_ptr := router.registered_clients
		if client_ptr == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, client_ptr.ToJson())
		}

	}

}

func (s *Server) delete_clients(router *Router, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client_id := []byte(vars["client_key"])

	var uuid_param uuid.UUID

	err := uuid_param.UnmarshalText(client_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

	} else {
		client_ptr := router.registered_clients.GetById(uuid_param)
		if client_ptr == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			router.RemoveClient(client_ptr)
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, (NewResponse(http.StatusOK, "client deleted")).ToJson())
		}

	}

}

func (s *Server) get_clients_channels(router *Router, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	client_id := []byte(vars["client_key"])

	var uuid_param uuid.UUID

	err := uuid_param.UnmarshalText(client_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

	} else {
		client_ptr := router.registered_clients.GetById(uuid_param)
		if client_ptr == nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
			node_list := client_ptr.GetClientNode()
			json_data, _ := json.Marshal(node_list)
			fmt.Fprint(w, string(json_data))
		}

	}
}

func (s *Server) serveWS(router *Router, wg *sync.WaitGroup) {
	var ws_addr = flag.String("ws_addr", fmt.Sprintf(":%v", Configuration.Websocket.Port), "http service address")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.server_websocket(router, w, r)
	})
	s.status <- WEBSOCKET
	err_ws := http.ListenAndServe(*ws_addr, nil)
	if err_ws != nil {
		Log.Errorf("error server Websocket: %v", err_ws)
	}
	wg.Done()
}

func (s *Server) serveAPI(router *Router, wg *sync.WaitGroup) {
	var api_addr = flag.String("api_addr", fmt.Sprintf(":%v", Configuration.Http.Port), "http service address")
	//token_auth := TokenAuthenticationMiddleware{}
	//s.server_router.Use(token_auth.Middleware)
	user_agent := UserAgentMiddleware{}
	token_auth := TokenAuthenticationMiddleware{}
	s.server_router.Use(user_agent.Middleware)
	s.server_router.Use(token_auth.Middleware)

	s.server_router.HandleFunc("/api/subscribe", func(w http.ResponseWriter, r *http.Request) {
		s.subscribe_channel(router, w, r)
	})

	s.server_router.HandleFunc("/api/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		s.unsubscribe_channel(router, w, r)
	})

	s.server_router.HandleFunc("/api/message", func(w http.ResponseWriter, r *http.Request) {
		s.send_message(router, w, r)
	}).Methods(http.MethodPost)

	s.server_router.HandleFunc("/api/namespaces/{namespace_key}", func(w http.ResponseWriter, r *http.Request) {
		s.get_nodes(router, w, r)
	})

	s.server_router.HandleFunc("/api/namespaces", func(w http.ResponseWriter, r *http.Request) {
		s.list_nodes(router, w, r)
	})

	s.server_router.HandleFunc("/api/namespaces/{namespace_key}/clients", func(w http.ResponseWriter, r *http.Request) {
		s.list_node_clients(router, w, r)
	}).Methods(http.MethodGet)

	s.server_router.HandleFunc("/api/clients", func(w http.ResponseWriter, r *http.Request) {
		s.list_clients(router, w, r)
	}).Methods(http.MethodGet)

	s.server_router.HandleFunc("/api/clients/{client_key}", func(w http.ResponseWriter, r *http.Request) {
		s.get_clients(router, w, r)
	}).Methods(http.MethodGet)

	s.server_router.HandleFunc("/api/clients/{client_key}", func(w http.ResponseWriter, r *http.Request) {
		s.delete_clients(router, w, r)
	}).Methods(http.MethodDelete)

	s.server_router.HandleFunc("/api/clients/{client_key}/channels", func(w http.ResponseWriter, r *http.Request) {
		s.get_clients_channels(router, w, r)
	})

	s.status <- API
	err_api := http.ListenAndServe(*api_addr,
		handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(s.server_router))
	if err_api != nil {
		Log.Errorf("error server API: %v ", err_api)
	}
	wg.Done()
}

func (s *Server) Init(router *Router, wg *sync.WaitGroup) {
	go s.serveAPI(router, wg)
	go s.serveWS(router, wg)

}
