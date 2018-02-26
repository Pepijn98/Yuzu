package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"yuzu/logger"

	"errors"
)

// Stuff
var (
	Token  string
	Prefix string
	StatusMsgChannel string
)

// Config all the config shiz
var Config configStruct

// Config maps how a general yuzu config should be
type configStruct struct {
	Release               string `json:"release"`
	Token                 string `json:"token"`
	DevToken              string `json:"dev_token"`
	Prefix                string `json:"prefix"`
	DevPrefix             string `json:"dev_prefix"`
	OwnerID               string `json:"owner_id"`
	Version               string `json:"version"`
	DBotsOrgKey           string `json:"DBotsOrgKey"`
	DBotsPWKey            string `json:"DBotsPWKey"`
	CleverbotUser         string `json:"cb_user"`
	CleverbotKey          string `json:"cb_key"`
	ErrWebhookToken       string `json:"errWebhookToken"`
	ErrWebhookID          string `json:"errWebhookID"`
	JoinLeaveWebhookToken string `json:"join_leaveWebhookToken"`
	JoinLeaveWebhookID    string `json:"join_leaveWebhookID"`
	StatusMessageChannel string `json:"statusMessageChannel"`
	DevStatusMessageChannel string `json:"devStatusMessageChannel"`
}

// Configure reads from disk the config file, unmarshal's it into the Go struct.
func Configure() error {
	logger.BOOT.L("Reading config file...")
	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		err := fmt.Sprintf("File error: %v\n", e)
		logger.ERROR.L(err)
		os.Exit(1)
	}
	err := json.Unmarshal(file, &Config)
	if err != nil {
		e := fmt.Sprintf("%s", err)
		logger.ERROR.L(e)
	}
	if Config.Release == "" {
		return errors.New("config requires the release field to be set to either 'dev' or 'prod'")
	}
	if Config.Release == "dev" {
		if Config.DevToken == "" || Config.Prefix == "" || Config.OwnerID == "" {
			return errors.New("config requires the dev_token, dev_prefix and owner_id fields to be set")
		}
		logger.DebugMode = true
		Token = Config.DevToken
		Prefix = Config.DevPrefix
		StatusMsgChannel = Config.DevStatusMessageChannel
	} else if Config.Release == "prod" {
		if Config.Token == "" || Config.Prefix == "" || Config.OwnerID == "" {
			return errors.New("config requires the token, prefix and owner_id fields to be set")
		}
		Token = Config.Token
		Prefix = Config.Prefix
		StatusMsgChannel = Config.StatusMessageChannel
	} else {
		return errors.New("config requires the release field to be set to either 'dev' or 'prod'")
	}
	return nil
}
