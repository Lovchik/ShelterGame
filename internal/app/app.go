package main

import (
	"ShelterGame/internal/config"
	"ShelterGame/internal/database/sqlite"
	"ShelterGame/internal/dtos"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"log/slog"
	"os"
)

func main() {
	bot, err := telego.NewBot(config.GetConfig().TelegramToken, telego.WithDefaultDebugLogger())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		user := dtos.User{
			Username: update.Message.From.Username,
			ChatId:   update.Message.Chat.ID,
		}
		sqlite.GetDB().Model(&dtos.User{}).Debug().Create(&user)
	}, th.CommandEqual("start"))

	/*for update := range updates {
		if update.Message != nil {
			chatID := tu.ID(update.Message.Chat.ID)

			if update.Message.Text == "Чебурек" {
				_, err = bot.SendMessage(tu.Message(chatID, "Ну дарова"))
				if err != nil {
					return
				}
				user := dtos.User{
					Username: update.Message.From.Username,
					ChatId:   chatID.ID,
				}
				sqlite.GetDB().Model(&dtos.User{}).Create(&user)

			} else {
				_, _ = bot.CopyMessage(tu.CopyMessage(chatID, chatID, update.Message.MessageID))

			}

		}

	}*/
	bh.Start()
}
