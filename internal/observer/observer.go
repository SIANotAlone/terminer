package observer

import (
	"log"
	"os"
	"terminer/internal/bot"

	"github.com/joho/godotenv"
)

const (
	TELEGRAM_BOT_TOKEN = "TELEGRAM_BOT_TOKEN"
)

type Observer interface {
	Notify(chatID string, message string)
}

// ConcreteObserver - конкретний спостерігач, який реагує на сповіщення від субʼєкта
type ConcreteObserver struct {
}

func (c *ConcreteObserver) Notify(chatID string, message string) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	telegram_bot := bot.NewTelegramBot(os.Getenv(TELEGRAM_BOT_TOKEN))
	telegram_bot.Notify(chatID, message)
}

func NewObserver(id string) *ConcreteObserver {
	return &ConcreteObserver{}
}
