// Description: This package contains functions for checking for new or removed assets and trading pairs.
package checkers

import (
	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	"github.com/rickstaa/crypto-listings-sniper/utils"
)

// Check Binance for new SPOT listings or de-listings and post telegram message.
func BinanceListingsCheck(oldBaseAssetsList *[]string, oldSymbolsList *[]string, binanceClient *binance.Client, telegramBot *telego.Bot, chatID int64) {
	baseAssets, symbols, symbolInfo := utils.RetrieveBinanceSpotAssets(binanceClient)

	// Check for new base assets and post telegram message.
	if len(*oldBaseAssetsList) != 0 { // Do not check if first run of the program.
		removed, newBasebaseAssets := utils.CompareLists(*oldBaseAssetsList, baseAssets)
		for _, s := range newBasebaseAssets {
			go utils.SendBaseAssetTelegramMessage(telegramBot, chatID, removed, s)
		}

		// Store updated base assets list in JSON files.
		if len(newBasebaseAssets) != 0 {
			utils.StoreOldListings(nil, symbols)
			*oldBaseAssetsList = baseAssets
		}
	}

	// Check for new trading pairs and post telegram message.
	if len(*oldSymbolsList) != 0 { // Do not check if first run of the program.
		removed, newSymbols := utils.CompareLists(*oldSymbolsList, symbols)
		for _, s := range newSymbols {
			go utils.SendTradingPairTelegramMessage(telegramBot, chatID, removed, s, symbolInfo)
		}

		// Store updated symbol list in JSON files.
		if len(newSymbols) != 0 {
			utils.StoreOldListings(nil, symbols)
			*oldSymbolsList = symbols
		}
	}
}
