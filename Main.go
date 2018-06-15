package main

import (
	"fmt"
	"strings"
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/TechReborn/DiscordBot/curse"
	"github.com/modmuss50/MCP-Diff/mcpDiff"
	"strconv"
	"github.com/modmuss50/MCP-Diff/utils"
	"net/url"
	"net/http"
	"bytes"
	"io/ioutil"
	"github.com/modmuss50/goutils"
	"github.com/TechReborn/DiscordBot/minecraft"
)

var (
	Token         string             //The token of the bot user
	BotID         string             //The id of the bot
	FirstCheck    bool               //If the application has not done its first check for a new version
	Connected     bool               //If the discord bot is has connected
	LastLatest    string             //The latest release version of minecraft
	LastSnapshot  string             //The latest snapshot of minecraft
	DiscordClient *discordgo.Session //The discord client
)

func main() {

	curse.Load()

	FirstCheck = true

	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for range ticker.C {
			if !Connected {
				return
			}
			var latest = minecraft.GetLatest()
			if FirstCheck == true {
				LastLatest = latest.Release
				LastSnapshot = latest.Snapshot
				FirstCheck = false

			} else if goutils.FileExists("channels.txt") {
				for _, element := range goutils.ReadLinesFromFile("channels.txt") {
					if latest.Release != LastLatest {
						DiscordClient.ChannelMessageSend(element, "A new release version of minecraft was just released! : "+latest.Release)
					}
					if latest.Snapshot != LastSnapshot {
						DiscordClient.ChannelMessageSend(element, "A new snapshot version of minecraft was just released! : "+latest.Snapshot)
					}
				}

				LastLatest = latest.Release
				LastSnapshot = latest.Snapshot
			}
		}
	}()

	LoadDiscord()
}

//LoadDiscord is based a lot off https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
func LoadDiscord() {

	Token = getToken()
	dg, err := discordgo.New("Bot " + Token)
	DiscordClient = dg
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}
	BotID = u.ID

	dg.AddHandler(handleMessage)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	Connected = true
	<-make(chan struct{})
	return
}

//Called when a message is posted
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel,_ := DiscordClient.Channel(m.ChannelID)
	channelName := channel.Name
	if channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
		channelName = m.Author.Username
	}
	fmt.Println("#" + channelName + " <" + m.Author.Username + ">:" + m.Content)
	utils.AppendStringToFile("#" + channelName + " <" + m.Author.Username + ">:" + m.Content, "discordLog.txt")
	if m.Author.ID == BotID {
		return
	}
	if m.Content == "!version" {
		var latest = minecraft.GetLatest()
		s.ChannelMessageSend(m.ChannelID, "Latest snapshot: "+latest.Snapshot)
		s.ChannelMessageSend(m.ChannelID, "Latest release: "+latest.Release)
	}

	if m.Content == "!verNotify" {
		if !isAuthorAdmin(m.Author) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permission to run that command.")
			return
		}
		goutils.AppendStringToFile(m.ChannelID, "channels.txt")
		s.ChannelMessageSend(m.ChannelID, "The bot will now announce new minecraft versions here!")
	}

	value, handled := handleTempMessage(m.Content)
	if handled {
		s.ChannelMessageSend(m.ChannelID, value)
	}

	if curse.HandleCurseMessage(s, m) {
		return
	}

	if m.Content == "!commands" || m.Content == "!help" {
		cmdList := ""
		for _, element := range goutils.ReadLinesFromFile("commands.txt") {
			command := "!" + strings.Split(element, "=")[0]
			cmdList = cmdList + "`" + command + "` "
		}
		if isAuthorAdmin(m.Author) {
			cmdList = cmdList + "`!addCom` "
			cmdList = cmdList + "`!verNotify` "
		}
		cmdList = cmdList + "`!mcpDiff <old> <new> (e.g: 20170614-1.12)` "
		cmdList = cmdList + "`!gm <srg>"
		cmdList = cmdList + "`!gf <srg>"
		cmdList = cmdList + "`!curse <username>"
		s.ChannelMessageSend(m.ChannelID, "The following commands are available for you to use. "+cmdList)
	}

	if m.Content == "!myID" {
		s.ChannelMessageSend(m.ChannelID, "You ID: `"+m.Author.ID+"`")
	}

	if strings.HasPrefix(m.Content, "!addCom") {
		if !isAuthorAdmin(m.Author) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permission to run that command.")
			return
		}
		text := strings.Replace(m.Content, "!addCom ", "", -1)
		textLine := strings.Replace(text, " ", "=", 1)
		goutils.AppendStringToFile(textLine, "commands.txt")
		s.ChannelMessageSend(m.ChannelID, "The command has been added!")
	}

	if strings.HasPrefix(m.Content, "!mcpDiff") {
		text := strings.Replace(m.Content, "!mcpDiff ", "", -1)
		split := strings.Split(text, " ")
		if len(split) != 2{
			s.ChannelMessageSend(m.ChannelID, "Usage: !mcpDiff <old> <new> `e.g: !mcpDiff 20170601-1.11 20170614-1.12` or `!mcpDiff stable-29-1.10.2 stable-32-1.11`you can find the list of MCP exports here: http://export.mcpbot.bspk.rs/")
			return
		}
		response, info, err := mcpDiff.GetMCPDiff(split[0], split[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, info)
		lines := strings.Split(response, "\n")
		if len(lines) -1 == 0 {
			s.ChannelMessageSend(m.ChannelID, "No changes in mappings between " + split[0] + " and " + split[1])
		} else {
			s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(lines) -1) + " changes in mappings, you can view them here: " + createPaste(response, info))
		}

	}

	if strings.HasPrefix(m.Content, "!gm") {
		text := strings.Replace(m.Content, "!gm ", "", -1)
		s.ChannelMessageSend(m.ChannelID, mcpDiff.LookupMethod(text))
	}

	if strings.HasPrefix(m.Content, "!gf") {
		text := strings.Replace(m.Content, "!gf ", "", -1)
		s.ChannelMessageSend(m.ChannelID, mcpDiff.LookupField(text))
	}

	if goutils.FileExists("commands.txt") {
		for _, element := range goutils.ReadLinesFromFile("commands.txt") {
			command := "!" + strings.Split(element, "=")[0]
			reply := strings.Split(element, "=")[1]
			if m.Content == command {
				s.ChannelMessageSend(m.ChannelID, reply)
			}
		}
	}

}

func createPaste(text string, title string) string {
	apiUrl := "https://paste.modmuss50.me"
	resource := "/api/create"
	data := url.Values{}
	data.Set("text", text)
	data.Set("title", title)
	data.Set("private", "1")
	data.Set("expire", "0")
	data.Set("name", "TechReborn Discord Bot")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Body)

	if resp.StatusCode == 200 { // OK
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return bodyString
	}
	return "An error occurred when getting paste bin"
}

func isAuthorAdmin(user *discordgo.User) bool {
	if user.ID != "98473211301212160" { //TODO have a file or some better way to do this.
		return false
	}
	return true
}

//Loads the token from the file
func getToken() string {
	return goutils.ReadStringFromFile("token.txt")
}
