package main

import (
	"fmt"
	"github.com/modmuss50/discordBot/minecraft"
)

func main() {

	var latest = minecraft.GetLatest()
	fmt.Println("Latest snapshot: " + latest.Snapshot)
	fmt.Println("Latest release: " + latest.Release)
}



