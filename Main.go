package main

import (
	"fmt"
	"github.com/TechReborn/DiscordBot/file"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

var (
	Token         string             //The token of the bot user
	BotID         string             //The id of the bot
	DiscordClient *discordgo.Session //The discord client
)

func main() {

	err := populateInitialJiraVersions()
	if err != nil {
		fmt.Println("Failed to load jira versions")
		panic(err)
	}

	err = populateInitialGameVersions()
	if err != nil {
		fmt.Println("Failed to load jira versions")
		panic(err)
	}

	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for range ticker.C {
			updateCheck()
		}
	}()

	err = LoadDiscord()
	if err != nil {
		fmt.Println("Failed to load discord")
		panic(err)
	}
}

func updateCheck() {
	go jiraUpdateCheck(postJiraMessage)
	go gameUpdateCheck(postGameMessage)
}

func postGameMessage(message string) error {
	fmt.Println(message)

	lines, err := file.ReadLines("channels.txt")

	if err != nil {
		return err
	}

	for _, element := range lines {
		DiscordClient.ChannelMessageSend(element, message)
	}

	return nil
}

func postJiraMessage(message string) error {
	fmt.Println(message)

	lines, err := file.ReadLines("jira_channels.txt")

	if err != nil {
		return err
	}

	for _, element := range lines {
		DiscordClient.ChannelMessageSend(element, message)
	}

	return nil
}

//LoadDiscord is based a lot off https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
func LoadDiscord() error {
	t, err := getToken()
	if err != nil {
		return err
	}

	Token = t

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		return err
	}

	DiscordClient = dg

	u, err := dg.User("@me")
	if err != nil {
		return err
	}

	BotID = u.ID

	dg.AddHandler(handleMessage)

	err = dg.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})
	return nil
}

//Called when a message is posted
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := DiscordClient.Channel(m.ChannelID)
	channelName := channel.Name
	if channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM {
		channelName = m.Author.Username
	}
	fmt.Println("#" + channelName + " <" + m.Author.Username + ">:" + m.Content)

	if m.Author.ID == BotID {
		return
	}

	if m.Content == "!gameNotify" {
		if isAuthorAdmin(m.Author) {
			err := file.AppendString(m.ChannelID, "channels.txt")
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, "The bot will now announce new minecraft versions here!")
			} else {
				log.Println("Failed to write game channels", err)
				s.ChannelMessageSend(m.ChannelID, "An error occurred, contact bot owner")
			}
		}
	}

	if m.Content == "!jiraNotify" {
		if isAuthorAdmin(m.Author) {
			err := file.AppendString(m.ChannelID, "jira_channels.txt")
			if err == nil {
				s.ChannelMessageSend(m.ChannelID, "The bot will now announce new jira versions here!")
			} else {
				log.Println("Failed to write jira channels", err)
				s.ChannelMessageSend(m.ChannelID, "An error occurred, contact bot owner")
			}
		}
	}

	if file.Exists("discord_muted.txt") {
		lines, err := file.ReadLines("discord_muted.txt")
		if err == nil {
			for _, str := range lines {
				if str == m.GuildID {
					return
				}
			}
		}
	}

	if m.Content == "!mute" {
		if !isAuthorAdmin(m.Author) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permission to run that command.")
			return
		}
		err := file.AppendString(m.GuildID, "discord_muted.txt")
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, "The bot will no longer response to commands in this server!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "An error occurred, contact bot owner")
		}
	}

	if m.Content == "!commands" || m.Content == "!help" {
		cmdList := ""
		lines, err := file.ReadLines("commands.txt")
		if err == nil {
			for _, element := range lines {
				command := "!" + strings.Split(element, "=")[0]
				cmdList = cmdList + "`" + command + "` "
			}
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

		err := file.AppendString(textLine, "commands.txt")
		if err == nil {
			s.ChannelMessageSend(m.ChannelID, "The command has been added!")
		} else {
			s.ChannelMessageSend(m.ChannelID, "An error occurred, contact bot owner")
		}

	}

	if file.Exists("commands.txt") {
		lines, err := file.ReadLines("commands.txt")
		if err == nil {
			for _, element := range lines {
				command := "!" + strings.Split(element, "=")[0]
				reply := strings.Split(element, "=")[1]
				if m.Content == command {
					s.ChannelMessageSend(m.ChannelID, reply)
				}
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
func getToken() (string, error) {
	return file.ReadString("token.txt")
}
