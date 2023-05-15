// Description: The messaging package contains general functions for sending messages to supported messaging services.

package messaging

import (
	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	dc "github.com/rickstaa/crypto-listings-sniper/messaging/discord"
	tg "github.com/rickstaa/crypto-listings-sniper/messaging/telegram"
)

// SendAssetMessage sends a message to the supported messaging services.
// NOTE: Currently only Telegram and Discord are supported.
func SendAssetMessage(binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool, removed bool, asset string, assetInfo binance.Symbol) {
	// Send telegram message.
	if enableTelegramMessage {
		go tg.SendAssetTelegramMessage(telegramBot, telegramChatID, removed, asset, assetInfo)
	}

	// Send discord message.
	if enableDiscordMessages {
		go dc.SendAssetDiscordMessage(discordBot, discordChannelIDs, removed, asset, assetInfo)
	}
}
