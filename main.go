package main

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	bn "github.com/rickstaa/crypto-listings-sniper/binance"
	dc "github.com/rickstaa/crypto-listings-sniper/discord"
	"github.com/rickstaa/crypto-listings-sniper/utils"
	"golang.org/x/time/rate"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	telegramBotKey, chatID, binanceKey, binanceSecret, discordBotKey, discordChannelIDs, discordAppID := utils.GetEnvVars()

	// Load Telegram telegramBot.
	telegramBot, err := telego.NewBot(telegramBotKey)
	if err != nil {
		log.Fatalf("Error loading Telegram telegramBot: %v", err)
	}
	defer telegramBot.Close()

	// Load Discord bot.
	discordBot, err := discordgo.New("Bot " + discordBotKey)
	if err != nil {
		log.Fatalf("Error loading Discord bot: %v", err)
	}

	// Log telegramBot and channel info.
	telegramBotInfo, err := telegramBot.GetMe()
	if err != nil {
		log.Fatalf("Error getting telegramBot info: %v", err)
	}
	telegramChat, err := telegramBot.GetChat(&telego.GetChatParams{ChatID: tu.ID(chatID)})
	if err != nil {
		log.Fatalf("Error getting telegramChat info: %v", err)
	}
	log.Printf("Authorized on account: %s", telegramBotInfo.Username)
	log.Printf("Bot id: %d", telegramBotInfo.ID)
	log.Printf("Chat id: %d", chatID)
	log.Printf("Chat type: %s", telegramChat.Type)
	log.Printf("Chat title: %s", telegramChat.Title)
	log.Printf("Chat username: %s", telegramChat.Username)
	log.Printf("Chat description: %s", telegramChat.Description)

	// Register slash commands.
	dc.SetupDiscordSlashCommands(discordBot, discordAppID, telegramChat.InviteLink)
	discordBot.Open()
	defer discordBot.Close()

	// Load Binance binanceClient.
	binanceClient := binance.NewClient(binanceKey, binanceSecret)
	log.Printf("Binance API endpoint: %s", binanceClient.BaseURL)

	// Retrieve old SPOT base assets and symbols and store them if they do not exist.
	oldBaseAssetsList, oldSymbolsList := utils.RetrieveOldListings()
	if len(oldBaseAssetsList) == 0 || len(oldSymbolsList) == 0 {
		baseAssets, symbols, _ := bn.RetrieveBinanceSpotAssets(binanceClient)

		// Store symbol and base assets lists in JSON files.
		utils.StoreOldListings(baseAssets, symbols)
	}

	// Check binance for new SPOT listings or de-listings and post telegram message.
	r := rate.Every(1 * time.Millisecond)
	limiter := rate.NewLimiter(r, 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.
		bn.BinanceListingsCheck(&oldBaseAssetsList, &oldSymbolsList, binanceClient, telegramBot, chatID, discordBot, discordChannelIDs)
	}
}
