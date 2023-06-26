package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var StartKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать заказ", "Создать заказ")))

var ToMainTheme = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню")))

var ToMainThemeOrSendCheck = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню")))
