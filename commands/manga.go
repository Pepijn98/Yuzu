package commands

import (
	"fmt"
	"strconv"
	"yuzu"
	"yuzu/functions"

	"strings"

	"github.com/KurozeroPB/kitsu-go"
)

// Manga command
type Manga struct{}

// IsOwnerOnly f
func (Manga) IsOwnerOnly() bool {
	return false
}

// Help f
func (Manga) Help() [2]string {
	return [2]string{"Get info on an manga from kitsu.io", "<manga_name>"}
}

// Process f
func (Manga) Process(ctx yuzu.Context) {
	m := strings.Join(ctx.Args, " ")

	if m == "" {
		_, err := ctx.Say("Please tell me which manga to search for.")
		if err != nil {
			return
		}
		return
	}

	manga, err := kitsu.SearchManga(m, 0)
	if err != nil {
		functions.ReportError(ctx.Session, fmt.Sprintf("%s", err), "/commands/manga.go")
		_, e := ctx.Say("Something went wrong retrieving the manga info:\n", err)
		if e != nil {
			return
		}
		return
	}

	var (
		start = "?"
		end   = "?"
		ageR  = "-"
		ageRG = "-"
	)

	if manga.Attributes.StartDate != "" {
		start = manga.Attributes.StartDate
	}
	if manga.Attributes.EndDate != "" {
		end = manga.Attributes.EndDate
	}
	if manga.Attributes.AgeRating != "" {
		ageR = manga.Attributes.AgeRating
	}
	if manga.Attributes.AgeRatingGuide != "" {
		ageRG = manga.Attributes.AgeRatingGuide
	}

	embed := yuzu.NewEmbed(manga.Attributes.Titles.JaJp)
	embed.Description = manga.Attributes.Synopsis
	embed.Field("Volumes", strconv.Itoa(manga.Attributes.VolumeCount), true)
	embed.Field("Chapters", strconv.Itoa(manga.Attributes.ChapterCount), true)
	embed.Field("Favorites", strconv.Itoa(manga.Attributes.FavoritesCount), true)
	embed.Field("Rating", strconv.Itoa(manga.Attributes.RatingRank), true)
	embed.Field("Popularity", strconv.Itoa(manga.Attributes.PopularityRank), true)
	embed.Field("Age Rating", ageR, true)
	embed.Field("Age Rating Guide", ageRG, true)
	embed.Field("Start/End", start+" until "+end, false)
	embed.Thumbnail(manga.Attributes.PosterImage.Original)
	_, er := ctx.SayEmbed(embed)
	if er != nil {
		_, e := ctx.Say("Something went wrong when sending the message:\n", er)
		if e != nil {
			return
		}
		return
	}
}
