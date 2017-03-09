package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
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
			processor(update, bot)
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
		msg := tgbotapi.NewMessage(chatID, "asdfasdf")
		msg.Text = "gagal nih!!! " + message
		bot.Send(msg)
	} else {
		log.Print("Chat id still 0")
	}
}

func processor(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		chatID := update.Message.Chat.ID
		msgID := update.Message.MessageID
		msg := tgbotapi.NewMessage(chatID, update.Message.Text)
		msg.ReplyToMessageID = msgID
		msg.Text = "echo " + update.Message.Text
		bot.Send(msg)
	}

	if update.Message == nil {
		return
	}

	// Parse responses
	response := strings.Split(update.Message.Text, " ")

	if response[0] == "" {
		sendTelegramBotMessage("Please enter in a command:\nstock [exchange] [symbol], trends", update.Message.Chat.ID)
		return
	}

	switch strings.ToLower(response[0]) {
	case "monitor":
		//@TODO need to check for valid url
		var url = response[1]
		var newWebsite = Website{Url: url, Interval: 5, ChatId: update.Message.Chat.ID}
		Websites = append(Websites, newWebsite)

		sendTelegramBotMessage("Added " + newWebsite.Url, update.Message.Chat.ID)
		break
	case "remove":
		//@TODO: add remove url entries here
		log.Printf("removing %s", update.Message.Text)

		break
	}

}
