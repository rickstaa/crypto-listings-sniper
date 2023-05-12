// Description: Contains functions that are related to Binance.
package exchanges

import (
	"context"
	"log"

	bn "github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	"github.com/rickstaa/crypto-listings-sniper/utils"
)

// Retrieve assets from Binance.
func RetrieveBinanceAssets(binanceClient *bn.Client) (assets []string) {
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

// Retrieve information about a given symbol.
func RetrieveSymbolInfo(binanceClient *bn.Client, symbol string) (symbolInfo bn.Symbol) {
	exchangeInfoService := binanceClient.NewExchangeInfoService()
	exchangeInfoService = exchangeInfoService.Symbols(symbol)
	exchangeInfo, err := exchangeInfoService.Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving Binance exchange info for symbol '%s': %v", symbol, err)
	}
	return exchangeInfo.Symbols[0]
}

// Check Binance for new listings or de-listings.
func BinanceListingsCheck(oldAssets *[]string, binanceClient *bn.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) (removed bool, changedAssets []string) {
	assets := RetrieveBinanceAssets(binanceClient)

	// Check for new assets and post channel messages.
	if len(*oldAssets) != 0 { // Do not check if first run of the program.
		removed, changedAssets = utils.CompareLists(*oldAssets, assets)

		// Store updated assets list in JSON files.
		if len(changedAssets) != 0 {
			utils.StoreOldListings(assets)
			*oldAssets = assets
		}
	}

	return removed, changedAssets
}
