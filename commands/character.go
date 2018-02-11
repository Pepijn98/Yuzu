package commands

import (
	"fmt"
	"strings"
	"yuzu"
	"yuzu/functions"

	"github.com/KurozeroPB/kitsu-go"
)

// Character command
type Character struct{}

// IsOwnerOnly f
func (Character) IsOwnerOnly() bool {
	return false
}

// Help f
func (Character) Help() [2]string {
	return [2]string{"description", "usage"}
}

// Process f
func (Character) Process(ctx yuzu.Context) {
	c := strings.Join(ctx.Args, " ")

	if c == "" {
		_, err := ctx.Say("Please tell me which character to search for.")
		if err != nil {
			return
		}
		return
	}

	char, err := kitsu.SearchCharacter(c)
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/character.go")
		_, e := ctx.Say("Something went wrong retrieving the character info:\n", err)
		if e != nil {
			return
		}
		return
	}

	desc := strings.Replace(char.Attributes.Description, "<br/>", "\n", -1)

	embed := yuzu.NewEmbed(char.Attributes.Name)
	embed.Description = desc
	embed.Thumbnail(char.Attributes.Image.Original)
	_, er := ctx.SayEmbed(embed)
	if er != nil {
		_, e := ctx.Say("Something went wrong when sending the message:\n", er)
		if e != nil {
			return
		}
		return
	}
}
