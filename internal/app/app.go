package main

import (
	"ShelterGame/internal/config"
	"ShelterGame/internal/database/sqlite"
	"ShelterGame/internal/dtos"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"gorm.io/gorm/utils"
	"log/slog"
	"os"
)

func main() {
	bot, err := telego.NewBot(config.GetConfig().TelegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	var gameMembers []string

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	keyboard := tu.Keyboard(getUsers()).WithSelective().WithResizeKeyboard().WithInputFieldPlaceholder("Select something")
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		var count *int64
		sqlite.GetDB().Model(&dtos.User{}).Where("chat_id=?", update.Message.Chat.ID).Count(count)
		if count == nil {
			user := dtos.User{
				Username: update.Message.From.Username,
				ChatId:   update.Message.Chat.ID,
			}
			sqlite.GetDB().Model(&dtos.User{}).Debug().Create(&user)
		} else {
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "You are already signed"))

			slog.Info("User with chat_id: ", update.Message.Chat.ID, "exists")
		}
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "You are already signed").WithReplyMarkup(keyboard))

	}, th.CommandEqual("w"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if gameMembers != nil {
			message := ""
			for _, member := range gameMembers {
				message += member + "\n"
			}
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), message))
		} else {
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "No any players"))

		}

	}, th.CommandEqual("checkPlayers"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		slog.Info(update.Message.Chat.Username)
		if update.Message.Chat.Username == "Lovchik1Kg" {
			gameMembers = nil
		}
	}, th.CommandEqual("clearList"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		if update.Message.Text != "" && update.Message.Text[:1] == "-" && isUserExists(update.Message.Text[1:]) {
			if utils.Contains(gameMembers, update.Message.Text[1:]) {
				_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), update.Message.Text[1:]+" already added"))
			} else {
				gameMembers = append(gameMembers, update.Message.Text[1:])
				_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), update.Message.Text[1:]+" successfully added"))
			}
		}
	}, th.Any())

	bh.Start()
}

func getUsers() []telego.KeyboardButton {
	var users []dtos.User
	sqlite.GetDB().Model(&dtos.User{}).Find(&users)
	var buttons []telego.KeyboardButton

	for _, user := range users {
		buttons = append(buttons, telego.KeyboardButton{
			Text:            "-" + user.Username,
			RequestUsers:    nil,
			RequestChat:     nil,
			RequestContact:  false,
			RequestLocation: false,
			RequestPoll:     nil,
			WebApp:          nil,
		})
	}
	return buttons
}

func isUserExists(username string) bool {
	var count int64
	sqlite.GetDB().Model(&dtos.User{}).Where("username=?", username).Count(&count)
	if count == 0 {
		return false
	}
	return true
}
