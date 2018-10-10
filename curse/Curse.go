package curse

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/modmuss50/CAV2"
	"sort"
	"strings"
	"time"
)

func Load() {
	//logs with with cav
	err := cav2.SetupDefaultConfig()
	if err != nil {
		panic(err)
	}
}

func HandleCurseMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if strings.HasPrefix(m.Content, "!curse") {
		messageSplit := strings.Split(m.Content, " ")

		if len(messageSplit) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Incorrect usage; `!curse <username>`")
			return true
		}
		username := messageSplit[1]

		start := time.Now()

		database, err := cav2.Search(username)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "There was an error searching for projects on curse")
			return true
		}

		sort.Sort(SortAddon(database))

		var downloads float64 = 0
		for _, addon := range database {
			for _, author := range addon.Authors {
				if strings.EqualFold(author.Name, username) {
					downloads += addon.DownloadCount
				}
			}
		}

		if downloads == 0 {
			s.ChannelMessageSend(m.ChannelID, "No downloads found for "+username)
			return true
		}

		count := 0
		projects := ""
		for _, addon := range database {
			for _, author := range addon.Authors {
				if strings.EqualFold(author.Name, username) {
					projects = projects + addon.Name + " : `" + humanize.Comma(round(addon.DownloadCount)) + "`\n"
					count++
				}
			}

		}

		s.ChannelMessageSend(m.ChannelID, projects)
		since := time.Since(start)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has `%s` total downloads over `%d` projects . Loaded in %s", username, humanize.Comma(round(downloads)), count, since))
	}
	return false
}

type SortAddon []cav2.Addon

func (c SortAddon) Len() int           { return len(c) }
func (c SortAddon) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortAddon) Less(i, j int) bool { return c[i].DownloadCount > c[j].DownloadCount }

func round(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}
