// Description: The messaging package contains general functions for sending messages to supported messaging services.
// Note: Currently only Discord and Telegram are supported.

package messaging

import (
	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	dc "github.com/rickstaa/crypto-listings-sniper/messaging/discord"

	tg "github.com/rickstaa/crypto-listings-sniper/messaging/telegram"
)

// SendAssetMessage sends a new/removed assets message to the supported messaging services.
func SendAssetMessage(telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool, removed bool, asset string, assetInfo binance.Symbol) {
	if enableTelegramMessage {
		go tg.SendAssetTelegramMessage(telegramBot, telegramChatID, removed, asset, assetInfo)
	}

	if enableDiscordMessages {
		go dc.SendAssetDiscordMessage(discordBot, discordChannelIDs, removed, asset, assetInfo)
	}
}

// SendAnnouncementMessage sends a new announcement message to the messaging services.
func SendAnnouncementMessage(telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool, announcementCode string, announcementTitle string) {
	if enableTelegramMessage {
		go tg.SendAnnouncementTelegramMessage(telegramBot, telegramChatID, announcementCode, announcementTitle)
	}

	if enableDiscordMessages {
		go dc.SendAnnouncementDiscordMessage(discordBot, discordChannelIDs, announcementCode, announcementTitle)
	}
}
