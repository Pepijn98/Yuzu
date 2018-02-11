package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"

	"github.com/KurozeroPB/go-weeb"
)

// Smug smug
type Smug struct{}

// IsOwnerOnly f
func (Smug) IsOwnerOnly() bool {
	return false
}

// Help f
func (Smug) Help() [2]string {
	return [2]string{"Sends s smug face", ""}
}

// Process f
func (Smug) Process(ctx yuzu.Context) {
	// Do the stuff
	img, err := weeb.GetImage("smug")
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/smug.go")
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
