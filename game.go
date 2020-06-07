package main

import (
	"encoding/json"
	"fmt"
	"time"
)

var gameVersionsCache []string

//Set the initial game versions in the cache
func populateInitialGameVersions() error {
	versions, err := getGameVersions()

	if err != nil {
		return err
	}

	for _, version := range versions {
		gameVersionsCache = append(gameVersionsCache, version.ID)
	}

	fmt.Printf("Loaded %d initial game versions\n", len(gameVersionsCache))

	return nil
}

func gameUpdateCheck(callback func(message string) error) error {
	currentVersions, err := getGameVersions()

	if err != nil {
		return err
	}

	for _, version := range currentVersions {
		if !contains(gameVersionsCache, version.ID) {
			gameVersionsCache = append(gameVersionsCache, version.ID)

			err = callback(gameVersionAsString(version))

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func gameVersionAsString(version Version) string {
	return fmt.Sprintf("A new %s version of minecraft was just released! : %s", version.Type, version.ID)
}

func getGameVersions() ([]Version, error) {
	jsonStr, err := getJson("https://launchermeta.mojang.com/mc/game/version_manifest.json")

	if err != nil {
		return nil, err
	}

	var meta GameMeta
	err = json.Unmarshal([]byte(jsonStr), &meta)
	return meta.Versions, err
}

type GameMeta struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	URL         string    `json:"url"`
	Time        time.Time `json:"time"`
	ReleaseTime time.Time `json:"releaseTime"`
}
