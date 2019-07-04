package core

import (
	"bytes"
	"github.com/op/go-logging"
	"os"
)

var Log = logging.MustGetLogger("RTlogger")
var Configuration Config
var Information = NewInfo()

var backend_file = logging.NewLogBackend(os.Stderr, "", 0)
var backend_console = logging.NewLogBackend(os.Stderr, "", 0)

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{color:reset} %{module} %{callpath} %{message}`,
)
var backend_console_formatter = logging.NewBackendFormatter(backend_console, format)

var backend_console_level_info = logging.AddModuleLevel(backend_console_formatter)
var backend_console_level_debug = logging.AddModuleLevel(backend_console_formatter)

var (
	buf bytes.Buffer
)
