package minecraft

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"unsafe"
)

//GetVersionJson returns the json file that contains all the minecraft versions
func GetVersionJson() string {
	return getJson("https://launchermeta.mojang.com/mc/game/version_manifest.json")
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string) string {
	r, err := myClient.Get(url)
	if err != nil {
		return "error"
	}
	defer r.Body.Close()

	var buf = new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	var b = buf.Bytes()
	var s = *(*string)(unsafe.Pointer(&b))
	return s
}

//GetLatest returns an instance of the latest minecraft version information
func GetLatest() Latest {
	res := &Json{}
	str := GetVersionJson()
	err := json.Unmarshal([]byte(str), res)
	if err != nil {
		log.Fatal(err)
	}
	return res.Latest
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
