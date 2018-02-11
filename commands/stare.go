package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"

	"github.com/KurozeroPB/go-weeb"
)

// Stare stare
type Stare struct{}

// IsOwnerOnly f
func (Stare) IsOwnerOnly() bool {
	return false
}

// Help f
func (Stare) Help() [2]string {
	return [2]string{"Stare~~~~", ""}
}

// Process f
func (Stare) Process(ctx yuzu.Context) {
	// Do the stuff
	img, err := weeb.GetImage("stare")
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/stare.go")
		_, e := ctx.Say("Error: ", err)
		if e != nil {
			return
		}
		return
	}
	embed := yuzu.NewEmbed("")
	embed.Image(img)
	_, er := ctx.SayEmbed(embed)
	if er != nil {
		return
	}
}
