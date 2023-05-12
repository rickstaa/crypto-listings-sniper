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
	binanceKey, binanceSecret, telegramBotKey, telegramChatID, enableTelegramMessage, discordBotKey, discordChannelIDs, discordAppID, enableDiscordMessages := utils.GetEnvVars()

	// Load Telegram bot.
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
	defer discordBot.Close()

	// Log Telegram bot and channel info.
	telegramBotInfo, err := telegramBot.GetMe()
	if err != nil {
		log.Fatalf("Error getting telegramBot info: %v", err)
	}
	telegramChat, err := telegramBot.GetChat(&telego.GetChatParams{ChatID: tu.ID(telegramChatID)})
	if err != nil {
		log.Fatalf("Error getting telegramChat info: %v", err)
	}
	log.Printf("Authorized on account: %s", telegramBotInfo.Username)
	log.Printf("Bot id: %d", telegramBotInfo.ID)
	log.Printf("Chat id: %d", telegramChatID)
	log.Printf("Chat type: %s", telegramChat.Type)
	log.Printf("Chat title: %s", telegramChat.Title)
	log.Printf("Chat username: %s", telegramChat.Username)
	log.Printf("Chat description: %s", telegramChat.Description)

	// Register slash commands.
	dc.SetupDiscordSlashCommands(discordBot, discordAppID, telegramChat.InviteLink)

	// Load Binance client.
	binanceClient := binance.NewClient(binanceKey, binanceSecret)
	log.Printf("Binance API endpoint: %s", binanceClient.BaseURL)

	// Retrieve old assets and store them if they do not exist.
	oldAssets := utils.RetrieveOldListings()
	if len(oldAssets) == 0 {
		assets := bn.RetrieveBinanceAssets(binanceClient)

		utils.StoreOldListings(assets)
	}

	// Check binance for new listings or de-listings and post Telegram/Discord message.
	r := rate.Every(1 * time.Millisecond)
	limiter := rate.NewLimiter(r, 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.
		bn.BinanceListingsCheck(&oldAssets, binanceClient, telegramBot, telegramChatID, enableTelegramMessage, discordBot, discordChannelIDs, enableDiscordMessages)
	}
}
