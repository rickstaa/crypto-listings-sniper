// Description: Contains general functions for sending messages to the supported messaging services.

package messaging

import (
	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	bc "github.com/rickstaa/crypto-listings-sniper/exchanges/binance"
	dc "github.com/rickstaa/crypto-listings-sniper/messaging/discord"
	tg "github.com/rickstaa/crypto-listings-sniper/messaging/telegram"
)

// SendAssetMessage sends a message to the messaging services.
func SendAssetMessage(binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool, removed bool, asset string) {
	// Retrieve asset information.
	symbolInfo := binance.Symbol{}
	if !removed {
		symbolInfo = bc.RetrieveSymbolInfo(binanceClient, asset)
	}

	// Send telegram message.
	if enableTelegramMessage {
		go tg.SendAssetTelegramMessage(telegramBot, telegramChatID, removed, asset, symbolInfo)
	}

	// Send discord message.
	if enableDiscordMessages {
		go dc.SendAssetDiscordMessage(discordBot, discordChannelIDs, removed, asset, symbolInfo)
	}
}
