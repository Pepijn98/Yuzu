package functions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"yuzu/config"
	"yuzu/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron"
)

var (
	timeNow   string
	imgFormat string
	imgSize   string
	avatarURL string
	iconURL   string
	// USERAGENT User Agent
	USERAGENT = "yuzu/" + config.Config.Version
)

// GetTime Gets the current time + date
func GetTime() string {
	timeNow = time.Now().Format("02/01/2006 15:04:05")
	return timeNow
}

// GET make the get request
func GET(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", USERAGENT)
	request.Header.Set("Accept", "application/vnd.api+json")
	request.Header.Set("Content-Type", "application/vnd.api+json")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Expected status %d; Got %d\nResponse: %#v", 200, response.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}

// CommandsList array of all the commands
func CommandsList() []string {
	commandsList := []string{}
	searchDir := "$HOME/go/src/yuzu/commands/"

	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		commandsList = append(commandsList, path)
		return nil
	})
	if err != nil {
		e := fmt.Sprintf("%s", err)
		logger.ERROR.L(e)
	}

	commandsList = append(commandsList[:0], commandsList[0+1:]...)

	for i := range commandsList {
		commandsList[i] = strings.Replace(commandsList[i], "$HOME\\go\\src\\yuzu\\commands\\", "", -1)
		commandsList[i] = strings.Replace(commandsList[i], ".go", "", -1)
	}
	return commandsList
}

// GetAvatarURL gets the user's avatar
func GetAvatarURL(user *discordgo.User) string {
	if user.Avatar == "" {
		avatarURL = "https://b.catgirlsare.sexy/ILuV.jpg"
		return avatarURL
	}
	imgFormat = "png"
	if strings.HasPrefix(user.Avatar, "a_") {
		imgFormat = "gif"
	}
	avatarURL = "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + "." + imgFormat + "?size=2048"
	return avatarURL
}

// GetGuildIcon gets the guild's icon
func GetGuildIcon(id string, icon string) string {
	if icon == "" {
		iconURL = "https://b.catgirlsare.sexy/ILuV.jpg"
		return iconURL
	}
	iconURL = "https://cdn.discordapp.com/icons/" + id + "/" + icon + ".png?size=2048"
	return iconURL
}

// PostGuildCount posts guild count to discordbots.org and bots.discord.pw every 15 minutes
func PostGuildCount(s *discordgo.Session) {
	if config.Config.Release == "dev" {
		return
	}
	c := cron.New()
	c.AddFunc("@every 15m", func() {
		gCount, err := getCurrentDBotsOrgStats(s)
		if err != nil {
			logger.ERROR.L(fmt.Sprintf("%s", err))
			return
		}
		guilds := s.State.Guilds
		guildCount := len(guilds)
		if gCount == guildCount {
			return
		}
		postDBotsOrg(s)
		postDBotsPW(s)
	})
	c.Start()
}

func postDBotsOrg(s *discordgo.Session) {
	if config.Config.DBotsOrgKey == "" {
		logger.WARNING.L("No token provided for discordbots.org")
		return
	}
	guilds := s.State.Guilds
	guildCount := len(guilds)
	url := "https://discordbots.org/api/bots/" + s.State.User.ID + "/stats"
	postGuildCount := []byte(`{"server_count":"` + strconv.Itoa(guildCount) + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postGuildCount))
	req.Header.Set("Authorization", config.Config.DBotsOrgKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Sprintf("error sending server count, %s", err)
		logger.ERROR.L(e)
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		logger.ERROR.L("Error posting server count to discordbots.org\nStatus code: " + resp.Status)
	} else {
		logger.INFO.L("Success posting server count to discordbots.org\nStatus code: " + resp.Status)
	}
}

func postDBotsPW(s *discordgo.Session) {
	if config.Config.DBotsPWKey == "" {
		logger.WARNING.L("No token provided for bots.discord.pw")
		return
	}
	guilds := s.State.Guilds
	guildCount := len(guilds)
	url := "https://bots.discord.pw/api/bots/" + s.State.User.ID + "/stats"
	postGuildCount := []byte(`{"server_count":"` + strconv.Itoa(guildCount) + `"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postGuildCount))
	req.Header.Set("Authorization", config.Config.DBotsPWKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Sprintf("error sending server count, %s", err)
		logger.ERROR.L(e)
	}
	defer resp.Body.Close()
	if resp.Status != "200 OK" {
		logger.ERROR.L("Error posting server count to bots.discord.pw\nStatus code: " + resp.Status)
	} else {
		logger.INFO.L("Success posting server count to bots.discord.pw\nStatus code: " + resp.Status)
	}
}

func getCurrentDBotsOrgStats(s *discordgo.Session) (int, error) {
	url := "https://discordbots.org/api/bots/" + s.State.User.ID + "/stats"
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return 0, e
	}

	type statsStruct struct {
		ServerCount int      `json:"server_count"`
		Shards      []string `json:"shards"`
	}

	// Stats variables used for command line parameters
	var Stats statsStruct

	jsonError := json.Unmarshal(body, &Stats)
	if jsonError != nil {
		return 0, jsonError
	}
	return Stats.ServerCount, jsonError
}

// UpdateStatus updates the bots playing message every 5 minutes
func UpdateStatus(s *discordgo.Session) {
	// Status update every 5 minutes
	c := cron.New()
	c.AddFunc("@every 5m", func() {
		var games = make([]string, 0)
		users := make(map[string]string)
		guilds := s.State.Guilds
		for _, guild := range guilds {
			for _, u := range guild.Members {
				users[u.User.ID] = u.User.Username
			}
		}
		var guildCount = strconv.Itoa(len(guilds))
		var userCount = strconv.Itoa(len(users))
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
	c.Start()
}

// ReportError executes the error webhook on my private guild to easily track errors
func ReportError(s *discordgo.Session, msg string, filePath string) {
	if strings.Contains(msg, "Missing Access") {
		return
	}

	webhook := &discordgo.WebhookParams{
		Username:  s.State.User.Username,
		AvatarURL: s.State.User.AvatarURL("2048"),
		Embeds: []*discordgo.MessageEmbed{
			{
				Color:       0xff0000,
				Title:       "ERROR",
				Description: "```glsl\n**" + filePath + "\n\n" + msg + "```",
			},
		},
	}

	err := s.WebhookExecute(config.Config.ErrWebhookID, config.Config.ErrWebhookToken, true, webhook)
	if err != nil {
		return
	}
}

// WebhookGuildCreate executes "joined guild" webhook
func WebhookGuildCreate(s *discordgo.Session, guild *discordgo.GuildCreate) {
	webhook := &discordgo.WebhookParams{
		Username:  s.State.User.Username,
		AvatarURL: s.State.User.AvatarURL("2048"),
		Embeds: []*discordgo.MessageEmbed{
			{
				Color:       0xF5C7D9,
				Title:       "Joined Guild:",
				Description: fmt.Sprintf("**%s (%s)**", guild.Name, guild.ID),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: GetGuildIcon(guild.ID, guild.Icon),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Members",
						Value:  strconv.Itoa(guild.MemberCount),
						Inline: true,
					}, {
						Name:   "Channels",
						Value:  strconv.Itoa(len(guild.Channels)),
						Inline: true,
					},
				},
			},
		},
	}

	err := s.WebhookExecute(config.Config.JoinLeaveWebhookID, config.Config.JoinLeaveWebhookToken, true, webhook)
	if err != nil {
		return
	}
}

// WebhookGuildDelete executes "left guild" webhook
func WebhookGuildDelete(s *discordgo.Session, guild *discordgo.GuildDelete) {
	webhook := &discordgo.WebhookParams{
		Username:  s.State.User.Username,
		AvatarURL: s.State.User.AvatarURL("2048"),
		Embeds: []*discordgo.MessageEmbed{
			{
				Color:       14038325,
				Title:       "Left Guild:",
				Description: fmt.Sprintf("**%s (%s)**", guild.Name, guild.ID),
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: GetGuildIcon(guild.ID, guild.Icon),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Members",
						Value:  strconv.Itoa(guild.MemberCount),
						Inline: true,
					}, {
						Name:   "Channels",
						Value:  strconv.Itoa(len(guild.Channels)),
						Inline: true,
					},
				},
			},
		},
	}

	err := s.WebhookExecute(config.Config.JoinLeaveWebhookID, config.Config.JoinLeaveWebhookToken, true, webhook)
	if err != nil {
		return
	}
}
