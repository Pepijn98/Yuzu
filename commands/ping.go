package commands

import (
	"strconv"
	"time"

	"yuzu"
)

// Ping is used for measuring time between the location of the bot and discord's servers.
type Ping struct{}

// IsOwnerOnly f
func (Ping) IsOwnerOnly() bool {
	return false
}

// Help f
func (Ping) Help() [2]string {
	return [2]string{"Check if the bot works and its response time", ""}
}

// Process f
func (Ping) Process(ctx yuzu.Context) {
	start := time.Now()

	m, err := ctx.Say("Pong!")
	if err != nil {
		return
	}

	end := time.Now()

	_, e := ctx.Edit(m.ID, "Pong!\n\nRest Latency: **"+strconv.Itoa(int(end.Sub(start)/time.Millisecond)/2)+"** ms\n")
	if e != nil {
		return
	}
}
