package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"net/url"
	"log"
	"strings"
	"fmt"
)

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

func sendTelegramBotMessage(message string, chatID int64) {
	if (chatID != 0) {
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
		bot.Send(msg)
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
			var newWebsite = Website{Url: targetUrl, Interval: 5, ChatId: update.Message.Chat.ID}
			Websites = append(Websites, newWebsite)

			sendTelegramBotMessage("Added: " + newWebsite.Url, update.Message.Chat.ID)
		}
		break
	case "/remove":
		var targetUrl = response[1]
		var tobeRemovedWebsite = Website{Url: targetUrl, Interval: 5, ChatId: update.Message.Chat.ID}
		Websites = remove(Websites, tobeRemovedWebsite)
		log.Printf("removing %s", update.Message.Text)

		sendTelegramBotMessage("Removed: " + tobeRemovedWebsite.Url, update.Message.Chat.ID)
		break
	case "/list":
		urls := []string{}
		for _, x := range Websites {
			if x.ChatId == update.Message.Chat.ID {
				urls = append(urls, x.Url)
			}
		}
		str := fmt.Sprintf("list: %s", urls)
		log.Printf("printing value %s", str)
		sendTelegramBotMessage(str, update.Message.Chat.ID)
		break
	}

}


func remove(s []Website, r Website) []Website {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
