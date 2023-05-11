// Description: Contains functions for interacting with the Telegram API.
package telegram

import (
	"log"

	"github.com/rickstaa/crypto-listings-sniper/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/rickstaa/crypto-listings-sniper/telegram/telegramMessages"
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

// Send a base asset Telegram message to the specified chat.
func SendBaseAssetTelegramMessage(telegramBot *telego.Bot, chatID int64, removed bool, symbol string) {
	message := telegramMessages.BaseAssetMessage(removed, symbol, utils.CreateBinanceURL(symbol)+"_USDT")
	SendTelegramMessage(telegramBot, chatID, message)
}

// Send a trading pair Telegram message to the specified chat.
func SendTradingPairTelegramMessage(telegramBot *telego.Bot, chatID int64, removed bool, symbol string, symbolInfo map[string]binance.Symbol) {
	message := telegramMessages.TradingPairMessage(removed, utils.CreateBinanceURL(symbol), symbolInfo[symbol].BaseAsset+"/"+symbolInfo[symbol].QuoteAsset)
	SendTelegramMessage(telegramBot, chatID, message)
}
