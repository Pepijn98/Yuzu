package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"yuzu"
	"yuzu/bot"
	"yuzu/commands"
	"yuzu/config"
	"yuzu/logger"
)

var requestOfflineUsers bool

func init() {
	flag.BoolVar(&requestOfflineUsers, "requestoffusers", false, "")
	flag.Parse()
}

func main() {
	err := config.Configure()
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("%s", err))
	}

	bot.Connect(config.Token)

	listOfCommands := map[string]yuzu.Command{
		"ping":    commands.Ping{},
		"stats":   commands.Stats{},
		"eval":    commands.Eval{},
		"cat":     commands.Cat{},
		"ddg":     commands.Duckduckgo{},
		"calc":    commands.Calc{},
		"about":   commands.About{},
		"avatar":  commands.Avatar{},
		"catgirl": commands.Catgirl{},
		"exec":    commands.Exec{},
		"pat":     commands.Pat{},
		"cuddle":  commands.Cuddle{},
		"hug":     commands.Hug{},
		"kiss":    commands.Kiss{},
		"slap":    commands.Slap{},
		"pout":    commands.Pout{},
		"smug":    commands.Smug{},
		"stare":   commands.Stare{},
		"ban":     commands.Ban{},
		"kick":    commands.Kick{},
		"talk":    commands.Talk{},
		"anime":   commands.Anime{},
		"manga":   commands.Manga{},
		"char":    commands.Character{},
	}
	listOfCommands["help"] = commands.Help{Commands: listOfCommands}
	bot.AddHandler(yuzu.Ready(requestOfflineUsers))
	bot.AddHandler(yuzu.GuildMemberChunk)
	bot.AddHandler(yuzu.MessageCreate(listOfCommands))
	bot.AddHandler(yuzu.GuildCreate)
	bot.AddHandler(yuzu.GuildDelete)

	bot.Start()

	// Clean exit
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill)
	<-channel

	logger.WARNING.L("Disconnecting...")
	bot.Stop()
}
