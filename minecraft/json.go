package minecraft

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"unsafe"
)

//GetVersionJson returns the json file that contains all the minecraft versions
func GetVersionJson() (string, error) {
	return getJson("https://launchermeta.mojang.com/mc/game/version_manifest.json")
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string) (string, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	var buf = new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var b = buf.Bytes()
	var s = *(*string)(unsafe.Pointer(&b))
	return s, nil
}

//GetLatest returns an instance of the latest minecraft version information
func GetLatest() (Latest, error) {
	res := &Json{}
	str, err := GetVersionJson()
	if err != nil {
		return Latest{}, err
	}

	err = json.Unmarshal([]byte(str), res)
	if err != nil {
		return Latest{}, err
	}

	return res.Latest, nil
}

//Latest contains the lastest versions of minecraft
type Latest struct {
	Snapshot string `json:"snapshot"`
	Release  string `json:"release"`
}

//Json is the main json object
type Json struct {
	Latest Latest `json:"latest"`
}
