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

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for range ticker.C {
			if !Connected {
				return
			}
			fmt.Println("Test1")
			var latest = minecraft.GetLatest()
			if FirstCheck == true {
				LastLatest = latest.Release
				LastSnapshot = latest.Snapshot
				FirstCheck = false

			} else {
				fmt.Println("Test2")
				for _,element := range fileutil.ReadLinesFromFile("channels.txt") {
					fmt.Println("Test3")
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
}

//Loads the token from the file
func getToken() string {
	return fileutil.ReadStringFromFile("token.txt")
}



