package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"net/url"
	"log"
	"os"
	"strings"
	"fmt"
	"github.com/satori/go.uuid"
	"strconv"
)

//Start the Telegram bot
func StartBot(){
	bot, err := tgbotapi.NewBotAPI(GetToken())
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for {
		for update := range updates {
			processor(update)
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

func sendTelegramBotMessage(message string, chatID int64) {
	if chatID != 0 {
		bot, err := tgbotapi.NewBotAPI(GetToken())
		if err != nil {
			log.Panic(err)
		}

		//bot.Debug = true
		bot.Debug = false

		log.Printf("Authorized on account %s", bot.Self.UserName)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		log.Printf("Have chat id %s", chatID)
		msg := tgbotapi.NewMessage(chatID, message)
		_, _ = bot.Send(msg)
	} else {
		log.Print("Chat id still 0")
	}
}

func processor(update tgbotapi.Update) {
	// Parse responses
	response := strings.Split(update.Message.Text, " ")

	if response[0] == "" {
		sendTelegramBotMessage("Please enter in a command:\nmonitor [url]", update.Message.Chat.ID)
		return
	}

	switch strings.ToLower(response[0]) {
	case "/monitor":
		_, err := url.ParseRequestURI(response[1])
		if err != nil {
			sendTelegramBotMessage("Invalid url: " + response[1], update.Message.Chat.ID)
		} else {
			var targetUrl = response[1]
			var u, _ = uuid.NewV4()
			var newWebsite = Website{Id: u, Url: targetUrl, Interval: 5, ChatId: update.Message.Chat.ID}
			Websites = append(Websites, newWebsite)

			webs <- Websites

			sendTelegramBotMessage("Added: " + newWebsite.Url, update.Message.Chat.ID)
		}
		break
	case "/remove":
		if (len(response) <= 1) {
			sendTelegramBotMessage("invalid remove command", update.Message.Chat.ID)
			break
		}
		var targetUrl = response[1]
		//var tobeRemovedWebsite = Website{Url: targetUrl, Interval: 5, ChatId: update.Message.Chat.ID}
		//Websites = remove(Websites, tobeRemovedWebsite)
		i, err := strconv.Atoi(targetUrl)
		if (err != nil) {
			sendTelegramBotMessage("invalid integer " + targetUrl, update.Message.Chat.ID)
			break
		}
		Websites = removeIndex(Websites, i)

		webs <- Websites

		log.Printf("removing %s", i)

		sendTelegramBotMessage(fmt.Sprintf("Removed: %v", i), update.Message.Chat.ID)
		break
	case "/list":
		urls := []string{}
		for i, x := range Websites {
			if x.ChatId == update.Message.Chat.ID {
				stri := strconv.Itoa(i)
				urls = append(urls, x.Url + " " + stri)
			}
		}
		str := fmt.Sprintf("list: %s", urls)
		log.Printf("printing value %s", Websites)
		sendTelegramBotMessage(str, update.Message.Chat.ID)
		break
	}

}


//func remove(s []Website, r Website) []Website {
//	for i, v := range s {
//		if v == r {
//			return append(s[:i], s[i+1:]...)
//		}
//	}
//	return s
//}

func removeIndex(slice []Website, s int) []Website {
	return append(slice[:s], slice[s+1:]...)
}
