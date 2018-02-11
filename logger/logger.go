package logger

import (
	"fmt"
	"time"
)

var (
	DebugMode = false
)

// L log the message
func (c LogLevel) L(msg string) {
	if c == DEBUG && DebugMode == false {
		return
	}

	fmt.Printf("[%s] "+colors[c].Color("(%s) %s\n"), time.Now().Format("15:04:05 02-01-2006"), nicenames[c], msg)
}

// C log commands used
func C(guild, user, command string) {
	fmt.Printf("[%s] "+colors[INFO].Color("(%s) ")+colors[MAGENTA].Color("%s ")+">> "+colors[INFO].Color("%s ")+"> "+colors[PLUGIN].Color("%s\n"), time.Now().Format("15:04:05 02-01-2006"), "COMMAND", guild, user, command)
}
