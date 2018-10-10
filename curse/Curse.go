package curse

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/modmuss50/CAV2"
	"sort"
	"strconv"
	"strings"
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

		database, err := cav2.Search(username)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "There was an error searching for projects on curse")
			return true
		}

		var addons []int
		for _, addon := range database {
			for _, author := range addon.Authors {
				if strings.EqualFold(author.Name, username) {
					addons = append(addons, addon.ID)
				}
			}
		}

		if len(addons) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No addons found for "+username)
			return true
		}

		var downloads float64 = 0
		addonInfo, err := cav2.GetAddons(addons)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something bad happened when loading detailed addon data from curse")
			return true
		}

		for _, addon := range addonInfo {
			downloads += addon.DownloadCount
		}

		if downloads == 0 {
			s.ChannelMessageSend(m.ChannelID, "No downloads found for "+username)
			return true
		}

		sort.Sort(SortAddon(addonInfo))
		projects := ""
		for _, addon := range addonInfo {
			projects = projects + addon.Name + " : `" + humanize.Comma(round(addon.DownloadCount)) + "`\n"
		}

		s.ChannelMessageSend(m.ChannelID, projects)

		fmt.Println(humanize.Comma(round(downloads)))
		s.ChannelMessageSend(m.ChannelID, username+" has `"+humanize.Comma(round(downloads))+"` total downloads over `"+strconv.Itoa(len(addons))+"` projects")
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
