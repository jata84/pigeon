package core

import (
	"github.com/op/go-logging"
	"sync"
)

const (
	ROUTER    = "Router"
	WEBSOCKET = "Websocket"
	API       = "Api"
	ENGINE    = "Engine"
)

var LogApp = logging.MustGetLogger("APP")

var status_list = []string{ROUTER, WEBSOCKET, ENGINE, API}

type App struct {
	status chan string
}

func NewApp() *App {
	return &App{
		status: make(chan string),
	}
}

func (a *App) Init() {
	var wg sync.WaitGroup
	wg.Add(2)
	router := NewRouter(a.status)
	router.Init()
	server := NewServer(a.status)
	server.Init(router, &wg)
	for i := 0; i < len(status_list); i++ {
		val := <-a.status
		Log.Infof("%v Ready", val)
	}
	Log.Infof("Waiting for connections...")
	wg.Wait()
}
