package commands

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"time"
	"yuzu/config"

	"yuzu"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
)

var uptime = time.Now()

// Stats f
type Stats struct{}

// IsOwnerOnly f
func (Stats) IsOwnerOnly() bool {
	return false
}

// Help f
func (Stats) Help() [2]string {
	return [2]string{"Shows stats about the bot", ""}
}

// Process f
func (Stats) Process(ctx yuzu.Context) {
	var users, channels int
	for _, g := range ctx.State.Guilds {
		users += len(g.Members)
		channels += len(g.Channels)
	}
	counter := yuzu.CommandCounter.Counter
	var mostused string
	var nums []int
	for _, num := range counter {
		nums = append(nums, num)
	}
	sort.Ints(nums)
	highest := nums[len(nums)-1]
	for command, num := range counter {
		if num == highest {
			mostused = fmt.Sprintf("%s (%d)", command, num)
		}
	}
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	ram := runtime.MemStats{}
	runtime.ReadMemStats(&ram)

	embed := yuzu.NewEmbed("Stats about " + ctx.State.User.Username)
	embed.Field("Bot version", "v"+config.Config.Version, true)
	embed.Field("Go version", runtime.Version()[2:], true)
	embed.Field("Lib", "discordgo ("+discordgo.VERSION+")", true)
	embed.Field("Uptime", formattedTime(time.Since(uptime)), true)
	embed.Field("Running Tasks", strconv.Itoa(runtime.NumGoroutine()), true)
	embed.Field("Users | Channels | Guilds", strconv.Itoa(users)+" | "+strconv.Itoa(channels)+" | "+strconv.Itoa(len(ctx.State.Guilds)), true)
	embed.Field("Memory Usage", "total / heap / garbage\n"+humanize.Bytes(memStats.Sys)+" / "+humanize.Bytes(memStats.Alloc)+" / "+humanize.Bytes(memStats.TotalAlloc), true)
	embed.Field("Used RAM", humanize.Bytes(ram.Alloc)+"/"+humanize.Bytes(ram.Sys), true)
	embed.Field("Most used command", mostused, true)
	embed.Thumbnail(discordgo.EndpointUserAvatar(ctx.State.User.ID, ctx.State.User.Avatar))
	_, err := ctx.SayEmbed(embed)
	if err != nil {
		_, e := ctx.Say("Error sending embed, ", err)
		if e != nil {
			return
		}
		return
	}
}

func formattedTime(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}
