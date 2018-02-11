package commands

import (
	"fmt"
	"yuzu"
	"yuzu/config"
	"yuzu/functions"
	"yuzu/logger"
)

// About sends some info about the bot
type About struct{}

// IsOwnerOnly f
func (About) IsOwnerOnly() bool {
	return false
}

// Help f
func (About) Help() [2]string {
	return [2]string{"Sends some info about the bot", ""}
}

// Process f
func (About) Process(ctx yuzu.Context) {
	u, err := ctx.Session.User("@me")
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error getting bot user, %s", err))
		functions.ReportError(ctx.Session, fmt.Sprintf("error getting bot user, %s", err), "/commands/about.go")
		return
	}
	var avatar = "https://cdn.discordapp.com/avatars/" + u.ID + "/" + u.Avatar + ".png?size=2048"

	embed := yuzu.NewEmbed(u.Username)
	embed.Description = "Yuzu is a simple bot created in Go using the DiscordGo lib."
	embed.Field("Owner", "Kurozero#0001", true)
	embed.Field("Version", "v"+config.Config.Version, true)
	embed.Thumbnail(avatar)

	_, e := ctx.SayEmbed(embed)
	if e != nil {
		logger.ERROR.L(fmt.Sprintf("error sending message, %s", e))
		functions.ReportError(ctx.Session, fmt.Sprintf("error sending message, %s", e), "/commands/about.go")
		return
	}
}
