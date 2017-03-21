package main

import (
	"time"
	"fmt"
	"net/http"
	"os"
	"log"
	"sync"
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

func main() {
	Websites = TestWebsites
	go StartBot()
	for {
		if(!reflect.DeepEqual(TempWebsites, Websites)){
			go StartWebsiteMonitors()
			TempWebsites = Websites
		}
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

func StartWebsiteMonitors() {
	var wg sync.WaitGroup
	log.Printf("websites now: %s", Websites)
	for _, website := range Websites {
		if !WebsiteExist(website, runningWebsites){
			wg.Add(1)
			time.Sleep(10 * time.Millisecond)
			go func(website Website) {
				defer wg.Done()
				MonitorWebsiteLoop(website)
			}(website)
		}
	}
	defer wg.Wait()
}

func MonitorWebsiteLoop(website Website) {
	for WebsiteExist(website, Websites) && !WebsiteExist(website, runningWebsites) {
		runningWebsites = append(runningWebsites, website)
		MonitorWebsite(website)
		log.Printf("waiting to check again for: %v second(s)", website.Interval)
		time.Sleep(time.Duration(website.Interval) * time.Second)
		log.Println("done waiting")
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

func WebsiteExist(websiteToCheck Website, websites []Website) bool {
	for _, existingWebsite := range websites {
		if existingWebsite == websiteToCheck {
			return true
		}
	}
	return false
}

