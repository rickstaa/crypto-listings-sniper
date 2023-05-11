package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/profile"
	"github.com/rickstaa/crypto-listings-sniper/utils"
	"github.com/rickstaa/crypto-listings-sniper/utils/checkers"
	"golang.org/x/time/rate"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: Run in loop/
// TODO: Create discord bot.
// TODO: Create telegram link commmand.
// TODO: Remove logging statements in checkers.go.

func main() {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	telegramBotKey, chatID, binanceKey, binanceSecret := utils.GetEnvVars()

	// Load Telegram telegramBot.
	telegramBot, err := telego.NewBot(telegramBotKey)
	if err != nil {
		log.Fatalf("Error loading Telegram telegramBot: %v", err)
	}

	// Log telegramBot and channel info.
	botInfo, err := telegramBot.GetMe()
	if err != nil {
		log.Fatalf("Error getting telegramBot info: %v", err)
	}
	chat, err := telegramBot.GetChat(&telego.GetChatParams{ChatID: tu.ID(chatID)})
	if err != nil {
		log.Fatalf("Error getting chat info: %v", err)
	}
	log.Printf("Authorized on account: %s", botInfo.Username)
	log.Printf("Bot id: %d", botInfo.ID)
	log.Printf("Chat id: %d", chatID)
	log.Printf("Chat type: %s", chat.Type)
	log.Printf("Chat title: %s", chat.Title)
	log.Printf("Chat username: %s", chat.Username)
	log.Printf("Chat description: %s", chat.Description)

	// Load Binance binanceClient.
	binanceClient := binance.NewClient(binanceKey, binanceSecret)
	log.Printf("Binance API endpoint: %s", binanceClient.BaseURL)

	// Retrieve old SPOT base assets and symbols and store them if they do not exist.
	oldBaseAssetsList, oldSymbolsList := utils.RetrieveOldListings()
	if len(oldBaseAssetsList) == 0 || len(oldSymbolsList) == 0 {
		baseAssets, symbols, _ := utils.RetrieveBinanceSpotAssets(binanceClient)

		// Store symbol and base assets lists in JSON files.
		utils.StoreOldListings(baseAssets, symbols)
	}

	// Check binance for new SPOT listings or de-listings and post telegram message.
	r := rate.Every(1 * time.Millisecond)
	limiter := rate.NewLimiter(r, 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.
		checkers.BinanceListingsCheck(&oldBaseAssetsList, &oldSymbolsList, binanceClient, telegramBot, chatID)
	}
}
