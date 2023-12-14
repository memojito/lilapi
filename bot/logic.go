package bot

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"github.com/memojito/lilapi/db"
)

type Transaction struct {
	name   string
	price  int
	userID int64
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func InitBot(token string, conn *pgx.Conn) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 60

	c := &db.Conn{Conn: conn}

	updateBot(bot, update, c)
}

func updateBot(bot *tgbotapi.BotAPI, update tgbotapi.UpdateConfig, conn *db.Conn) {
	updates := bot.GetUpdatesChan(update)

	for u := range updates {
		if u.Message != nil { // If we got a message
			if _, err := conn.GetUser(u.Message.From.ID); err != nil { // check if user is new
				conn.SaveUser(u.Message.From.ID, u.Message.From.FirstName, u.Message.From.LastName)
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Hi, "+u.Message.From.FirstName+"! You're new around here!")
				bot.Send(msg)
			}
			if u.Message.IsCommand() {
				handleCommand(bot, u.Message, conn)
			} else {
				transaction, err := handleTransaction(bot, u.Message)
				category, _ := conn.GetCategoryByName(transaction.Name, transaction.UserID)
				if err != nil {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "Invalid Transaction!")
					bot.Send(msg)
				} else if category.ID == 0 {
					msg := tgbotapi.NewMessage(u.Message.Chat.ID, "\""+transaction.Name+"\" is a new category!")
					bot.Send(msg)

					handleCategory(bot, u.Message, transaction.Name, transaction.Value)
				} else {
					//save transaction
					conn.SaveTransaction(transaction.Name, transaction.Value, transaction.UserID, category.ID)
					transactions, _ := conn.GetTransactions(u.Message.From.ID, time.Now())

					//count daily total
					total := float64(countDailyTotal(transactions)) / 100
					handleTotal(bot, u.Message.Chat.ID, total)
				}
			}
		} else if u.CallbackQuery != nil {
			handleCallbackQuery(bot, u.CallbackQuery, conn)
		}
	}
}

func countDailyTotal(transactions []db.Transaction) int {
	var total int
	for _, t := range transactions {
		total += t.Value
	}
	return total
}

func countWeeclyTotal(transactions []db.Transaction) int {
	return 0
}

func handleTotal(bot *tgbotapi.BotAPI, chatID int64, total float64) {
	totalString := strconv.FormatFloat(total, 'f', 2, 64) //
	msg := tgbotapi.NewMessage(chatID, "Spent today: "+totalString+"€")
	bot.Send(msg)
}

func handleCategory(bot *tgbotapi.BotAPI, message *tgbotapi.Message, name string, value int) {
	yesButton := tgbotapi.NewInlineKeyboardButtonData("Yes", "Saving category "+name+" "+strconv.Itoa(value))
	noButton := tgbotapi.NewInlineKeyboardButtonData("No", "Reset")

	//log.Printf("Value is: %s", strconv.Itoa(value))

	msg := tgbotapi.NewMessage(message.Chat.ID, "Do you want to save \""+name+"\" as a new category?")
	button := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(yesButton),
		tgbotapi.NewInlineKeyboardRow(noButton))
	msg.ReplyMarkup = button
	bot.Send(msg)
}

func handleTransaction(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (db.Transaction, error) {
	parts := strings.Fields(message.Text)
	if len(parts) != 2 {
		return db.Transaction{}, errors.New("Invalid Transaction!")
	}

	name := parts[0]

	// validate and transform transaction
	priceStr := parts[1]
	if !strings.Contains(priceStr, ",") && !strings.Contains(priceStr, ".") {
		priceStr = priceStr + "00"
	} else {
		pattern := `[,.].{2}`
		match, _ := regexp.MatchString(pattern, priceStr)
		if !match {
			return db.Transaction{}, errors.New("Invalid Transaction!")
		}
	}
	priceStr = strings.ReplaceAll(priceStr, ",", "")
	priceStr = strings.ReplaceAll(priceStr, ".", "")

	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return db.Transaction{}, errors.New("Invalid Transaction!")
	}

	return db.Transaction{
		Name:   strings.ToLower(name),
		Value:  price,
		UserID: message.From.ID,
	}, nil
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

func handleCallbackQuery(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery, conn *db.Conn) {
	var answer tgbotapi.CallbackConfig

	answer = tgbotapi.NewCallback(cq.ID, cq.Data)
	bot.Send(answer)

	//log.Println(cq.Message.Text)

	s := cq.Data

	switch {
	case strings.Contains(s, "I will!"):
		genGuide(bot, cq.Message)
	case strings.Contains(s, "Have fun!"):
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "What did you spend today?")
		bot.Send(msg)
	case strings.Contains(s, "Saving category"):
		parts := strings.Fields(cq.Data)
		name := parts[2]
		valueStr := parts[3]
		value, _ := strconv.Atoi(valueStr)
		id := conn.SaveCategory(name, cq.From.ID)
		conn.SaveTransaction(name, value, cq.From.ID, id)

		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "Category \""+name+"\" saved!")
		bot.Send(msg)

		transactions, _ := conn.GetTransactions(cq.From.ID, time.Now())

		//count daily total
		total := float64(countDailyTotal(transactions)) / 100
		handleTotal(bot, cq.Message.Chat.ID, total)
	default:
		msg := tgbotapi.NewMessage(cq.Message.Chat.ID, "What did you spend today?")
		bot.Send(msg)
	}
}

func handleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, conn *db.Conn) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	var msg tgbotapi.MessageConfig

	helpButton := tgbotapi.NewInlineKeyboardButtonData("Help me", "I will!")
	skipButton := tgbotapi.NewInlineKeyboardButtonData("Skip", "Have fun!")

	switch message.Text {
	case "/start":
		msg = tgbotapi.NewMessage(message.Chat.ID, "Welcome! This bot helps you to get more control over your expenses")
		bot.Send(msg)
		msg = tgbotapi.NewMessage(message.Chat.ID, "Do you want to know how to use me?")
		button := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(helpButton),
			tgbotapi.NewInlineKeyboardRow(skipButton))
		msg.ReplyMarkup = button
		bot.Send(msg)
	case "/help":
		genGuide(bot, message)
	case "/weekly":
		//transactions, _ := conn.GetTransactionsByCategory(message.From.ID, GetStartDayOfWeek(time.Now()), time.Now(), 0)
	}
}
