package main

import (
	"fmt"
	"github.com/TechReborn/DiscordBot/minecraft"
	"github.com/bwmarrin/discordgo"
	"github.com/modmuss50/goutils"
	"strings"
	"time"
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
			lat, err := minecraft.GetLatest()
			if err != nil {
				fmt.Println(err)
				return
			}

			var latest = lat
			DiscordClient.UpdateStatus(0, latest.Snapshot)
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
	if goutils.FileExists("discord_muted.txt") {
		for _, str := range goutils.ReadLinesFromFile("discord_muted.txt") {
			if str == m.GuildID {
				return
			}
		}
	}
	channel, _ := DiscordClient.Channel(m.ChannelID)
	channelName := channel.Name
	if channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
		channelName = m.Author.Username
	}
	fmt.Println("#" + channelName + " <" + m.Author.Username + ">:" + m.Content)

	if m.Author.ID == BotID {
		return
	}

	if m.Content == "!version" {
		version, err := minecraft.GetLatest()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "An error occurred.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Latest snapshot: "+version.Snapshot)
			s.ChannelMessageSend(m.ChannelID, "Latest release: "+version.Release)
		}

	}

	if m.Content == "!verNotify" {
		if !isAuthorAdmin(m.Author) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permission to run that command.")
			return
		}
		goutils.AppendStringToFile(m.ChannelID, "channels.txt")
		s.ChannelMessageSend(m.ChannelID, "The bot will now announce new minecraft versions here!")
	}

	if m.Content == "!mute" {
		if !isAuthorAdmin(m.Author) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permission to run that command.")
			return
		}
		goutils.AppendStringToFile(m.GuildID, "discord_muted.txt")
		s.ChannelMessageSend(m.ChannelID, "The bot will no longer response to commands in this server!")
	}

	value, handled := handleTempMessage(m.Content)
	if handled {
		s.ChannelMessageSend(m.ChannelID, value)
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
