package commands

import (
	"fmt"
	"strings"
	"yuzu/functions"

	"yuzu"

	"github.com/ajanicij/goduckgo/goduckgo"
)

// Duckduckgo uses duckduckgo's api to search things.
type Duckduckgo struct{}

// IsOwnerOnly f
func (Duckduckgo) IsOwnerOnly() bool {
	return false
}

// Help f
func (Duckduckgo) Help() [2]string {
	return [2]string{"Searches the web with duckduckgo", "<query>"}
}

// Process f
func (Duckduckgo) Process(ctx yuzu.Context) {
	query := strings.Join(ctx.Args, " ")
	msg, err := goduckgo.Query(query)
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/duckduckgo.go")
		_, e := ctx.Error(err)
		if e != nil {
			return
		}
		return
	}
	thumbnail := "https://images.duckduckgo.com/iu/?u=http%3A%2F%2Fcore2.staticworld.net%2Fimages%2Farticle%2F2014%2F05%2Fduckduckgo-logo-100266737-large.png&f=1"
	var onlyRedirect bool
	if strings.HasPrefix(query, "!") {
		onlyRedirect = true
		thumbnail = ""
	}
	embed := yuzu.NewEmbed("Results for your query")
	var description string
	if msg.Redirect != "" && onlyRedirect {
		description = msg.Redirect
	}
	if msg.Results != nil && len(msg.Results) != 0 && !onlyRedirect {
		for index, result := range msg.Results {
			if len(msg.Results) > 4 && index > 4 {
				break
			}
			if !result.Icon.IsEmpty() {
				thumbnail = result.Icon.URL
			}
			embed.Field("Result", "["+result.Text+"]("+result.FirstURL+")", true)
		}
	}
	if msg.AbstractText != "" && !onlyRedirect {
		embed.Field("Summary", msg.AbstractText, true)
	} else if msg.AbstractURL != "" && msg.AbstractText != "" && !onlyRedirect {
		embed.Field("Summary", "[Link]"+msg.AbstractURL+")\n\n"+msg.AbstractText, true)
	}
	if msg.RelatedTopics != nil && len(msg.RelatedTopics) != 0 && !onlyRedirect {
		var topics []string
		for index, topic := range msg.RelatedTopics {
			if len(msg.RelatedTopics) > 1 && index > 1 {
				break
			}
			topics = append(topics, fmt.Sprintf("[%s](%s)", topic.Text, topic.FirstURL))
		}
		embed.Field("Related topics", strings.Join(topics, "\n\n"), true)
	}
	if msg.Type != "" && !onlyRedirect {
		embed.Field("Category", toCategory(msg.Type), true)
	}
	if msg.Heading != "" && !onlyRedirect {
		embed.Field("Main topic", msg.Heading, true)
	}
	if msg.Answer != "" && !onlyRedirect {
		embed.Field("Answer", msg.Answer, true)
	}
	if msg.Definition != "" && !onlyRedirect {
		embed.Field("Definition", msg.Definition, true)
	}
	if msg.Image != "" && !onlyRedirect {
		thumbnail = msg.Image
	}
	embed.Thumbnail(thumbnail)
	embed.Color = 0xF5C7D9
	if onlyRedirect {
		embed.Title = "Redirection for your query"
		embed.Description = description
	}
	ctx.SayEmbed(embed)
}

func toCategory(msgtype string) string {
	switch msgtype {
	case "A":
		return "article"
	case "C":
		return "category"
	case "D":
		return "disambiguation"
	case "E":
		return "exclusive"
	case "N":
		return "name"
	default:
		return "nothing"
	}
}
