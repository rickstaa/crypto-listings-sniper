// Decription: This package contains functions for creating discord embeds.
package discordEmbeds

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	DEFAULT_EMBED = discordgo.MessageEmbed{
		URL:   "https://www.binance.com/en/support/announcement/360050980292",
		Color: HexColorToInt("F3BA2F"),
		Image: &discordgo.MessageEmbedImage{URL: "https://t4.ftcdn.net/jpg/04/46/35/17/360_F_446351747_WHAenLH7njEwEAuDf3aJ7Q3WFX9FM18s.jpg"},
	}
)

// Convert hex color to int.
func HexColorToInt(color string) int {
	colorInt, err := strconv.ParseUint(color, 16, 64)
	if err != nil {
		log.Fatalf("Error parsing color: %v", err)
	}
	return int(colorInt)
}

// Returns a string containing a message for a new SPOT trading pair.
func newTradingPairMessage(symbol string) discordgo.MessageEmbed {
	embed := DEFAULT_EMBED
	embed.Title = fmt.Sprintf("⚖️ Binance listed new SPOT trading pair (%s)", symbol)
	return embed
}

// Return a string containing a message for a removed SPOT trading pair.
func removedTradingPairMessage(symbol string) discordgo.MessageEmbed {
	embed := DEFAULT_EMBED
	embed.Title = fmt.Sprintf("🗑 Binance removed SPOT trading pair (%s)\n", symbol)
	embed.Image = nil
	return embed
}

// Returns a string containing a message for a new SPOT base asset.
func newBaseAssetMessage(symbol string) discordgo.MessageEmbed {
	embed := DEFAULT_EMBED
	embed.Title = fmt.Sprintf("💎 Binance listed new SPOT asset (%s)", symbol)
	return embed
}

// Returns a string containing a message for a removed SPOT base asset.
func removedBaseAssetMessage(symbol string) discordgo.MessageEmbed {
	embed := DEFAULT_EMBED
	embed.Title = fmt.Sprintf("🗑 Binance removed SPOT asset (%s)\n", symbol)
	embed.Image = nil
	return embed
}

// Returns a string containing a Discord message for a new or removed SPOT base asset.
func BaseAssetEmbed(removed bool, symbol string) discordgo.MessageEmbed {
	switch removed {
	case true:
		return removedBaseAssetMessage(symbol)
	default:
		return newBaseAssetMessage(symbol)
	}
}

// Returns a string containing a Discord message for a new or removed SPOT trading pair.
func TradingPairEmbed(removed bool, symbol string) discordgo.MessageEmbed {
	switch removed {
	case true:
		return removedTradingPairMessage(symbol)
	default:
		return newTradingPairMessage(symbol)
	}
}
