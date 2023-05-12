// Description: Contains functions that are related to Binance.
package binance

import (
	"context"
	"log"

	dc "github.com/rickstaa/crypto-listings-sniper/discord"
	tg "github.com/rickstaa/crypto-listings-sniper/telegram"
	"github.com/rickstaa/crypto-listings-sniper/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
)

// Retrieve assets from Binance.
func RetrieveBinanceAssets(binanceClient *binance.Client) (assets []string) {
	priceService := binanceClient.NewListPricesService()
	listingPrices, err := priceService.Do(context.Background())
	assets = make([]string, len(listingPrices))
	if err != nil {
		log.Fatalf("Error retrieving Binance listing prices: %v", err)
	}
	for i, s := range listingPrices {
		assets[i] = s.Symbol
	}

	return assets
}

// Check Binance for new listings or de-listings and post Telegram message.
func BinanceListingsCheck(oldAssets *[]string, binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) {
	assets := RetrieveBinanceAssets(binanceClient)

	// Check for new assets and post channel messages.
	if len(*oldAssets) != 0 { // Do not check if first run of the program.
		removed, newAssets := utils.CompareLists(*oldAssets, assets)
		for _, s := range newAssets {
			// Send telegram message.
			if enableTelegramMessage {
				go tg.SendAssetTelegramMessage(telegramBot, telegramChatID, removed, s)
			}

			// Send discord message.
			if enableDiscordMessages {
				go dc.SendAssetDiscordMessage(discordBot, discordChannelIDs, removed, s)
			}
		}

		// Store updated assets list in JSON files.
		if len(newAssets) != 0 {
			utils.StoreOldListings(assets)
			*oldAssets = assets
		}
	}
}
