package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type telegramBot struct {
	ChatID  string
	Token   string
	Message string
}

func (t *telegramBot) Notify(chatID string, message string) error {
	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		return err
	}
	// Створюємо бота
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Fatalf("Помилка при створенні бота: %v", err)
	}

	// Створюємо повідомлення
	msg := tgbotapi.NewMessage(chatIDInt, message)
	msg.ParseMode = "MarkdownV2"
	// Відправляємо повідомлення
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatalf("Помилка при відправленні повідомлення: %v", err)
	}

	// fmt.Println("Повідомлення успішно відправлено")
	return nil
}

func NewTelegramBot(token string) *telegramBot {
	return &telegramBot{
		Token: token,
	}
}
