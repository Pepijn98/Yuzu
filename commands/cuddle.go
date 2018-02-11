package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"
	"yuzu/logger"

	"github.com/KurozeroPB/go-weeb"
)

// Cuddle cuddle someone
type Cuddle struct{}

// IsOwnerOnly f
func (Cuddle) IsOwnerOnly() bool {
	return false
}

// Help f
func (Cuddle) Help() [2]string {
	return [2]string{"Give someone a sweet cuddle", "<@mention>"}
}

// Process f
func (Cuddle) Process(ctx yuzu.Context) {
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
		_, er := ctx.Say("Please mention a member to cuddle.")
		if er != nil {
			logger.ERROR.L(fmt.Sprintf("error sending message, %s", er))
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
			functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/cuddle.go")
			_, e := ctx.Say("Error: ", err)
			if e != nil {
				return
			}
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
			functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/cuddle.go")
			_, e := ctx.Say("Error: ", err)
			if e != nil {
				return
			}
			return
		}
		if mentionedMember.Nick == "" {
			mentioned = ctx.Message.Mentions[0].Username
		} else {
			mentioned = mentionedMember.Nick
		}
		// Do the stuff
		img, err := weeb.GetImage("cuddle")
		if err != nil {
			functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/cuddle.go")
			_, e := ctx.Say("Error: ", err)
			if e != nil {
				return
			}
			return
		}
		embed := yuzu.NewEmbed(author + " gave " + mentioned + " a sweet cuddle")
		embed.Image(img)
		_, er := ctx.SayEmbed(embed)
		if er != nil {
			_, e := ctx.Say("error sending message, ", er)
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
