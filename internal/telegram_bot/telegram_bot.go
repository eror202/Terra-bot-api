package telegram_bot

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	//"github.com/joho/godotenv"

	clientchanel "Terra-bot-api/internal/client_chanel"
	"Terra-bot-api/internal/keyboards"

	_ "github.com/lib/pq"
	//"tg-bot-for-ts/repository"
	//"tg-bot-for-ts/service"
	"Terra-bot-api/tgUtil"
)

type Bot struct {
	Logger  *zap.Logger
	Bot     *tgbotapi.BotAPI
	Updates tgbotapi.UpdatesChannel
}

func NewBot(logger *zap.Logger) (*Bot, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	bot, err := tgbotapi.NewBotAPI("6056577101:AAHUPXIbya0XBSkYC6z4LXOnF2vaJ_-XmQU")
	if err != nil {
		logger.Error("failed to create a new BotAPI instance", zap.Error(err))
		return nil, err
	}
	bot.Debug = true

	logger.Debug("Authorized on ", zap.String("account", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 3600

	updates := bot.GetUpdatesChan(u)

	return &Bot{
		Bot:     bot,
		Updates: updates,
	}, err
}

func (b *Bot) sendBotMessage(msg tgbotapi.MessageConfig) {
	if _, err := b.Bot.Send(msg); err != nil {
		b.Logger.Error("", zap.Error(err))
	}
}

func (b *Bot) BotWorker(cc *clientchanel.Beluga) error {
	for update := range b.Updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		var chatID int64

		if update.Message != nil {
			chatID = update.Message.Chat.ID
		}
		if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}

		if f, ok := cc.Get(chatID); ok {
			f(update, b.Bot)
			continue
		}

		if update.Message != nil {
			b.switcherMessage(update)
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := b.Bot.Request(callback); err != nil {
				b.Logger.Error("sending Chattable error", zap.Error(err))
				return err
			}

			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			b.switcherCallback(update, msg, cc)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.UnrecognizedСommand)
			tgUtil.SendBotMessage(msg, b.Bot)
		}
	}
	return nil
}

func (b *Bot) switcherMessage(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.StartReply)
		msg.ReplyMarkup = keyboards.StartKeyBoard

		tgUtil.SendBotMessage(msg, b.Bot)

	case "commands":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.CommandsReply)
		tgUtil.SendBotMessage(msg, b.Bot)

	case "main":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.MainReply)
		msg.ReplyMarkup = keyboards.StartKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.DefReply)
		tgUtil.SendBotMessage(msg, b.Bot)
	}
}

func (b *Bot) switcherCallback(update tgbotapi.Update, msg tgbotapi.MessageConfig, cc *clientchanel.Beluga) {

	switch update.CallbackQuery.Data {

	case "Создать заказ":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Введите верную сумму, минимальный лимит на создание транзакции - 300 рублей, максимальный"+
			" - 300к рублей")
		msg.ReplyMarkup = keyboards.ToMainThemeOrSendCheck
		tgUtil.SendBotMessage(msg, b.Bot)
		cc.Add(update.CallbackQuery.Message.Chat.ID, cc.DigitalSignature)

	case "Меню":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
		msg.ReplyMarkup = keyboards.StartKeyBoard
		tgUtil.SendBotMessage(msg, b.Bot)

	default:
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.InvalidReqReply)
		tgUtil.SendBotMessage(msg, b.Bot)
	}
}
