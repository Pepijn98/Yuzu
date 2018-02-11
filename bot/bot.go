package bot

import (
	"fmt"
	"time"
	"yuzu/logger"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	// BotID   string
	Session *discordgo.Session

	// StartTime when the application starts
	StartTime = time.Now()
)

// Connect to discord
func Connect(token string) {
	// Create a new Discord session using the provided bot token.
	var err error
	Session, err = discordgo.New("Bot " + token)
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error creating Discord session, %s", err))
		return
	}
	// Get the account information.
	_, err = Session.User("@me")
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error obtaining account details, %s", err))
	}
	// Store the account ID for later use.
	// BotID = u.ID
	logger.BOOT.L("Bot connected")
}

// Start discord session
func Start() {
	// Open the websocket and begin listening.
	err := Session.Open()
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error opening connection, %s", err))
		return
	}
	Session.Lock()
	Session.LogLevel = discordgo.LogError
	Session.ShouldReconnectOnError = true
	Session.Unlock()

	return
}

// Stop discord session
func Stop() {
	err := Session.Close()
	if err != nil {
		logger.ERROR.L(fmt.Sprintf("error closing session, %s", err))
		return
	}
}

// AddHandler adds the handler
func AddHandler(handler interface{}) {
	Session.AddHandler(handler)
}
