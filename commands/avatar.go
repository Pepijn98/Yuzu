package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"
	"yuzu/logger"
)

// Avatar sends someones avatar
type Avatar struct{}

// IsOwnerOnly f
func (Avatar) IsOwnerOnly() bool {
	return false
}

// Help f
func (Avatar) Help() [2]string {
	return [2]string{"Sends your avatar", ""}
}

// Process f
func (Avatar) Process(ctx yuzu.Context) {
	avatar := functions.GetAvatarURL(ctx.Message.Author)

	embed := yuzu.NewEmbed(ctx.Message.Author.Username + "#" + ctx.Message.Author.Discriminator)
	embed.Image(avatar)

	_, err := ctx.SayEmbed(embed)
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error sending message, %s", err))
		return
	}
}
