package main

import (
	"fmt"
	"log"
	"telegram"
	post "telegram/server"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const eventRegBotToken = "603277423:AAGXuWe-J2czAX2AVgqH0dT44extnFFFksA" //bot's token

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🏠 Home"),
		tgbotapi.NewKeyboardButton("🗓 Join event"),
	),
)

var eventsList = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Event 1"),
		tgbotapi.NewKeyboardButton("Event 2"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Event 3"),
		tgbotapi.NewKeyboardButton("Event 4"),
	),
)

var EventSignMap map[int]*telegram.EventSign

func init() {
	EventSignMap = make(map[int]*telegram.EventSign)
}

func main() {
	var (
		bot        *tgbotapi.BotAPI
		err        error
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
		botUser    tgbotapi.User
	)
	bot, err = tgbotapi.NewBotAPI(eventRegBotToken)
	if err != nil {
		log.Panic("bot init error", err.Error())
		return
	}

	botUser, err = bot.GetMe()
	if err != nil {
		log.Panic("bot getme error", err.Error())
		return
	}

	fmt.Printf("auth ok! bot is: %s\n", botUser.FirstName)

	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel, err = bot.GetUpdatesChan(updConfig)

	if err != nil {
		log.Panic("update channel error", err.Error())
	}

	for {
		update = <-updChannel
		if update.Message != nil {
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "test" {
					msgConfig := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Test command !")
					bot.Send(msgConfig)
				} else if cmdText == "menu" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Menu")
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)
				}
			} else {
				if update.Message.Text == mainMenu.Keyboard[0][1].Text {

					EventSignMap[update.Message.From.ID] = new(telegram.EventSign)
					EventSignMap[update.Message.From.ID].State = 0

					fmt.Printf(
						"ID %d; from %d; message: %s; chatID: %d\n",
						update.Message.MessageID,
						update.Message.From.ID,
						update.Message.Text,
						update.Message.Chat.ID)

					msgConfig := tgbotapi.NewMessage(
						update.Message.Chat.ID,
						"Enter your email")
					msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(msgConfig)
				} else {
					es, ok := EventSignMap[update.Message.From.ID]
					if ok {
						if es.State == 0 {
							es.Email = update.Message.Text
							msgConfig := tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"Enter your phone number")

							bot.Send(msgConfig)
							es.State = 1
						} else if es.State == 1 {
							es.Phone = update.Message.Text
							msgConfig := tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"Enter Event name:")
							msgConfig.ReplyMarkup = eventsList
							bot.Send(msgConfig)
							es.State = 2
						} else if es.State == 2 {
							es.Event = update.Message.Text
							msgConfig := tgbotapi.NewMessage(
								update.Message.Chat.ID,
								"successfully !")
							msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							bot.Send(msgConfig)
							delete(EventSignMap, update.Message.From.ID)
							//post to site
							err = post.SendPost(es)
							if err != nil {
								fmt.Printf("send post error: %v\n", err)
							}

						}
						fmt.Printf("state: %+v\n", es)
					} else {
						msgConfig := tgbotapi.NewMessage(
							update.Message.Chat.ID,
							"Nothing happened !")
						msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						bot.Send(msgConfig)
					}
				}
			}
		} else {
			fmt.Printf("not a message... %+v\n", update)
		}
	}
}
