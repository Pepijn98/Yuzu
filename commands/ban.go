package commands

import (
	"strings"
	"yuzu"

	"github.com/bwmarrin/discordgo"
)

// Ban f
type Ban struct{}

// IsOwnerOnly f
func (Ban) IsOwnerOnly() bool {
	return false
}

// Help f
func (Ban) Help() [2]string {
	return [2]string{"Bans the mentioned user", "<@mention> <reason>"}
}

// Process f
func (Ban) Process(ctx yuzu.Context) {
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
	if uPerms&discordgo.PermissionBanMembers == 0 {
		_, err := ctx.Say("You don't have the **ban members** permission.")
		if err != nil {
			_, e := ctx.Error(err)
			if e != nil {
				return
			}
			return
		}
		return
	}

	if bPerms&discordgo.PermissionBanMembers == 0 {
		_, err := ctx.Say("I don't have the **ban members** permission.")
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
		_, err := ctx.Say("Please mention a member to ban")
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
		_, err := ctx.Say("Please mention a member to ban")
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

	e := ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID,
		ctx.Message.Mentions[0].ID,
		strings.Join(reason, " ")+". Responsible moderator: "+ctx.Author.Username+"#"+ctx.Author.Discriminator, 7)
	if e != nil {
		_, err := ctx.Error(e)
		if err != nil {
			return
		}
		return
	}
}
