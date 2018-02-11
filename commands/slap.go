package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"
	"yuzu/logger"

	"github.com/KurozeroPB/go-weeb"
)

// Slap slap someone
type Slap struct{}

// IsOwnerOnly f
func (Slap) IsOwnerOnly() bool {
	return false
}

// Help f
func (Slap) Help() [2]string {
	return [2]string{"Slap someone", "<@mention>"}
}

// Process f
func (Slap) Process(ctx yuzu.Context) {
	channel, err := ctx.Session.Channel(ctx.Message.ChannelID)
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("%s", err))
		return
	}
	guildID := channel.GuildID
	guild, err := ctx.Session.Guild(guildID)
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("%s", err))
		return
	}

	mentionedUsers := len(ctx.Message.Mentions)
	if mentionedUsers == 0 {
		_, er := ctx.Say("Please mention a member to slap.")
		if er != nil {
			return
		}
	} else if mentionedUsers == 1 {
		var (
			author    string
			mentioned string
		)
		// Get the Member of the Author user
		authMember, err := ctx.Session.GuildMember(guild.ID, ctx.Message.Author.ID)
		if err != nil {
			logger.ERROR.L(fmt.Sprintf("%s", err))
			return
		}
		if authMember.Nick == "" {
			author = ctx.Message.Author.Username
		} else {
			author = authMember.Nick
		}
		// Get the Member of the mentioned user
		mentionedMember, err := ctx.Session.GuildMember(guild.ID, ctx.Message.Mentions[0].ID)
		if err != nil {
			logger.ERROR.L(fmt.Sprintf("%s", err))
			return
		}
		if mentionedMember.Nick == "" {
			mentioned = ctx.Message.Mentions[0].Username
		} else {
			mentioned = mentionedMember.Nick
		}
		// Do the stuff
		img, err := weeb.GetImage("slap")
		if err != nil {
			functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/slap.go")
			_, e := ctx.Say("Error: ", err)
			if e != nil {
				return
			}
			return
		}
		embed := yuzu.NewEmbed(author + " slapped " + mentioned)
		embed.Image(img)
		_, er := ctx.SayEmbed(embed)
		if er != nil {
			_, e := ctx.Say("Error sending embed, ", er)
			if e != nil {
				return
			}
			return
		}
	} else if mentionedUsers > 1 {
		_, err := ctx.Say("Please only mention one user at a time.")
		if err != nil {
			return
		}
	}
}
