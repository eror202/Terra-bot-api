package keyboards

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var StartKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать заказ", "Создать заказ"),
		tgbotapi.NewInlineKeyboardButtonData("Прикрепить чек", "Прикрепить чек")))

var ToMainTheme = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню"),
		tgbotapi.NewInlineKeyboardButtonData("Прикрепить чек", "Прикрепить чек")))

var ToMainThemeOrSendCheck = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню"),
		tgbotapi.NewInlineKeyboardButtonData("Прикрепить чек", "Прикрепить чек")))
