// Description: Package binance contains functions that are related to Binance.
package exchanges

import (
	"context"
	"log"
	"time"

	bn "github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	"github.com/rickstaa/crypto-listings-sniper/messaging"
	"github.com/rickstaa/crypto-listings-sniper/utils"
	"golang.org/x/time/rate"
)

// BinanceListingsChecker is a class that when started checks Binance for new listings or de-listings.
type BinanceListingsChecker struct {
	BinanceClient         *bn.Client
	TelegramBot           *telego.Bot
	TelegramChatID        int64
	EnableTelegramMessage bool
	DiscordBot            *discordgo.Session
	DiscordChannelIDs     []string
	EnableDiscordMessages bool
	OldAssets             *[]string
}

// NewBinanceListingsChecker creates a new BinanceListingsChecker.
func NewBinanceListingsChecker(binanceClient *bn.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) *BinanceListingsChecker {
	return &BinanceListingsChecker{
		BinanceClient:         binanceClient,
		TelegramBot:           telegramBot,
		TelegramChatID:        telegramChatID,
		EnableTelegramMessage: enableTelegramMessage,
		DiscordBot:            discordBot,
		DiscordChannelIDs:     discordChannelIDs,
		EnableDiscordMessages: enableDiscordMessages,
	}
}

// retrieveBinanceAssets retrieves a list with the available assets from Binance.
func (blc *BinanceListingsChecker) retrieveBinanceAssets() (assets []string) {
	priceService := blc.BinanceClient.NewListPricesService()
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

// retrieveSymbolInfo retrieves information about a given symbol from Binance.
func (blc *BinanceListingsChecker) retrieveSymbolInfo(symbol string) (assetInfo bn.Symbol) {
	exchangeInfoService := blc.BinanceClient.NewExchangeInfoService()
	exchangeInfoService = exchangeInfoService.Symbols(symbol)
	exchangeInfo, err := exchangeInfoService.Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving Binance exchange info for symbol '%s': %v", symbol, err)
	}
	return exchangeInfo.Symbols[0]
}

// BinanceListingsChecker checks Binance for new listings or de-listings.
func (blc *BinanceListingsChecker) binanceListingsCheck(oldAssets *[]string) (removed bool, changedAssets []string) {
	assets := blc.retrieveBinanceAssets()

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

// Start starts the BinanceListingsChecker.
func (blc *BinanceListingsChecker) Start() {
	// Retrieve old assets list.
	oldAssets := utils.RetrieveOldListings()

	// Check binance for new listings or de-listings and post Telegram/Discord message.
	r := rate.Every(1 * time.Millisecond)
	limiter := rate.NewLimiter(r, 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.

		// retrieve changed assets.
		removed, changedAssets := blc.binanceListingsCheck(&oldAssets)

		// Post messages.
		for _, asset := range changedAssets {
			// Retrieve symbol info.
			assetInfo := blc.retrieveSymbolInfo(asset)

			// Post telegram and discord messages.
			go messaging.SendAssetMessage(blc.BinanceClient, blc.TelegramBot, blc.TelegramChatID, blc.EnableTelegramMessage, blc.DiscordBot, blc.DiscordChannelIDs, blc.EnableDiscordMessages, removed, asset, assetInfo)
		}
	}
}
