package main

import (
	"fmt"
	"github.com/modmuss50/discordBot/minecraft"
	"github.com/modmuss50/discordBot/fileutil"
	"github.com/bwmarrin/discordgo"
	"time"
)

var (
	Token string
	BotID string
	FirstCheck bool
	Connected bool
	LastLatest string
	LastSnapshot string
	DiscordClient *discordgo.Session
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
				for _,element := range fileutil.ReadLinesFromFile("channels.txt") {
					if latest.Release != LastLatest{
						DiscordClient.ChannelMessageSend(element, "A new release version of minecraft was just released! : " + latest.Release)
					}
					if latest.Snapshot != LastSnapshot{
						DiscordClient.ChannelMessageSend(element, "A new snapshot version of minecraft was just released! : " + latest.Snapshot)
					}
				}

				LastLatest = latest.Release
				LastSnapshot = latest.Snapshot
			}
		}
	}()


	LoadDiscord()
}

//Based a lot off https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
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
		_, _ = s.ChannelMessageSend(m.ChannelID, "Latest snapshot: " + latest.Snapshot)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Latest release: " + latest.Release)
	}
	if m.Content == "!verNotify" {
		fileutil.AppendStringToFile(m.ChannelID, "channels.txt")
		_, _ = s.ChannelMessageSend(m.ChannelID, "The bot will now annouce new minecraft versions here!")
	}

	if m.Content == "!commands" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "The following commands are advalibe for you to use. `!version`, `!issue`, `!wiki`, `!jei`")
	}

	if m.Content == "!issuse" || m.Content == "!issue" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "You can report an bug on our issuse tracker here: https://github.com/TechReborn/TechReborn/issues Please take a quick look to check that your isssus hasnt been reported before.")
	}
	if m.Content == "!wiki" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "We have a wiki located here: https://wiki.techreborn.ovh/ Please not not all the content is present at the current time.")
	}
	if m.Content == "!jei" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "JEI is a great mod to use to findout how to craft something, TechReborn has full support. You can download JEI from here: https://minecraft.curseforge.com/projects/just-enough-items-jei")
	}
}

//Loads the token from the file
func getToken() string {
	return fileutil.ReadStringFromFile("token.txt")
}



