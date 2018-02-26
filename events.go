package yuzu

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"yuzu/bot"
	"yuzu/config"
	"yuzu/functions"
	"yuzu/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dshardmanager"
	"strconv"
)

var Manager *dshardmanager.Manager

// Ready handles the READY event.
func Ready(chunkMembers bool) func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		if chunkMembers {
			for _, g := range r.Guilds {
				if g.Large {
					s.RequestGuildMembers(g.ID, "", 0)
				}
			}
		}
		functions.PostGuildCount(s)
		functions.UpdateStatus(s)
		time.AfterFunc(5*time.Second, func() {
			// Count guilds, channels and users
			users := make(map[string]string)
			channels := 0
			guilds := s.State.Guilds

			for _, guild := range guilds {
				channels += len(guild.Channels)

				for _, u := range guild.Members {
					users[u.User.ID] = u.User.Username
				}
			}

			var guildCount = strconv.Itoa(len(guilds))
			var userCount = strconv.Itoa(len(users))

			// Initial status update
			var games = make([]string, 0)
			games = append(games,
				"prefix is y:",
				"with Senpai",
				"in "+guildCount+" servers",
				"with "+userCount+" users",
				"y:help")
			rand.Seed(time.Now().Unix())
			game := games[rand.Intn(len(games))]
			err := s.UpdateStatus(0, game)
			if err != nil {
				e := fmt.Sprintf("error setting status, %s", err)
				logger.ERROR.L(e)
			}
		})

		if s.ShardID == s.ShardCount-1 {
			logger.BOOT.L(fmt.Sprintf(
				"Logged in as: %s#%s (%s)",
				r.User.Username,
				r.User.Discriminator,
				r.User.ID),
			)
		}
	}
}

// GuildMemberChunk handles the GUILD_MEMBER_CHUNK event.
func GuildMemberChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
	for _, g := range s.State.Guilds {
		if g.ID == c.GuildID {
			g.Members = append(g.Members, c.Members...)
		}
	}
}

// MessageCreate handles the MESSAGE_CREATE event while also handling command parsing, execution, etc.
func MessageCreate(commands map[string]Command) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		defer func() {
			if r := recover(); r != nil {
				logger.INFO.L("Recovered")
			}
		}()
		if checkMessage(msg.Message, s.State.User.ID) {
			return
		}

		args := strings.Fields(msg.Content[len(config.Prefix):])
		command, args := args[0], args[1:]
		cmd, ok := commands[strings.ToLower(command)]
		if !ok {
			return
		}
		if cmd.IsOwnerOnly() && msg.Author.ID != config.Config.OwnerID {
			_, err := s.ChannelMessageSend(msg.ChannelID, "This command is for the bot owner only.")
			if err != nil {
				logger.ERROR.L(fmt.Sprintf("%s", err))
				functions.ReportError(s, fmt.Sprintf("%s", err), "/events.go")
				return
			}
			return
		}
		channel, err := s.State.Channel(msg.ChannelID)
		if err != nil {
			logger.ERROR.L(fmt.Sprintf("%s", err))
			functions.ReportError(s, fmt.Sprintf("%s", err), "/events.go")
			return
		}
		guild := guildFromState(channel, s.State)

		LogString(msg.Author, command, channel, guild, strings.Join(args, " "))
		CommandCounter.Update(command)
		go cmd.Process(NewContext(s, msg.Message, channel, guild, args))
	}
}

// GuildCreate discord event
func GuildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	if guild.Unavailable == false {
		joinedTime, _ := guild.JoinedAt.Parse()
		joinedSince := int64(time.Since(joinedTime) / time.Second)
		if joinedSince >= 20 {
			return
		}
		uptime := time.Since(bot.StartTime)
		sec := int64(uptime / time.Second)
		if sec <= 20 {
			return
		}
		logger.INFO.L(fmt.Sprintf("Joined guild: %s", guild.Name))
		functions.WebhookGuildCreate(s, guild)
	} else {
		return
	}
}

// GuildDelete discord event
func GuildDelete(s *discordgo.Session, guild *discordgo.GuildDelete) {
	if guild.Unavailable == false {
		uptime := time.Since(bot.StartTime)
		sec := int64(uptime / time.Second)
		if sec <= 15 {
			return
		}
		logger.INFO.L(fmt.Sprintf("Left guild: %s", guild.Name))
		functions.WebhookGuildDelete(s, guild)
	} else {
		return
	}
}
