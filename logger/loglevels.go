package logger

import (
	"github.com/ttacon/chalk"
)

// LogLevel levels for the logger
type LogLevel int

// Constants
const (
	ERROR LogLevel = iota
	WARNING
	PLUGIN
	BOOT
	INFO
	VERBOSE
	DEBUG
	MAGENTA
)

var nicenames = map[LogLevel]string{
	ERROR:   "ERROR",
	WARNING: "WARNING",
	PLUGIN:  "PLUGIN",
	BOOT:    "BOOT",
	INFO:    "INFO",
	VERBOSE: "VERBOSE",
	DEBUG:   "DEBUG",
	MAGENTA: "MAGENTA",
}

var colors = map[LogLevel]chalk.Color{
	ERROR:   chalk.Red,
	WARNING: chalk.Yellow,
	PLUGIN:  chalk.Cyan,
	BOOT:    chalk.Blue,
	INFO:    chalk.Green,
	VERBOSE: chalk.White,
	DEBUG:   chalk.ResetColor,
	MAGENTA: chalk.Magenta,
}
