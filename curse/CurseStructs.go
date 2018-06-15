package curse

import "encoding/json"

type AddonDatabase struct {
	TimeStamp int64 `json:"timestamp"`
	Addons []AddonDB `json:"data"`
}

type AddonDB struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Authors []Author `json:"authors"`
	DownloadCount json.Number `json:"downloadCount"`
}

type Author struct {
	Name string `json:"name"`
}