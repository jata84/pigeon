package core

import (
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type Info struct {
	Hosts map[string]time.Time
}

func NewInfo() Info {
	return Info{
		Hosts: make(map[string]time.Time),
	}
}

func (i *Info) SetHost(s string) {
	i.Hosts[s] = time.Now()
}

func (i *Info) UnsetHost(s string) {
	i.Hosts[s] = time.Now()
}

type Config struct {
	Websocket struct {
		Port   int
		Scheme string
	}
	Http struct {
		Port   int
		Scheme string
	}

	Engine struct {
		Engine_type    string
		Engine_address string
		Engine_port    string
	}

	Authorization struct {
		Authorization_type  string
		Authorization_token string
	}
	Main struct {
		NodeBeat    time.Duration //Time host is signaling Node activity
		NodeTimeout time.Duration //Time host is disabled when a node stay NodeTimeout seconds without transmit information
		Hostname    string
		Gracetime   int
		Logpath     string
		LogLevel    string
	}
}

func LoadConfig() {

	b, err_file := ioutil.ReadFile("config.yaml")
	if err_file != nil {
		log.Fatalf("error: %v", err_file)
	}
	err_config := yaml.Unmarshal(b, &Configuration)
	if err_config != nil {
		log.Fatalf("error: %v", err_config)
	}

	backend_console_level_info.SetLevel(logging.INFO, "")
	backend_console_level_debug.SetLevel(logging.DEBUG, "")

	//backend_console_level.SetLevel(logging.ERROR, "")

	//CONFIG

	Log.SetBackend(backend_console_level_info)
	Log.SetBackend(backend_console_level_debug)
	logging.SetBackend(backend_console, backend_console_formatter)

}
