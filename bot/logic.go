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
				handleCommand(bot, u.Message)
			}
		} else if u.CallbackQuery != nil {
			handleCallbackQuery(bot, u.CallbackQuery)
		}
	}
}

func genReply(bot *tgbotapi.BotAPI, message *tgbotapi.Message) tgbotapi.MessageConfig {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	return msg
}

func genGuide(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var msg tgbotapi.MessageConfig
	msg = tgbotapi.NewMessage(message.Chat.ID, "Right now I can help you to count your daily expenses and give you some statistics.")
	bot.Send(msg)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Start by entering the name of the expense and the price in euros. Here is an example:")
	bot.Send(msg)

	msg = tgbotapi.NewMessage(message.Chat.ID, "coffee 3,29")
	bot.Send(msg)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	var answer tgbotapi.CallbackConfig

	answer = tgbotapi.NewCallback(cq.ID, cq.Data)
	bot.Send(answer)

	switch cq.Data {
	case "I will!":
		genGuide(bot, cq.Message)
	case "Have fun!":
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "What did you spend today?")
		bot.Send(msg)
	}
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	var msg tgbotapi.MessageConfig

	helpButton := tgbotapi.NewInlineKeyboardButtonData("Help me", "I will!")
	skipButton := tgbotapi.NewInlineKeyboardButtonData("Skip", "Have fun!")

	if message.Text == "/start" {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Welcome! This bot helps you to get more control over your expenses")
		bot.Send(msg)
		msg = tgbotapi.NewMessage(message.Chat.ID, "Do you want to know how to use me?")
		button := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(helpButton),
			tgbotapi.NewInlineKeyboardRow(skipButton))
		msg.ReplyMarkup = button
		bot.Send(msg)
	} else if message.Text == "/help" {
		genGuide(bot, message)
	}
}
