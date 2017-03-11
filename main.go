package main

import (
	"time"
	"fmt"
	"net/http"
	"os"
	"log"
)

// start it as:
// SimpleUptimeBot "token:forBot"
// or
// BOT_TOKEN="token:forBot" SimpleUptimeBot
//

var TestWebsites = []Website{
	Website{Url: `http://www.google.com/robots.txt`, Interval: 5, ChatId: 0},
}

func main() {
	Websites = TestWebsites
	go StartWebsiteMonitors()

	StartBot()
}

func GetToken() string {
	if len(os.Args) > 1 {
		log.Println("got token from command line arg")
		return os.Args[1]
	}
	v := os.Getenv("BOT_TOKEN")
	if v != "" {
		log.Println("got token from envvar")
		return v
	}
	log.Fatal("token not set. set it as commandline arg or in BOT_TOKEN envvar")
	return ""
}

func StartWebsiteMonitors() {
	for {
		log.Printf("websites now: %s", Websites)
		for _, website := range Websites {
			time.Sleep(10 * time.Millisecond)
			go MonitorWebsite(website)
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func MonitorWebsite(website Website) {
	code, err := GetStatusCode(website.Url)

	if code == 200 {
		log.Printf("url %s is okay", website.Url)
	} else {
		fmt.Println(err)
		fmt.Println("Sending failure notification about:\n" + website.Url)
		y := fmt.Sprintf("Got %s with status code (%d). Please check", website.Url, code)
		sendTelegramBotMessage(y, website.ChatId)
	}

	time.Sleep(time.Duration(website.Interval) * time.Second)
}

func GetStatusCode(url string) (int, error) {
	res, err := http.Get(url)
	fmt.Println("url: " + url + "--> " + res.Status)
	return res.StatusCode, err
}

