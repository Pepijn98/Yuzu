package commands

import (
	"fmt"
	"yuzu"
	"yuzu/functions"

	"github.com/KurozeroPB/go-weeb"
)

// Hug hug someone
type Hug struct{}

// IsOwnerOnly f
func (Hug) IsOwnerOnly() bool {
	return false
}

// Help f
func (Hug) Help() [2]string {
	return [2]string{"Give someone a nice hug", "<@mention>"}
}

// Process f
func (Hug) Process(ctx yuzu.Context) {
	channel, err := ctx.Session.Channel(ctx.Message.ChannelID)
	if err != nil {
		return
	}
	guildID := channel.GuildID
	guild, err := ctx.Session.Guild(guildID)
	if err != nil {
		return
	}

	mentionedUsers := len(ctx.Message.Mentions)
	if mentionedUsers == 0 {
		_, er := ctx.Say("Please mention a member to hug.")
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
			_, e := ctx.Say(err)
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
			_, e := ctx.Say(err)
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
		img, err := weeb.GetImage("hug")
		if err != nil {
			functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/hug.go")
			_, e := ctx.Say("Error: ", err)
			if e != nil {
				return
			}
			return
		}
		embed := yuzu.NewEmbed(author + " hugs " + mentioned)
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
