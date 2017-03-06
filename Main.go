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
)

func main() {
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for range ticker.C {
			//TODO check to see if the version has changed
		}
	}()
	LoadDiscord()
}

//Based a lot off https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go
func LoadDiscord() {

	Token = getToken()
	dg, err := discordgo.New("Bot " + Token)
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
}

//Loads the token from the file
func getToken() string {
	return fileutil.ReadStringFromFile("token.txt")
}



