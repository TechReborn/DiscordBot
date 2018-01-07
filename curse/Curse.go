package curse

import (
	"log"
	"os"
	"fmt"
	"compress/bzip2"
	"bytes"
	"github.com/modmuss50/goutils"
	"encoding/json"
	"strings"
	"strconv"
	"github.com/dustin/go-humanize"
	"github.com/bwmarrin/discordgo"
	"sort"
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	Cache *cache.Cache
)


func Load(){
	Cache = cache.New(90*time.Minute, 1*time.Minute)
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
			s.ChannelMessageSend(m.ChannelID, "Curse Cache Expired, re-downloading")
			goutils.DownloadURL("http://clientupdate-v6.cursecdn.com/feed/addons/432/v10/complete.json.bz2", "curse.json.bz2")
			jsonStr := LoadFromBz2("curse.json.bz2")
			var database AddonDatabase
			err := json.Unmarshal([]byte(jsonStr), &database)
			if err != nil {
				log.Fatal(err)
				s.ChannelMessageSend(m.ChannelID, "failed to read json file")
			}

			data = database

			Cache.Set("data", database, cache.DefaultExpiration)
		}

		var database = data.(AddonDatabase)

		var downloads int64 = 0

		var addons []Addon

		for _,addon := range database.Addons {
			for _,author := range addon.Authors {
				if strings.EqualFold(author.Name, username){
					addons = append(addons, addon)

					downloads = downloads + getDownloadCount(addon.DownloadCount)
				}
			}
		}

		if downloads == 0 {
			s.ChannelMessageSend(m.ChannelID, "No downloads found for " + username)
			return true
		}

		sort.Sort(SortAddon(addons))
		projects := ""
		for i,addon := range addons {
			if i < 10 {
				projects = projects + addon.Name + " : `" + humanize.Comma(getDownloadCount(addon.DownloadCount)) + "`\n"

			}
		}
		s.ChannelMessageSend(m.ChannelID, projects)

		fmt.Println(humanize.Comma(downloads))
		s.ChannelMessageSend(m.ChannelID, username +  " has `" + humanize.Comma(downloads) + "` total downloads")
	}
	return false
}

type SortAddon []Addon

func (c SortAddon) Len() int           { return len(c) }
func (c SortAddon) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortAddon) Less(i, j int) bool { return getDownloadCount(c[i].DownloadCount ) > getDownloadCount(c[j].DownloadCount ) }

func getDownloadCount(number json.Number) int64 {
	intstr := strings.Split(number.String(), ".")[0]
	int, _ := strconv.ParseInt(intstr, 10, 64)
	return int
}

func LoadFromBz2(filename string) string {
	bzip_file, err := os.Open(filename)
	defer bzip_file.Close()
	if err != nil {
		log.Panic(err)
	}
	reader := bzip2.NewReader(bzip_file)
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	s := buf.String()
	bzip_file.Close()
	return s

}