package main

import (
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"pidgeon/core"
	"strings"
	"time"
)

var App = grumble.New(&grumble.Config{
	Name:                  "pidgeon",
	Description:           "Pidgeon Command Line Tool",
	HistoryFile:           "/tmp/pidgeon.hist",
	Prompt:                "pidgeon » ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,

	Flags: func(f *grumble.Flags) {
		f.String("d", "directory", "DEFAULT", "set an alternative root directory path")
		f.Bool("v", "verbose", false, "enable verbose mode")
	},
})

func init() {

	ws := core.NewWS("ws://localhost:4444")
	ws.Connect()
	ws.Auth()

	App.AddCommand(&grumble.Command{
		Name:      "daemon",
		Help:      "run the daemon",
		Aliases:   []string{"run"},
		Usage:     "daemon [OPTIONS]",
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Run: func(c *grumble.Context) error {
			fmt.Println("timeout:", c.Flags.Duration("timeout"))
			fmt.Println("directory:", c.Flags.String("directory"))
			fmt.Println("verbose:", c.Flags.Bool("verbose"))

			// Handle args.
			fmt.Println("args:")
			fmt.Println(strings.Join(c.Args, "\n"))

			return nil
		},
	})

	namespaceCommand := &grumble.Command{
		Name:     "namespaces",
		Help:     "namespaces commands",
		LongHelp: "Namespaces commands",
	}

	App.AddCommand(namespaceCommand)

	namespaceCommand.AddCommand(&grumble.Command{
		Name:      "list",
		Help:      "list of conncetions",
		AllowArgs: false,
		Run: func(c *grumble.Context) error {
			var connection_list []string
			ret := ws.GetConnections()
			for _, elem := range ret {
				connection_list = append(connection_list, elem.Uuid.String())
				fmt.Println(elem.Uuid.String())
			}
			//fmt.Println(connection_list)

			return nil
		},
	})

	connectionsCommand := &grumble.Command{
		Name:     "connection",
		Help:     "Connection commands",
		LongHelp: "Connections commands",
	}

	App.AddCommand(connectionsCommand)

	connectionsCommand.AddCommand(&grumble.Command{
		Name: "list",
		Help: "list of conncetions",
		Completer: func(prefix string, args []string) []string {
			return []string{
				"test",
				"test2",
			}
		},
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", time.Second, "timeout duration")
		},
		Run: func(c *grumble.Context) error {
			var connection_list []string
			ret := ws.GetConnections()
			for _, elem := range ret {
				connection_list = append(connection_list, elem.Uuid.String())
				fmt.Println(elem.Uuid.String())
			}
			//fmt.Println(connection_list)

			return nil
		},
	})

	connectionsCommand.AddCommand(&grumble.Command{
		Name: "kill",
		Help: "kill connection",
		Completer: func(prefix string, args []string) []string {
			return []string{
				"test",
				"test2",
			}
		},
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.String("u", "uuid", "", "uuid for the string")
			f.String("a", "all", "", "kill all the connections")
		},
		Run: func(c *grumble.Context) error {

			uuid_local := c.Flags.String("uuid")
			if uuid_local != "" {
				if ws.KillConnection(uuid_local) {
					fmt.Print("DELETED")
				} else {
					fmt.Print("NOT FOUND")
				}
			}

			//fmt.Println(connection_list)

			return nil
		},
	})

}

func main() {
	grumble.Main(App)
}

/*
func main(){
	//var username string
	//var password string
	shell := ishell.New()

	// display welcome info.
	shell.Println("Sample Interactive Shell")
	ws := core.NewWS("ws://localhost:4444")
	ws.Connect()
	ws.Auth()

	// register a function for "greet" command.
	shell.AddCmd(&ishell.Cmd{
		Name: "greet",
		Help: "greet user",
		Func: func(c *ishell.Context) {
			c.Println("Hello", strings.Join(c.Args, " "))
		},
	})


	shell.AddCmd(&ishell.Cmd{
		Name: "connections",
		Help: "simulate a login",
		Func: func(c *ishell.Context) {
			var connection_list []string
			ret := ws.GetConnections()
			for _,elem := range ret{
				connection_list = append(connection_list,elem.Uuid.String())
			}
			choice := c.MultiChoice(connection_list, "Listado de Conexiones")
			if choice == 1 {
				c.Println("You got it!")
			} else {
				c.Println("Sorry, you're wrong.")
			}



		},
	})


	// run shell
	shell.Run()
}

*/
