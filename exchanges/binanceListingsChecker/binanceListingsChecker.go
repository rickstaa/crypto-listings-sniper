// Description: Package binanceListingsChecker contains a class that when started periodically checks Binance for new listings or de-listings and posts a message in set message channels.
package binanceListingsChecker

import (
	"context"
	"log"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/common"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	"github.com/rickstaa/crypto-listings-sniper/messaging"
	"github.com/rickstaa/crypto-listings-sniper/utils"
	"golang.org/x/time/rate"
)

// BinanceListingsChecker is a class that when started checks Binance for new listings or de-listings and posts a message in set message channels
type BinanceListingsChecker struct {
	BinanceClient             *binance.Client
	TelegramBot               *telego.Bot
	TelegramChatID            int64
	EnableTelegramMessage     bool
	DiscordBot                *discordgo.Session
	DiscordChannelIDs         []string
	EnableDiscordMessages     bool
	OldAssets                 *[]string
	maxRate                   float64
	lastAssetsWarningTime     time.Time
	lastSymbolInfoWarningTime time.Time
}

// NewBinanceListingsChecker creates a new BinanceListingsChecker.
func NewBinanceListingsChecker(binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) *BinanceListingsChecker {
	return &BinanceListingsChecker{
		BinanceClient:             binanceClient,
		TelegramBot:               telegramBot,
		TelegramChatID:            telegramChatID,
		EnableTelegramMessage:     enableTelegramMessage,
		DiscordBot:                discordBot,
		DiscordChannelIDs:         discordChannelIDs,
		EnableDiscordMessages:     enableDiscordMessages,
		lastAssetsWarningTime:     time.Now(),
		lastSymbolInfoWarningTime: time.Now(),
	}
}

// retrieveBinanceAssets retrieves a list with the available assets from Binance.
// NOTE: Retry if failed and throw warning every minute.
func (blc *BinanceListingsChecker) retrieveBinanceAssets() (assets []string) {
	// Retrieve listing prices from Binance.
	priceService := blc.BinanceClient.NewListPricesService()
	listingPrices, err := priceService.Do(context.Background())
	assets = make([]string, len(listingPrices))

	// Log warning if failed.
	if err != nil {
		if time.Since(blc.lastAssetsWarningTime) > time.Minute { // Only log every minute.
			log.Printf("WARNING: Error retrieving Binance listing prices: %v", err)
			blc.lastAssetsWarningTime = time.Now()
		}
	}

	// Return assets.
	for i, s := range listingPrices {
		assets[i] = s.Symbol
	}
	return assets
}

// retrieveSymbolInfo retrieves information about a given symbol from Binance.
// NOTE: Try for 1 minutes before continuing.
func (blc *BinanceListingsChecker) retrieveSymbolInfo(symbol string) (assetInfo binance.Symbol) {
	tStart := time.Now()
	limiter := rate.NewLimiter(rate.Limit(blc.maxRate), 1)
	for time.Since(tStart) < 1*time.Minute {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.

		// Retrieve symbol info from Binance.
		exchangeInfoService := blc.BinanceClient.NewExchangeInfoService()
		exchangeInfoService = exchangeInfoService.Symbols(symbol)
		exchangeInfoTmp, err := exchangeInfoService.Do(context.Background())

		// Retry if symbol was not found (i.e. error -1121).
		if err != nil {
			if err.(*common.APIError).Code == -1121 {
				if time.Since(blc.lastSymbolInfoWarningTime) > 10*time.Second { // Only log every 10 seconds.
					log.Printf("Warning: Error retrieving Binance symbol info: %v", err)
					blc.lastSymbolInfoWarningTime = time.Now()
				}
				continue
			}
			log.Fatalf("WARNING: Error retrieving Binance symbol info: %v", err)
		}

		assetInfo = exchangeInfoTmp.Symbols[0]
		break
	}

	return assetInfo
}

// changedListings checks whether the listings on Binance have changed.
func (blc *BinanceListingsChecker) changedListings(oldAssets *[]string) (removed bool, changedAssets []string) {
	assets := blc.retrieveBinanceAssets()

	// Return if no assets are available.
	if len(assets) == 0 {
		return false, []string{}
	}

	// Return changed assets.
	removed, changedAssets = utils.CompareLists(*oldAssets, assets)
	*oldAssets = assets
	return removed, changedAssets
}

// Post messages in Telegram and Discord if new listings or de-listings are found.
func (blc *BinanceListingsChecker) postMessages(removed bool, changedAssets []string, oldAssets []string) {
	for _, asset := range changedAssets {
		assetInfo := blc.retrieveSymbolInfo(asset)

		// Log new listing or de-listing.
		if removed {
			log.Printf("De-listing found: %v", asset)
		} else {
			log.Printf("New listing found: %v", asset)
		}

		// Post telegram and discord messages.
		go messaging.SendAssetMessage(blc.TelegramBot, blc.TelegramChatID, blc.EnableTelegramMessage, blc.DiscordBot, blc.DiscordChannelIDs, blc.EnableDiscordMessages, removed, asset, assetInfo)

		utils.StoreOldListings(oldAssets)
	}
}

// Start starts the BinanceListingsChecker.
func (blc *BinanceListingsChecker) Start(maxRate float64) {
	blc.maxRate = maxRate

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
		go blc.postMessages(removed, changedAssets, oldAssets)
	}
}
