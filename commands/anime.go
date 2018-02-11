package commands

import (
	"fmt"
	"strconv"
	"yuzu"
	"yuzu/functions"

	"strings"

	"github.com/KurozeroPB/kitsu-go"
)

// Anime command
type Anime struct{}

// IsOwnerOnly f
func (Anime) IsOwnerOnly() bool {
	return false
}

// Help f
func (Anime) Help() [2]string {
	return [2]string{"Get info on an anime from kitsu.io", "<anime_name>"}
}

// Process f
func (Anime) Process(ctx yuzu.Context) {
	ani := strings.Join(ctx.Args, " ")

	if ani == "" {
		_, err := ctx.Say("Please tell me which anime to search for.")
		if err != nil {
			return
		}
		return
	}

	anime, err := kitsu.SearchAnime(ani, 0)
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/anime.go")
		_, e := ctx.Say("Something went wrong retrieving the anime info:\n", err)
		if e != nil {
			return
		}
		return
	}

	var (
		status = "-"
		start  = "?"
		end    = "?"
		ageR   = "-"
		ageRG  = "-"
	)

	if anime.Attributes.Status != "" {
		status = anime.Attributes.Status
	}
	if anime.Attributes.StartDate != "" {
		start = anime.Attributes.StartDate
	}
	if anime.Attributes.EndDate != "" {
		end = anime.Attributes.EndDate
	}
	if anime.Attributes.AgeRating != "" {
		ageR = anime.Attributes.AgeRating
	}
	if anime.Attributes.AgeRatingGuide != "" {
		ageRG = anime.Attributes.AgeRatingGuide
	}

	embed := yuzu.NewEmbed(anime.Attributes.Titles.JaJp)
	embed.Description = anime.Attributes.Synopsis
	embed.Field("Status", status, true)
	embed.Field("Episodes", strconv.Itoa(anime.Attributes.EpisodeCount), true)
	embed.Field("Favorites", strconv.Itoa(anime.Attributes.FavoritesCount), true)
	embed.Field("Rating", strconv.Itoa(anime.Attributes.RatingRank), true)
	embed.Field("Popularity", strconv.Itoa(anime.Attributes.PopularityRank), true)
	embed.Field("Age Rating", ageR, true)
	embed.Field("Age Rating Guide", ageRG, true)
	embed.Field("Start/End", start+" until "+end, false)
	embed.Thumbnail(anime.Attributes.PosterImage.Original)
	_, er := ctx.SayEmbed(embed)
	if er != nil {
		_, e := ctx.Say("Something went wrong when sending the message:\n", er)
		if e != nil {
			return
		}
		return
	}
}
