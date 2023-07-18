package clientchanel

import (
	"Terra-bot-api/dto"
	"Terra-bot-api/internal/keyboards"
	"Terra-bot-api/tgUtil"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	// "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type UpdateHandlerFunc func(update tgbotapi.Update, bot *tgbotapi.BotAPI)

type Beluga struct {
	Logger *zap.Logger
	State  map[int64]UpdateHandlerFunc
	Sword  sync.Mutex //rw mutex
}

func New() *Beluga {
	bl := &Beluga{
		State: make(map[int64]UpdateHandlerFunc),
		Sword: sync.Mutex{},
	}
	return bl
}

func (b *Beluga) Get(chadID int64) (f UpdateHandlerFunc, ok bool) {
	// b.Sword.Lock()
	f, ok = b.State[chadID]
	// b.Sword.Unlock()
	return f, ok
}

func (b *Beluga) Add(chadID int64, f func(update tgbotapi.Update, bot *tgbotapi.BotAPI)) {
	// b.Sword.Lock()
	b.State[chadID] = f
	// b.Sword.Unlock()
}

func (b *Beluga) DigitalSignature(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			// logrus.Error(err)
		}
		switch update.CallbackQuery.Data {
		case "Меню":
			// logrus.Printf(update.CallbackQuery.Data)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
			msg.ReplyMarkup = keyboards.StartKeyBoard
			tgUtil.SendBotMessage(msg, bot)
			delete(b.State, update.CallbackQuery.Message.Chat.ID)
		}

	} else {

		amount := update.Message.Text
		// log.Printf(txt)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Заявка на сумму: "+amount+" \n\n Переведите точную сумму на эти реквизиты \n\n"+request(amount))
		msg.ReplyMarkup = keyboards.ToMainTheme
		tgUtil.SendBotMessage(msg, bot)

	}
}

func (b *Beluga) Feedback(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery != nil {

		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := bot.Request(callback); err != nil {
			// logrus.Error(err)
		}
		switch update.CallbackQuery.Data {
		case "Меню":
			// logrus.Printf(update.CallbackQuery.Data)
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, keyboards.MenuReply)
			msg.ReplyMarkup = keyboards.StartKeyBoard
			tgUtil.SendBotMessage(msg, bot)
			delete(b.State, update.CallbackQuery.Message.Chat.ID)
		}

	} else {
		txt := update.Message.Chat.FirstName + " " + update.Message.Chat.LastName + "\n" + "@" + update.Message.Chat.UserName + "\n"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, keyboards.GratitudeReply)
		msg.ReplyMarkup = keyboards.ToMainTheme
		tgUtil.SendBotMessage(msg, bot)
		/*photoArray := update.Message.Photo
		photoLastIndex := len(photoArray) - 1
		photo := photoArray[photoLastIndex]
		reqBodyBytes := new(bytes.Buffer)
		json.NewEncoder(reqBodyBytes).Encode(photo)
		photoFileBytes := tgbotapi.FileBytes{
			Name:  "123",
			Bytes: reqBodyBytes.Bytes(),
		}*/
		msg = tgbotapi.NewMessageToChannel("-936178018", txt)
		//msg1 := tgbotapi.NewPhotoToChannel("-936178018", photoFileBytes) // канал вынести в отдельную переменную окружения
		tgUtil.SendBotMessage(msg, bot)
		//tgUtil.SendBotMessage1(msg1, bot)
	}
}

func request(amount string) string {

	url := "https://admin.paylama.io/api/api/payment/p2p_form"

	payload := strings.NewReader("{\n  \"payerID\": \"test\",\n  \"amount\": " + amount + ",\n  \"expireAt\": 3600,\n  \"comment\": \"comment\",\n  \"clientIP\": \"172.0.0.1\",\n  \"currencyID\": 1\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("API-Key", "oWHnJumI2IjLNxlrMtw6WDvnPCSQUmdkMsSzHf88pxCT6SCn4giuCBOd9zaTD5guK3VhbmpW206fTDJeTwxJjDLXKIcqwwIKZGIF8nqqHRD7Yi46Te2tYtuL7cNmKgqszc7jfR9fkzXaTkXD2mMS2knNTmYBLsM0QGEE52npo6J3zAj3yxDATN4M8UbV49tWeOBmcnoZJMyMWQiOU9tJMNrXaFiyt2sg604Yv7vVFdvo6jjkzljkw3MZKdrU4GEl")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var huy dto.Response
	err := json.Unmarshal(body, &huy)
	if err != nil {
		return "Что-то в коде сломалось, напишите @Glav_Control_IT"
	}
	if huy.Success == true {
		return huy.FormURL + "\nExternalID - " + huy.ExternalID
	} else {
		if huy.Cause == "No mid found for client!" {
			return "Проблемы, напишите @Glav_Control_IT"
		} else if huy.Cause == "An error has occurred. Please try again later." {
			return "Сломалась лама, напишите @Glav_Control_IT"
		} else if huy.Cause == "minimum amount is 300" || huy.Cause == "maximum amount is 300000" {
			return "Лимиты на создание от 300 до 300к, проверьте правильность введенной суммы"
		} else if huy.Cause == "Validation error." {
			return "Напишите сумму цифрами"
		} else {
			return "Ничего не вернулось"
		}
	}
}


