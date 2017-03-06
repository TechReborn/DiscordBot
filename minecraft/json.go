package minecraft

import (
	"net/http"
	"time"
	"bytes"
	"unsafe"
	"encoding/json"
	"log"
)

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

func GetLatest() Latest {
	res := &Json{}
	str := GetVersionJson()
	err := json.Unmarshal([]byte(str), res)
	if(err!=nil) {
		log.Fatal(err)
	}
	return res.Latest
}

type Latest struct {
	Snapshot string `json:"snapshot"`
	Release string `json:"release"`
}

type Json struct {
	Latest Latest `json:"latest"`
}
