package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/modmuss50/discordBot/fileutil"
	"github.com/modmuss50/discordBot/minecraft"
	"github.com/modmuss50/MCP-Diff/mcpDiff"
	"strconv"
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

			} else if fileutil.FileExists("channels.txt") {
				for _, element := range fileutil.ReadLinesFromFile("channels.txt") {
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
	fmt.Println("#" + channel.Name + " <" + m.Author.Username + ">:" + m.Content)
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
		fileutil.AppendStringToFile(m.ChannelID, "channels.txt")
		s.ChannelMessageSend(m.ChannelID, "The bot will now annouce new minecraft versions here!")
	}

	if m.Content == "!commands" || m.Content == "!help" {
		cmdList := ""
		for _, element := range fileutil.ReadLinesFromFile("commands.txt") {
			command := "!" + strings.Split(element, "=")[0]
			cmdList = cmdList + "`" + command + "` "
		}
		if isAuthorAdmin(m.Author) {
			cmdList = cmdList + "`!addCom` "
			cmdList = cmdList + "`!verNotify` "
		}
		cmdList = cmdList + "`!mcpDiff <old> <new> (e.g: 20170614-1.12)` "
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
		fileutil.AppendStringToFile(textLine, "commands.txt")
		s.ChannelMessageSend(m.ChannelID, "The command has been added!")
	}

	if strings.HasPrefix(m.Content, "!mcpDiff") {
		text := strings.Replace(m.Content, "!mcpDiff ", "", -1)
		split := strings.Split(text, " ")
		response := mcpDiff.GetMCPDiff(split[0], split[1])
		lines := strings.Split(response, "\n")
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa(len(lines)) + " changes in mappings")
		if len(lines) > 100 {
			s.ChannelMessageSend(m.ChannelID, "There are over 100 changes, bot will not message you them all. A fix is coming soon")
			return 
		}
		if channel.IsPrivate {
			if len(lines) == 0{
				s.ChannelMessageSend(m.ChannelID, "No mappings changed")
			} else {
				for i,line := range lines{
					s.ChannelMessageSend(m.ChannelID, "```#" + strconv.Itoa(i) + "	" + line + "```")
				}
			}
		} else {
			if len(lines) == 0{
				s.ChannelMessageSend(m.ChannelID, "No mappings changed")
			} else if len(lines) > 20 {
				s.ChannelMessageSend(m.ChannelID, "Mappings changes are larger than 20 lines, please private message the bot the request")
			} else {
				for i,line := range lines{
					s.ChannelMessageSend(m.ChannelID, "```#" + strconv.Itoa(i) + "	" + line + "```")
				}
			}
		}
	}

	if fileutil.FileExists("commands.txt") {
		for _, element := range fileutil.ReadLinesFromFile("commands.txt") {
			command := "!" + strings.Split(element, "=")[0]
			reply := strings.Split(element, "=")[1]
			if m.Content == command {
				s.ChannelMessageSend(m.ChannelID, reply)
			}
		}
	}

}

func isAuthorAdmin(user *discordgo.User) bool {
	if user.ID != "98473211301212160" { //TODO have a file or some better way to do this.
		return false
	}
	return true
}

//Loads the token from the file
func getToken() string {
	return fileutil.ReadStringFromFile("token.txt")
}
