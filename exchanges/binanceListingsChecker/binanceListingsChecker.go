// Description: Package binanceListingsChecker contains a class that when started periodically checks Binance for new listings or de-listings and posts a message in set message channels.
package binanceListingsChecker

import (
	"context"
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	"github.com/rickstaa/crypto-listings-sniper/messaging"
	"github.com/rickstaa/crypto-listings-sniper/utils"
	"golang.org/x/time/rate"
)

// BinanceListingsChecker is a class that when started checks Binance for new listings or de-listings and posts a message in set message channels
type BinanceListingsChecker struct {
	BinanceClient         *binance.Client
	TelegramBot           *telego.Bot
	TelegramChatID        int64
	EnableTelegramMessage bool
	DiscordBot            *discordgo.Session
	DiscordChannelIDs     []string
	EnableDiscordMessages bool
	OldAssets             *[]string
}

// NewBinanceListingsChecker creates a new BinanceListingsChecker.
func NewBinanceListingsChecker(binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) *BinanceListingsChecker {
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
func (blc *BinanceListingsChecker) retrieveSymbolInfo(symbol string) (assetInfo binance.Symbol) {
	exchangeInfoService := blc.BinanceClient.NewExchangeInfoService()
	exchangeInfoService = exchangeInfoService.Symbols(symbol)
	exchangeInfo, err := exchangeInfoService.Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving Binance exchange info for symbol '%s': %v", symbol, err)
	}
	return exchangeInfo.Symbols[0]
}

// changedListings checks whether the listings on Binance have changed.
func (blc *BinanceListingsChecker) changedListings(oldAssets *[]string) (removed bool, changedAssets []string) {
	assets := blc.retrieveBinanceAssets()

	removed, changedAssets = utils.CompareLists(*oldAssets, assets)
	*oldAssets = assets

	return removed, changedAssets
}

// Start starts the BinanceListingsChecker.
func (blc *BinanceListingsChecker) Start(maxRate float64) {
	// Retrieve (old) stored Binance listings.
	oldAssets := utils.RetrieveOldListings()
	if len(oldAssets) == 0 { // Get from Binance if no old listings are stored.
		oldAssets = blc.retrieveBinanceAssets()
		utils.StoreOldListings(oldAssets)
	}

	// Check binance for new listings or de-listings and post Telegram/Discord message.
	limiter := rate.NewLimiter(rate.Limit(maxRate), 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.

		// Check for new listings or de-listings.
		removed, changedAssets := blc.changedListings(&oldAssets)

		// Post messages.
		for _, asset := range changedAssets {
			assetInfo := blc.retrieveSymbolInfo(asset)

			// Post telegram and discord messages.
			go messaging.SendAssetMessage(blc.TelegramBot, blc.TelegramChatID, blc.EnableTelegramMessage, blc.DiscordBot, blc.DiscordChannelIDs, blc.EnableDiscordMessages, removed, asset, assetInfo)

			utils.StoreOldListings(oldAssets)
		}
	}
}
