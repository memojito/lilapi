package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func InitBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	updateBot(bot, update)
}

func updateBot(bot *tgbotapi.BotAPI, update tgbotapi.UpdateConfig) {
	updates := bot.GetUpdatesChan(update)

	for u := range updates {
		if u.Message != nil { // If we got a message
			if u.Message.IsCommand() {
				HandleCommand(u.Message)
			}
			msg := GenReply(u.Message)

			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Failed to send message %v", err)
				return
			}
		}
	}
}

func GenReply(message *tgbotapi.Message) tgbotapi.MessageConfig {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	return msg
}

func HandleCommand(message *tgbotapi.Message) tgbotapi.MessageConfig {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	var msg tgbotapi.MessageConfig

	if message.Text == "/start" {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Welcome! This bot helps you to get more control over your expenses")
	} else if message.Text == "/help" {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Do you want to know how to use me?")
		//button := tgbotapi.InlineKeyboardButton{Text: "yes", CallbackData: "help"}
	}

	return msg
}
