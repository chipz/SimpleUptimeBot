package main

import (
	"time"
	"fmt"
	"net/http"
	"os"
	"log"
	"reflect"
)

// start it as:
// SimpleUptimeBot "token:forBot"
// or
// BOT_TOKEN="token:forBot" SimpleUptimeBot
//

var TestWebsites = []Website{
	//Website{Url: `http://www.google.com/robots.txt`, Interval: 5, ChatId: 0},
	//Website{Url: `https://news.ycombinator.com/item?id=13816627`, Interval: 10, ChatId: 0},
}

var c1 chan Website = make(chan Website)
var killingPills chan bool = make(chan bool)
var webs chan []Website = make(chan []Website)

func main() {
	Websites = TestWebsites
	go StartBot()
	go MonitorListOfWebsite()
	go MonitorWebsitesChannel()

	select {

	}
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

func MonitorListOfWebsite(){
	for {
		lw := <- webs
		log.Println("got list web")
		if(!reflect.DeepEqual(TempWebsites, lw)){
			for _, website := range lw {
				time.Sleep(10 * time.Millisecond)
				go func(website Website) {
					log.Println("sending killing pills..")
					killingPills <- true
				}(website)
			}
			TempWebsites = lw
			for _, website := range lw {
				time.Sleep(600 * time.Millisecond)
				go func(website Website) {
					log.Println("sending monitoring task..")
					go MonitorWebsitesChannel()
					c1 <- website
				}(website)
			}
		}
	}
}

func MonitorWebsitesChannel() {
	for {
		select {
			case <- killingPills:
				log.Println("Killing me softly..")
				return
			case website := <- c1:
				log.Println("Got monitoring task..")
				MonitorWebsite(website)
				log.Printf("waiting to check again for: %v second(s)", website.Interval)
				time.Sleep(time.Duration(website.Interval) * time.Second)
				log.Println("done waiting")
				go func(website Website) {
					c1 <- website
				}(website)
		}
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
}

func GetStatusCode(url string) (int, error) {
	res, err := http.Get(url)
	fmt.Println("url: " + url + "--> " + res.Status)
	return res.StatusCode, err
}

