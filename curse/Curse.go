package curse

import (
	"log"
	"fmt"
	"bytes"
	"github.com/modmuss50/goutils"
	"encoding/json"
	"strings"
	"github.com/dustin/go-humanize"
	"github.com/bwmarrin/discordgo"
	"sort"
	"github.com/patrickmn/go-cache"
	"time"
	"compress/bzip2"
	"github.com/modmuss50/CAV2"
	"strconv"
)

var (
	Cache *cache.Cache
)


func Load(){
	Cache = cache.New(90*time.Minute, 1*time.Minute)

	//logs with with cav
	err := cav2.SetupDefaultConfig()
	if err != nil {
		panic(err)
	}
}

func HandleCurseMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if strings.HasPrefix(m.Content, "!curse"){
		messageSplit := strings.Split(m.Content, " ")

		if len(messageSplit) != 2{
			s.ChannelMessageSend(m.ChannelID, "Incorrect usage; `!curse <username>`")
			return true
		}
		username := messageSplit[1]

		data, found := Cache.Get("data")
		if !found {
			m, merr := s.ChannelMessageSend(m.ChannelID, "Curse Cache Expired, give me a second ...")
			s.ChannelTyping(m.ChannelID)
			reader, dwerr := goutils.Download("http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2")
			if dwerr != nil{
				s.ChannelMessageSend(m.ChannelID, "failed to download curse data")
				return true
			}
			jsonStr := LoadFromBz2(reader)
			var database AddonDatabase
			err := json.Unmarshal([]byte(jsonStr), &database)
			if err != nil {
				log.Fatal(err)
				s.ChannelMessageSend(m.ChannelID, "failed to read json file")
				return true
			}

			data = database

			Cache.Set("data", database, cache.DefaultExpiration)

			if merr == nil { //Remove that message
				s.ChannelMessageDelete(m.ChannelID, m.ID)
			}
		}

		var database = data.(AddonDatabase)



		var addons []int

		for _,addon := range database.Addons {
			for _,author := range addon.Authors {
				if strings.EqualFold(author.Name, username){
					addons = append(addons, addon.Id)
				}
			}
		}

		if len(addons) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No addons found for " + username)
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
			s.ChannelMessageSend(m.ChannelID, "No downloads found for " + username)
			return true
		}

		sort.Sort(SortAddon(addonInfo))
		projects := ""
		for _,addon := range addonInfo {
			projects = projects + addon.Name + " : `" + humanize.Comma(round(addon.DownloadCount)) + "`\n"
		}

		s.ChannelMessageSend(m.ChannelID, projects)

		fmt.Println(humanize.Comma(round(downloads)))
		s.ChannelMessageSend(m.ChannelID, username +  " has `" + humanize.Comma(round(downloads)) + "` total downloads over `" + strconv.Itoa(len(addons)) + "` projects")
	}
	return false
}

type SortAddon []cav2.Addon

func (c SortAddon) Len() int           { return len(c) }
func (c SortAddon) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortAddon) Less(i, j int) bool { return c[i].DownloadCount  > c[j].DownloadCount }

func LoadFromBz2(byteArray []byte) string {
	reader := bzip2.NewReader(bytes.NewReader(byteArray))
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	s := buf.String()
	return s
}

func round(val float64) int64 {
	if val < 0 { return int64(val-0.5) }
	return int64(val+0.5)
}