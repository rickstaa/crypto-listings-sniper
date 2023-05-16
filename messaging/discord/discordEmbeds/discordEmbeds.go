// Description: The discordEmbeds package contains functions for creating Discord embeds.
package discordEmbeds

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/rickstaa/crypto-listings-sniper/utils"
)

// Initialise default asset and announcement embeds.
var (
	ASSET_EMBED = discordgo.MessageEmbed{
		Color: utils.HexColorToInt("F3BA2F"),
		Image: &discordgo.MessageEmbedImage{URL: "https://t4.ftcdn.net/jpg/04/46/35/17/360_F_446351747_WHAenLH7njEwEAuDf3aJ7Q3WFX9FM18s.jpg"},
	}
	ANNOUNCEMENT_EMBED = discordgo.MessageEmbed{
		Color: utils.HexColorToInt("F3BA2F"),
		Image: &discordgo.MessageEmbedImage{URL: "https://t4.ftcdn.net/jpg/04/46/35/17/360_F_446351747_WHAenLH7njEwEAuDf3aJ7Q3WFX9FM18s.jpg"},
	}
)

// newAssetMessage returns a new asset embed.
func newAssetMessage(asset string, symbolInfo binance.Symbol) discordgo.MessageEmbed {
	embed := ASSET_EMBED
	embed.Title = fmt.Sprintf("💎 Binance listed new asset (%s)", asset)
	embed.Description = fmt.Sprintf("• **Base Asset:** %s\n", symbolInfo.BaseAsset) +
		fmt.Sprintf("• **Quota Asset:** %s\n", symbolInfo.QuoteAsset)
	embed.URL = utils.CreateBinanceURL(asset)
	return embed
}

// removedAssetMessage return a removed asset embed.
func removedAssetMessage(asset string) discordgo.MessageEmbed {
	embed := ASSET_EMBED
	embed.Title = fmt.Sprintf("🗑 Binance removed asset (%s)\n", asset)
	embed.Image = nil
	return embed
}

// AssetEmbed returns a asset Discord embed.
func AssetEmbed(removed bool, asset string, assetInfo binance.Symbol) discordgo.MessageEmbed {
	if removed {
		return removedAssetMessage(asset)
	}
	return newAssetMessage(asset, assetInfo)
}

// AnnouncementEmbed returns a new announcement embed.
func AnnouncementEmbed(url string, title string) discordgo.MessageEmbed {
	embed := ANNOUNCEMENT_EMBED
	embed.Title = fmt.Sprintf("📢 %s", title)
	embed.URL = url
	return embed
}
