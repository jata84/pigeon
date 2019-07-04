package core

import (
	"errors"
	"github.com/op/go-logging"
	"sync"
	"time"
)

const (
	ROUTER        = "Router"
	WEBSOCKET     = "Websocket"
	API           = "Api"
	ENGINE        = "Engine"
	ADMIN_CHANNEL = "Admin"
)

var LogApp = logging.MustGetLogger("APP")

var status_list = []string{ROUTER, WEBSOCKET, ENGINE, API}

type App struct {
	status chan string
	router *Router
}

func NewApp() *App {
	return &App{
		status: make(chan string),
	}
}

func (a *App) Beat() {
	//router.SendSignal(NewNodeSignal("test"))
	ticker := time.NewTicker(Configuration.Main.NodeBeat * time.Second)
	go func() {

		for t := range ticker.C {
			a.router.SendSignal(NewNodeSignal(Configuration.Main.Hostname))
			for k := range Information.Hosts {
				timestamp := Information.Hosts[k]
				if timestamp.Sub(t)*(-1) > Configuration.Main.NodeTimeout*time.Second {
					delete(Information.Hosts, k)

				}

			}
		}
	}()
}

func (a *App) Init() error {
	var wg sync.WaitGroup
	var err error
	wg.Add(2)
	a.router, err = NewRouter(a.status)
	if err != nil {
		return errors.New("Error during router initialization ")
		wg.Done()
	}
	a.router.Init()
	server := NewServer(a.status)
	server.Init(a.router, &wg)
	for i := 0; i < len(status_list); i++ {
		val := <-a.status
		Log.Infof("%v Ready", val)
	}
	a.Beat()
	Log.Infof("Waiting for connections...")

	wg.Wait()

	return nil
}
