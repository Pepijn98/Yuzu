package commands

import (
	"strings"
	"yuzu"

	"github.com/bwmarrin/discordgo"
)

// Kick f
type Kick struct{}

// IsOwnerOnly f
func (Kick) IsOwnerOnly() bool {
	return false
}

// Help f
func (Kick) Help() [2]string {
	return [2]string{"Kicks the mentioned user", "<@mention> [reason]"}
}

// Process f
func (Kick) Process(ctx yuzu.Context) {
	uPerms, err := ctx.ChannelPermissions(ctx.Message.Author.ID, ctx.Channel.ID)
	if err != nil {
		_, e := ctx.Error(err)
		if e != nil {
			return
		}
		return
	}
	bPerms, err := ctx.ChannelPermissions()
	if err != nil {
		_, e := ctx.Error(err)
		if e != nil {
			return
		}
		return
	}
	if uPerms&discordgo.PermissionKickMembers == 0 {
		_, err := ctx.Say("You don't have the **kick members** permission.")
		if err != nil {
			_, e := ctx.Error(err)
			if e != nil {
				return
			}
			return
		}
		return
	}

	if bPerms&discordgo.PermissionKickMembers == 0 {
		_, err := ctx.Say("I don't have the **kick members** permission.")
		if err != nil {
			_, e := ctx.Error(err)
			if e != nil {
				return
			}
			return
		}
		return
	}

	if len(ctx.Args) < 1 {
		_, err := ctx.Say("Please mention a member to kick")
		if err != nil {
			_, e := ctx.Error(err)
			if e != nil {
				return
			}
			return
		}
		return
	}

	if len(ctx.Message.Mentions) < 1 {
		_, err := ctx.Say("Please mention a member to kick")
		if err != nil {
			_, e := ctx.Error(err)
			if e != nil {
				return
			}
			return
		}
		return
	}

	var reason = make([]string, 0)

	if len(ctx.Args) == 1 {
		reason = append(reason, "No reason provided. Responsible moderator: "+ctx.Author.Username+"#"+ctx.Author.Discriminator)
	}

	first := true
	for _, val := range ctx.Args {
		if first {
			first = false
			continue
		}
		reason = append(reason, val)
	}

	e := ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID,
		ctx.Message.Mentions[0].ID,
		strings.Join(reason, " ")+". Responsible moderator: "+ctx.Author.Username+"#"+ctx.Author.Discriminator)
	if e != nil {
		_, err := ctx.Error(e)
		if err != nil {
			return
		}
		return
	}
}
