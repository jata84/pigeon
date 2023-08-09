package main

import (
	"flag"
	"pidgeon/core"
)

var GitCommit string

const Version = "0.1.0"

var VersionPrerelease = ""
var BuildDate = ""

func main() {
	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		core.Log.Info("Git Commit:", GitCommit)
		core.Log.Info("Version:", Version)
		if VersionPrerelease != "" {
			core.Log.Info("Version PreRelease:", VersionPrerelease)
		}
	}
	core.LoadConfig()

	core.Log.Info("Init Application")
	app := core.NewApp()
	app.Init()
	core.Log.Info("Close Application")

}
