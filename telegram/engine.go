package telegram

import (
	"github.com/aktelion/quiz-test/quizplease"
	site_grabber "github.com/aktelion/quiz-test/quizplease/site-grabber"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BotType int

const (
	PollBot BotType = iota
	Webhook
)

type Command string

const (
	Help         = "help"
	Schedule     = "schedule"
	FullSchedule = "full_schedule"
	Rating       = "rating"
)

const RatingUrl = "https://moscow.quizplease.ru/rating?QpRaitingSearch%5Bgeneral%5D=1&QpRaitingSearch%5Bleague%5D=1&QpRaitingSearch%5Btext%5D=%D0%98%D0%BC%D0%B1%D0%B8%D1%80%D0%BD%D0%B0%D1%8F+%D0%BA%D0%B0%D0%BC%D0%B1%D0%B0%D0%BB%D0%B0"
const ScheduleUrl = "https://moscow.quizplease.ru/schedule"

var gameFilter = quizplease.GameFilter{
	FilterOnline:         true,
	FilterSubject:        true,
	FilterUnwantedPlaces: false,
}

var baseKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Расписание", Schedule),
		tgbotapi.NewInlineKeyboardButtonData("Рейтинг", Rating),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Полное расписание", FullSchedule),
	),
)

func replyCommand(msg *tgbotapi.MessageConfig, command string) {
	switch command {
	case Schedule, FullSchedule:
		msg.ReplyMarkup = baseKeyboard
		schedule, err := site_grabber.ParseSchedule(ScheduleUrl)
		if err != nil {
			log.Println("Can't load or parse schedule")
		}
		msg.Text = FormatSchedule(schedule, &gameFilter)
	case Rating:
		msg.ReplyMarkup = baseKeyboard
		rating, err := site_grabber.ParseRating(RatingUrl)
		if err != nil {
			log.Println("Can't load or parse rating")
		}
		msg.Text = FormatTeam(&quizplease.Team{
			Name:   "Имбирная Камбала",
			Rank:   quizplease.NewRank(rating.AllScores),
			Rating: *rating,
		})
	//case Help:
	//	msg.ReplyMarkup = baseKeyboard
	//	msg.Text = "Need some help? Ok..."
	default:
		msg.ReplyMarkup = baseKeyboard
		msg.Text = "..."
	}
	msg.ParseMode = "HTML"
}

func StartBot(token string, botType BotType) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			// Construct a new message from the given chat ID and containing
			// the text that we received.
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if update.Message.IsCommand() {
				replyCommand(&msg, update.Message.Command())
			}

			// Send the message.
			if _, err = bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
			replyCommand(&msg, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}
