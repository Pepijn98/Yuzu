package bot

import (
	"fmt"
	"time"
	"yuzu/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dshardmanager"
	"yuzu/config"
)

// Variables used for command line parameters
var (
	// Manager the shard manager
	Manager *dshardmanager.Manager

	// StartTime when the application starts
	StartTime = time.Now()
)

// Connect to discord
func Connect(token string) {
	// Create a new Discord session using the provided bot token.
	Manager = dshardmanager.New("Bot " + token)

	shardCount, err := Manager.GetRecommendedCount()
	if err != nil {
		shardCount = 2
		logger.ERROR.L(fmt.Sprintf("Failed getting recommended shard count, using static: %v", shardCount))
	}
	Manager.SetNumShards(shardCount)
	Manager.Name = "Yuzu | 柚"
	Manager.LogChannel = "417460849561698308"
	Manager.StatusMessageChannel = config.StatusMsgChannel

	Manager.SessionFunc = func(token string) (Session *discordgo.Session, err error) {
		Session, err = discordgo.New(token)
		if err != nil {
			logger.ERROR.L("Failed to start a session")
			return
		}

		Session.LogLevel = discordgo.LogWarning
		Session.ShouldReconnectOnError = true
		Session.SyncEvents = true
		return
	}

	logger.BOOT.L("Bot connected")
}

// Start discord session
func Start() {
	// Open the websocket and begin listening.
	err := Manager.Start()
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error opening connection, %s", err))
		return
	}
	return
}

// Stop discord session
func Stop() {
	err := Manager.StopAll()
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error closing sessions, %s", err))
		return
	}
}

// AddHandler adds the handler
func AddHandler(handler interface{}) {
	Manager.AddHandler(handler)
}
