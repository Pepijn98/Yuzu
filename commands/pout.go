package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"

	"github.com/KurozeroPB/go-weeb"
)

// Pout pout
type Pout struct{}

// IsOwnerOnly f
func (Pout) IsOwnerOnly() bool {
	return false
}

// Help f
func (Pout) Help() [2]string {
	return [2]string{"Sends an anime character with a pout face", ""}
}

// Process f
func (Pout) Process(ctx yuzu.Context) {
	// Do the stuff
	img, err := weeb.GetImage("pout")
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/pout.go")
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
		_, e := ctx.Say("Error sending embed, ", er)
		if e != nil {
			return
		}
		return
	}
}
