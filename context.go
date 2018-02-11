package yuzu

import (
	"fmt"
	"io"

	"github.com/bwmarrin/discordgo"
)

// Context f
type Context struct {
	Session *discordgo.Session
	State   *discordgo.State
	Message *discordgo.Message
	Author  *discordgo.User
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
	Args    []string
}

// NewContext creates a new `Context` instance.
func NewContext(
	session *discordgo.Session,
	message *discordgo.Message,
	channel *discordgo.Channel,
	guild *discordgo.Guild,
	args []string,
) Context {
	return Context{
		Session: session,
		State:   session.State,
		Message: message,
		Author:  message.Author,
		Channel: channel,
		Guild:   guild,
		Args:    args,
	}
}

// Say ayy lmao
func (context Context) Say(a ...interface{}) (*discordgo.Message, error) {
	return context.Session.ChannelMessageSend(context.Channel.ID, fmt.Sprint(a...))
}

// Sayf ayy lmao
func (context Context) Sayf(format string, a ...interface{}) (*discordgo.Message, error) {
	return context.Session.ChannelMessageSend(context.Channel.ID, fmt.Sprintf(format, a...))
}

// Error ayy lmao
func (context Context) Error(err error) (*discordgo.Message, error) {
	return context.Sayf("An error occurred in the command: %s", err)
}

// SayEmbed ayy lmao
func (context Context) SayEmbed(msg *Embed) (*discordgo.Message, error) {
	return context.Session.ChannelMessageSendEmbed(context.Channel.ID, msg.MessageEmbed)
}

// SayFile ayy lmao
func (context Context) SayFile(name string, reader io.Reader) (*discordgo.Message, error) {
	return context.Session.ChannelFileSend(context.Channel.ID, name, reader)
}

// Edit ayy lmao
func (context Context) Edit(msgID, text string) (*discordgo.Message, error) {
	return context.Session.ChannelMessageEdit(context.Channel.ID, msgID, text)
}

// EditEmbed ayy lmao
func (context Context) EditEmbed(sentMsgID string, msg *Embed) (*discordgo.Message, error) {
	return context.Session.ChannelMessageEditEmbed(context.Channel.ID, sentMsgID, msg.MessageEmbed)
}

// GetMessages ayy lmao
func (context Context) GetMessages(limit int) ([]*discordgo.Message, error) {
	return context.Session.ChannelMessages(context.Channel.ID, limit, "", "", "")
}

// BulkDelete ayy lmao
func (context Context) BulkDelete(messages []*discordgo.Message) error {
	var ids []string
	for _, msg := range messages {
		ids = append(ids, msg.ID)
	}
	return context.Session.ChannelMessagesBulkDelete(context.Channel.ID, ids)
}

// ChannelPermissions 0/1 args is for bot perms, 2 args is for user perms
func (context Context) ChannelPermissions(args ...string) (int, error) {
	switch {
	case len(args) == 0:
		return context.State.UserChannelPermissions(context.State.User.ID, context.Channel.ID)
	case len(args) == 1:
		return context.State.UserChannelPermissions(context.State.User.ID, args[0])
	default:
		return context.State.UserChannelPermissions(args[0], args[1])
	}
}

// JoinVoiceChannel ayy lmao
func (context Context) JoinVoiceChannel(id string) (*discordgo.VoiceConnection, error) {
	return context.Session.ChannelVoiceJoin(context.Guild.ID, id, false, false)
}
