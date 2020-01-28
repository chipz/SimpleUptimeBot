package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
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

var webs chan []Website = make(chan []Website)

func main() {
	Websites = TestWebsites
	go StartBot()

	ticker := time.NewTicker(5 * time.Second)
	pchan := make(chan interface{})

	go func() {
		for t := range ticker.C {
			log.Println("Tick at", t)
			select {
			case web := <- pchan:
				log.Printf("Got list website %v", web.([]Website))
				Websites = web.([]Website)
			default:
				log.Printf("no activity")
				for _, element := range Websites {
					go func(website Website) {
						log.Printf("Got website %v", website.Url)
						MonitorWebsite(website)
					}(element)
				}
			}
		}
	}()

	//for {
	//	lw := <- webs
	//
	//	source := observable.Just(lw)
	//
	//	onNext := handlers.NextFunc(func(item interface{}) {
	//		if item, ok := item.([]Website); ok {
	//			pchan <- item
	//		}
	//	})
	//
	//	_ = source.Subscribe(observers.New(onNext))
	//}


	//for {
	//	lw := <- webs
	//
	//	flavio := make(chan string)
	//
	//
	//}
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

