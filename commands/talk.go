package commands

import (
	"strings"
	"yuzu"
	"yuzu/functions"
	"github.com/Jeffail/gabs"
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

	url := "http://api.program-o.com/v2/chatbot/?bot_id=12&say=" + strings.Join(ctx.Args, " ") + "&convo_id=" + "YUZU-" + ctx.Message.Author.ID + "&format=json"
	res, er := functions.GET(url)
	if er != nil {
		_, e := ctx.Error(er)
		if e != nil {
			return
		}
	}
	json, e := gabs.ParseJSON(res)
	if e != nil {
		_, er := ctx.Error(e)
		if er != nil {
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
			return
		}
		return
	}
}
