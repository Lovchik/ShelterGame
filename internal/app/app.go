package main

import (
	"ShelterGame/internal/config"
	"ShelterGame/internal/database/sqlite"
	"ShelterGame/internal/dtos"
	"ShelterGame/internal/render"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"gorm.io/gorm/utils"
	"log/slog"
	"os"
	"strconv"
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
		var count int64
		sqlite.GetDB().Model(&dtos.User{}).Where("chat_id=?", update.Message.Chat.ID).Debug().Count(&count)
		if count == 0 {
			user := dtos.User{
				Username: update.Message.From.Username,
				ChatId:   update.Message.Chat.ID,
			}
			sqlite.GetDB().Model(&dtos.User{}).Debug().Create(&user)
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "Congratulation! You are now connected to the chat."))
		} else {
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "You are already signed"))

			slog.Info("User with chat_id: ", update.Message.Chat.ID, "exists")
		}
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		/*if len(gameMembers) <= 5 {
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "Cant start game.Add more members.").WithReplyMarkup(keyboard))
		 return
		}*/

		for index, member := range gameMembers {
			var chatId int64
			sqlite.GetDB().Raw("select chat_id from users where username=?", member).Debug().Scan(&chatId)
			result := render.Render(strconv.Itoa(index + 1))
			_, _ = bot.SendMessage(tu.Message(telego.ChatID{
				ID: chatId,
			}, result))

		}

	}, th.CommandEqual("startGame"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		keyboard = tu.Keyboard(getUsers()).WithSelective().WithResizeKeyboard().WithInputFieldPlaceholder("Select something")

		if gameMembers != nil && len(gameMembers) != 0 {
			message := ""
			for index, member := range gameMembers {
				message += strconv.Itoa(index+1) + ". " + member + "\n"
			}
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), message).WithReplyMarkup(keyboard))
		} else {
			_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), "No any players").WithReplyMarkup(keyboard))

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
				gameMembers = remove(gameMembers, getNumber(gameMembers, update.Message.Text[1:]))
				_, _ = bot.SendMessage(tu.Message(update.Message.Chat.ChatID(), update.Message.Text[1:]+" removed"))
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

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func getNumber(slice []string, value string) int {
	for i, s := range slice {
		if s == value {
			return i
		}
	}
	return -1
}
