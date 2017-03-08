package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/modmuss50/discordBot/fileutil"
	"github.com/modmuss50/discordBot/minecraft"
	"time"
)

var (
	Token         string //The token of the bot user
	BotID         string //The id of the bot
	FirstCheck    bool //If the application has not done its first check for a new version
	Connected     bool //If the discord bot is has connected
	LastLatest    string //The latest release version of minecraft
	LastSnapshot  string //The latest snapshot of minecraft
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

			} else {
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
	fmt.Println(m.Author.Username + ":" + m.Content)
	if m.Author.ID == BotID {
		return
	}
	if m.Content == "!version" {
		var latest = minecraft.GetLatest()
		s.ChannelMessageSend(m.ChannelID, "Latest snapshot: "+latest.Snapshot)
		s.ChannelMessageSend(m.ChannelID, "Latest release: "+latest.Release)
	}
	if m.Content == "!verNotify" {
		fileutil.AppendStringToFile(m.ChannelID, "channels.txt")
		s.ChannelMessageSend(m.ChannelID, "The bot will now annouce new minecraft versions here!")
	}

	if m.Content == "!commands" {
		s.ChannelMessageSend(m.ChannelID, "The following commands are available for you to use. `!version`, `!issue`, `!wiki`, `!jei`")
	}

	if m.Content == "!issuse" || m.Content == "!issue" {
		s.ChannelMessageSend(m.ChannelID, "You can report an bug on our issuse tracker here: https://github.com/TechReborn/TechReborn/issues Please take a quick look to check that your isssus hasnt been reported before.")
	}
	if m.Content == "!wiki" {
		s.ChannelMessageSend(m.ChannelID, "We have a wiki located here: https://wiki.techreborn.ovh/ Please not not all the content is present at the current time.")
	}
	if m.Content == "!jei" {
		s.ChannelMessageSend(m.ChannelID, "JEI is a great mod to use to findout how to craft something, TechReborn has full support. You can download JEI from here: https://minecraft.curseforge.com/projects/just-enough-items-jei")
	}
}

//Loads the token from the file
func getToken() string {
	return fileutil.ReadStringFromFile("token.txt")
}
