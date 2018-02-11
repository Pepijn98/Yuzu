package commands

import (
	"fmt"
	"strings"
	"yuzu"
	//"yuzu/config"
	"yuzu/functions"

	//"github.com/CleverbotIO/go-cleverbot.io"
	"github.com/Jeffail/gabs"
	"yuzu/logger"
)

// Talk things
type Talk struct{}

// IsOwnerOnly f
func (Talk) IsOwnerOnly() bool {
	return false
}

// Help f
func (Talk) Help() [2]string {
	return [2]string{"Talk with Yuzu", "<query>"}
}

// Process f
func (Talk) Process(ctx yuzu.Context) {
	if len(ctx.Args) == 0 {
		_, err := ctx.Say("What do you want to talk about?")
		if err != nil {
			_, e := ctx.Error(err)
			if e != nil {
				return
			}
			return
		}
		return
	}

	ctx.Session.ChannelTyping(ctx.Channel.ID)

	/* For whenever program-o is down
	s, e := cleverbot.New(config.Config.CleverbotUser, config.Config.CleverbotKey)
	if e != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", e), "/commands/talk.go")
		_, err := ctx.Say("An error occured:\n", e)
		if err != nil {
			return
		}
		return
	}

	reply, err := s.Ask(strings.Join(ctx.Args, " "))
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/talk.go")
		_, e := ctx.Say("An error occured:\n", err)
		if e != nil {
			return
		}
		return
	}

	_, er := ctx.Say(reply)
	if er != nil {
		return
	}
	*/

	url := "http://api.program-o.com/v2/chatbot/?bot_id=12&say=" + strings.Join(ctx.Args, " ") + "&convo_id=" + "YUZU-" + ctx.Message.Author.ID + "&format=json"
	res, er := functions.GET(url)
	if er != nil {
		_, e := ctx.Say("Something went wrong while getting an image %s", er)
		if e != nil {
			return
		}
	}
	json, e := gabs.ParseJSON(res)
	if e != nil {
		_, er := ctx.Error(e)
		if er != nil {
			logger.ERROR.L(fmt.Sprintf("%s", er))
			return
		}
		return
	}
	reply := json.Path("botsay").Data().(string)
	reply = strings.Replace(reply, "Elizabeth", "Kurozero#0001", -1)
	reply = strings.Replace(reply, "Chatmundo", ctx.State.User.Username, -1)
	reply = strings.Replace(reply, "<br/>", "\n", -1)
	_, err := ctx.Say(reply)
	if err != nil {
		_, e := ctx.Error(err)
		if e != nil {
			logger.ERROR.L(fmt.Sprintf("%s", err))
			return
		}
		return
	}
}
