package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"yuzu"
	"yuzu/functions"
	"yuzu/logger"

	"github.com/Jeffail/gabs"
)

// Catgirl sends sexy catgirl
type Catgirl struct{}

// IsOwnerOnly f
func (Catgirl) IsOwnerOnly() bool {
	return false
}

// Help f
func (Catgirl) Help() [2]string {
	return [2]string{"Sends a cute catgirl (sfw or nsfw)", "<sfw/nsfw>"}
}

// Process f
func (Catgirl) Process(ctx yuzu.Context) {
	var sites [2]string
	sites[0] = "nekos-life"
	sites[1] = "nekos"
	var suffix string
	if len(ctx.Args) == 0 {
		suffix = ""
	} else {
		suffix = ctx.Args[0]
	}
	// SFW
	if strings.ToLower(suffix) == "sfw" || suffix == "" {
		rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
		randSite := fmt.Sprint(sites[rand.Intn(len(sites))])

		if randSite == "nekos-life" {
			url := "https://nekos.life/api/neko"
			res, err := functions.GET(url)
			if err != nil {
				functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/catgirl.go")
				_, e := ctx.Say("Something went wrong while getting an image ", err)
				if e != nil {
					return
				}
			}
			json, e := gabs.ParseJSON(res)
			if e != nil {
				logger.ERROR.L(fmt.Sprintf("%s", e))
				functions.ReportError(ctx.Session, fmt.Sprintf("%s", e), "/commands/catgirl.go")
				return
			}
			img := json.Path("neko").Data().(string)

			embed := yuzu.NewEmbed("")
			embed.Image(img)
			_, er := ctx.SayEmbed(embed)
			if er != nil {
				_, e := ctx.Say("error sending embed, ", er)
				if e != nil {
					return
				}
				return
			}
		} else if randSite == "nekos" {
			url := "https://nekos.brussell.me/api/v1/random/image?nsfw=false"
			res, err := functions.GET(url)
			if err != nil {
				functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/catgirl.go")
				_, e := ctx.Say("Something went wrong while getting an image ", err)
				if e != nil {
					return
				}
				return
			}
			json, e := gabs.ParseJSON(res)
			if e != nil {
				logger.ERROR.L(fmt.Sprintf("%s", e))
				functions.ReportError(ctx.Session, fmt.Sprintf("%s", e), "/commands/catgirl.go")
				return
			}
			inter := json.Path("images").Path("id").Data().([]interface{})
			img := inter[0].(string)

			embed := yuzu.NewEmbed("")
			embed.Image("https://nekos.brussell.me/image/" + img)
			_, er := ctx.SayEmbed(embed)
			if er != nil {
				_, e := ctx.Say("error sending embed, ", er)
				if e != nil {
					return
				}
				return
			}
		} else {
			_, e := ctx.Say("Something went wrong while picking the catgirl site")
			if e != nil {
				return
			}
			return
		}
	} else if strings.ToLower(suffix) == "nsfw" {
		channel, err := ctx.Session.Channel(ctx.Message.ChannelID)
		if err != nil {
			_, e := ctx.Say("error sending embed, ", err)
			if e != nil {
				return
			}
			return
		}
		if channel.NSFW != true {
			_, err := ctx.Say("The NSFW option can only be used in NSFW marked channels.\nPlease do this in the channel settings of Discord.")
			if err != nil {
				_, e := ctx.Say("error sending message, ", err)
				if e != nil {
					return
				}
				return
			}
		} else {
			rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
			randSite := fmt.Sprint(sites[rand.Intn(len(sites))])

			if randSite == "nekos-life" {
				url := "https://nekos.life/api/lewd/neko"
				res, err := functions.GET(url)
				if err != nil {
					functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/catgirl.go")
					_, e := ctx.Say("Something went wrong while getting an image ", err)
					if e != nil {
						return
					}
					return
				}
				json, e := gabs.ParseJSON(res)
				if e != nil {
					functions.ReportError(ctx.Session, fmt.Sprintf("%s", e), "/commands/catgirl.go")
					logger.ERROR.L(fmt.Sprintf("%s", e))
					return
				}
				img := json.Path("neko").Data().(string)

				embed := yuzu.NewEmbed("")
				embed.Image(img)
				_, er := ctx.SayEmbed(embed)
				if er != nil {
					_, e := ctx.Say("error sending embed, ", er)
					if e != nil {
						return
					}
					return
				}
			} else if randSite == "nekos" {
				url := "https://nekos.brussell.me/api/v1/random/image?nsfw=true"
				res, err := functions.GET(url)
				if err != nil {
					functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/catgirl.go")
					_, e := ctx.Say("Something went wrong while getting an image, ", err)
					if e != nil {
						return
					}
					return
				}
				json, e := gabs.ParseJSON(res)
				if e != nil {
					functions.ReportError(ctx.Session, fmt.Sprintf("%s", e), "/commands/catgirl.go")
					logger.ERROR.L(fmt.Sprintf("%s", e))
					return
				}
				inter := json.Path("images").Path("id").Data().([]interface{})
				img := inter[0].(string)

				embed := yuzu.NewEmbed("")
				embed.Image("https://nekos.brussell.me/image/" + img)
				_, er := ctx.SayEmbed(embed)
				if er != nil {
					_, e := ctx.Say("error sending embed, %s", er)
					if e != nil {
						return
					}
					return
				}
			} else {
				_, e := ctx.Say("Something went wrong while picking the catgirl site")
				if e != nil {
					return
				}
				return
			}
		}
	}
}
