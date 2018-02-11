package yuzu

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"yuzu/config"
	"yuzu/logger"

	"github.com/bwmarrin/discordgo"
)

// CommandCounter is the default instance of the commandCounter struct.
var CommandCounter = commandCounter{
	Counter: map[string]int{},
}

type commandCounter struct {
	sync.Mutex
	Counter map[string]int
}

// Update updates the underlaying counter with the provided command.
func (c *commandCounter) Update(command string) {
	c.Lock()
	c.Counter[command]++
	c.Unlock()
}

// Panicf is a shorthand to `panic(fmt.Sprintf(...))`.
func Panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

// ISOTimestamp returns an ISO8601 based timestamp from the current time.
func ISOTimestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.070000")
}

// VoiceChannelID gets the voice channel id of the author.
func VoiceChannelID(context Context) string {
	for _, vc := range context.Guild.VoiceStates {
		if vc.UserID == context.Author.ID {
			return vc.ChannelID
		}
	}
	return ""
}

// LogString gets the dynamically changed log string for logging messages.
func LogString(author *discordgo.User, command string, channel *discordgo.Channel, guild *discordgo.Guild, args string) {
	switch {
	case guild != nil:
		logger.C(guild.Name, author.Username, config.Prefix+command+" "+args)
	case channel.Type == 1:
		logger.C("(dm)", author.Username, config.Prefix+command+" "+args)
	default:
		return
	}
}

func checkMessage(msg *discordgo.Message, userID string) bool {
	return msg.Author.ID == userID || msg.Author.Bot || !strings.HasPrefix(msg.Content, config.Prefix)
}

func guildFromState(channel *discordgo.Channel, state *discordgo.State) *discordgo.Guild {
	if channel.GuildID == "" && channel.Type == 1 {
		return nil
	}
	guild, _ := state.Guild(channel.GuildID)
	return guild
}

// Embed is a little wrapper around discordgo's `MessageEmbed`.
type Embed struct {
	*discordgo.MessageEmbed
}

// Field is a shortcut to field boilerplate.
func (e *Embed) Field(name, value string, inline bool) {
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{Name: name, Value: value, Inline: inline})
}

// Thumbnail yet another shortcut.
func (e *Embed) Thumbnail(url string) {
	e.MessageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: url}
}

// Image yet another shortcut.
func (e *Embed) Image(url string) {
	e.MessageEmbed.Image = &discordgo.MessageEmbedImage{URL: url}
}

// Footer YEEEET another shortcut.
func (e *Embed) Footer(text string) {
	e.MessageEmbed.Footer = &discordgo.MessageEmbedFooter{Text: text}
}

// NewEmbed creates `Embed`.
// Assings `Color` to the colour of the blue colour and Timestamp to the current time.
//
// Accepts an optional `title`
func NewEmbed(title string) *Embed {
	embed := &Embed{&discordgo.MessageEmbed{Color: 0xF5C7D9}}
	if title == "" {
		return embed
	}
	embed.Title = title
	return embed
}
