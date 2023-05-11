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

// Retrieve SPOT assets and symbols from Binance.
func RetrieveBinanceSpotAssets(binanceClient *binance.Client) (baseAssets []string, symbols []string, symbolInfo map[string]binance.Symbol) {
	// Retrieve SPOT exchange info.
	exchangeInfoService := binanceClient.NewExchangeInfoService()
	exchangeInfoService = exchangeInfoService.Permissions("SPOT")
	exchangeInfo, err := exchangeInfoService.Do(context.Background())

	// Retrieve SPOT symbols.
	exchangeSymbols := exchangeInfo.Symbols
	if err != nil {
		log.Fatalf("Error retrieving Binance 'SPOT' symbols: %v", err)
	}
	symbols = make([]string, len(exchangeSymbols))
	symbolInfo = map[string]binance.Symbol{}
	for i, s := range exchangeSymbols {
		symbols[i] = s.Symbol
		symbolInfo[s.Symbol] = s
	}

	// Retrieve SPOT base assets.
	k := make(map[string]bool)
	for _, s := range exchangeSymbols {
		if _, value := k[s.BaseAsset]; !value {
			k[s.BaseAsset] = true
			baseAssets = append(baseAssets, s.BaseAsset)
		}
	}

	return baseAssets, symbols, symbolInfo
}

// Check Binance for new SPOT listings or de-listings and post telegram message.
func BinanceListingsCheck(oldBaseAssetsList *[]string, oldSymbolsList *[]string, binanceClient *binance.Client, telegramBot *telego.Bot, chatID int64, discordBot *discordgo.Session, discordChannelIDs []string) {
	baseAssets, symbols, symbolInfo := RetrieveBinanceSpotAssets(binanceClient)

	// Check for new base assets and post telegram message.
	if len(*oldBaseAssetsList) != 0 { // Do not check if first run of the program.
		removed, newBasebaseAssets := utils.CompareLists(*oldBaseAssetsList, baseAssets)
		for _, s := range newBasebaseAssets {
			// Send telegram message.
			go tg.SendBaseAssetTelegramMessage(telegramBot, chatID, removed, s)

			// Send discord message..
			go dc.SendBaseAssetDiscordMessage(discordBot, discordChannelIDs, removed, s)
		}

		// Store updated base assets list in JSON files.
		if len(newBasebaseAssets) != 0 {
			utils.StoreOldListings(baseAssets, nil)
			*oldBaseAssetsList = baseAssets
		}
	}

	// Check for new trading pairs and post telegram message.
	if len(*oldSymbolsList) != 0 { // Do not check if first run of the program.
		removed, newSymbols := utils.CompareLists(*oldSymbolsList, symbols)
		for _, s := range newSymbols {
			// Send telegram message.
			go tg.SendTradingPairTelegramMessage(telegramBot, chatID, removed, s, symbolInfo)

			// Send discord message.
			go dc.SendTradingPairDiscordMessage(discordBot, discordChannelIDs, removed, s)
		}

		// Store updated symbol list in JSON files.
		if len(newSymbols) != 0 {
			utils.StoreOldListings(nil, symbols)
			*oldSymbolsList = symbols
		}
	}
}
