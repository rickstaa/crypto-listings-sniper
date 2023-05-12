// Description: Contains functions for interacting with the Telegram API.
package telegram

import (
	"log"

	"github.com/rickstaa/crypto-listings-sniper/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/rickstaa/crypto-listings-sniper/messaging/telegram/telegramMessages"
)

// Send a Telegram message to the specified chat.
func SendTelegramMessage(telegramBot *telego.Bot, chatID int64, message string) {
	msg := tu.Message(tu.ID(chatID), message)
	msg.ParseMode = telego.ModeHTML
	_, err := telegramBot.SendMessage(msg)
	if err != nil {
		log.Printf("WARNING: Error sending message '%s' to channel '%d': %v", msg.Text, chatID, err)
	}
}

// Send a asset Telegram message to the specified chat.
func SendAssetTelegramMessage(telegramBot *telego.Bot, chatID int64, removed bool, asset string, symbolInfo binance.Symbol) {
	message := telegramMessages.AssetMessage(removed, asset, utils.CreateBinanceURL(asset), symbolInfo)
	SendTelegramMessage(telegramBot, chatID, message)
}
