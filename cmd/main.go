package main

import (

	//"github.com/joho/godotenv"

	clientchanel "Terra-bot-api/internal/client_chanel"
	"Terra-bot-api/internal/telegram_bot"
	"log"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	//"tg-bot-for-ts/repository"
	//"tg-bot-for-ts/service"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	bot, err := telegram_bot.NewBot(logger)
	if err != nil {
		logger.Fatal("failed to create a new BotAPI instance", zap.Error(err))
	}

	cc := clientchanel.New()

	err = bot.BotWorker(cc)
	if err != nil {
		// logger.Fatal("")
	}
}
