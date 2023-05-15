// Description: The telegram package contains functions for interacting with the Telegram API.
package telegram

import (
	"log"

	"github.com/rickstaa/crypto-listings-sniper/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/rickstaa/crypto-listings-sniper/messaging/telegram/telegramMessages"
)

// SendTelegramMessage sends a Telegram message to a specified chat.
func SendTelegramMessage(telegramBot *telego.Bot, chatID int64, message string) {
	msg := tu.Message(tu.ID(chatID), message)
	msg.ParseMode = telego.ModeHTML
	_, err := telegramBot.SendMessage(msg)
	if err != nil {
		log.Printf("WARNING: Error sending message '%s' to channel '%d': %v", msg.Text, chatID, err)
	}
}

// SendAssetTelegramMessage send a new/removed asset Telegram message to a specified chat.
func SendAssetTelegramMessage(telegramBot *telego.Bot, chatID int64, removed bool, asset string, assetInfo binance.Symbol) {
	message := telegramMessages.AssetMessage(removed, asset, utils.CreateBinanceURL(asset), assetInfo)
	SendTelegramMessage(telegramBot, chatID, message)
}

// Send a announcement Telegram message to the specified chat.
func SendAnnouncementTelegramMessage(telegramBot *telego.Bot, chatID int64, announcementCode, announcementTitle string) {
	message := telegramMessages.AnnouncementMessage(utils.CreateBinanceArticleURL(announcementCode, announcementTitle), announcementTitle)
	SendTelegramMessage(telegramBot, chatID, message)
}
