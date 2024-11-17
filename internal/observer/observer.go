package observer

import (
	"log"
	"os"
	"terminer/internal/bot"

	"github.com/joho/godotenv"
)

type Observer interface {
	Notify(chatID string, message string)
}

// ConcreteObserver - конкретний спостерігач, який реагує на сповіщення від субʼєкта
type ConcreteObserver struct {
}

func (c *ConcreteObserver) Notify(chatID string, message string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	telegram_bot := bot.NewTelegramBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	telegram_bot.Notify(chatID, message)
}

func NewObserver(id string) *ConcreteObserver {
	return &ConcreteObserver{}
}
