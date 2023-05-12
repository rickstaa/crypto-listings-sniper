// Description: This package contains functions for creating Discord embeds.
package discordEmbeds

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rickstaa/crypto-listings-sniper/utils"
)

// Convert hex color to int.
func HexColorToInt(color string) int {
	colorInt, err := strconv.ParseUint(color, 16, 64)
	if err != nil {
		log.Fatalf("Error parsing color: %v", err)
	}
	return int(colorInt)
}

// Asset embed.
var (
	ASSET_EMBED = discordgo.MessageEmbed{
		Color: HexColorToInt("F3BA2F"),
		Image: &discordgo.MessageEmbedImage{URL: "https://t4.ftcdn.net/jpg/04/46/35/17/360_F_446351747_WHAenLH7njEwEAuDf3aJ7Q3WFX9FM18s.jpg"},
	}
)

// Returns a string containing a message for a new asset.
func newAssetMessage(asset string) discordgo.MessageEmbed {
	embed := ASSET_EMBED
	embed.Title = fmt.Sprintf("💎 Binance listed new asset (%s)", asset)
	embed.URL = utils.CreateBinanceURL(asset)
	return embed
}

// Return a string containing a message for a removed asset.
func removedAssetMessage(asset string) discordgo.MessageEmbed {
	embed := ASSET_EMBED
	embed.Title = fmt.Sprintf("🗑 Binance removed asset (%s)\n", asset)
	embed.Image = nil
	return embed
}

// Returns a string containing a Discord message for a new or removed asset.
func AssetEmbed(removed bool, asset string) discordgo.MessageEmbed {
	switch removed {
	case true:
		return removedAssetMessage(asset)
	default:
		return newAssetMessage(asset)
	}
}
