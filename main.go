package main

import (
	"log"
	"runtime"

	"github.com/bwmarrin/discordgo"
	bn "github.com/rickstaa/crypto-listings-sniper/exchanges/binance"
	dc "github.com/rickstaa/crypto-listings-sniper/messaging/discord"
	"github.com/rickstaa/crypto-listings-sniper/utils"

	"github.com/adshao/go-binance/v2"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// TODO: Create tests.
// TODO: Cleanup code.

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
	binanceClient.SetApiEndpoint("https://api4.binance.com")
	log.Printf("Binance API endpoint: %s", binanceClient.BaseURL)

	// Initialize checkers.
	binanceListingsChecker := bn.NewBinanceListingsChecker(binanceClient, telegramBot, telegramChatID, enableTelegramMessage, discordBot, discordChannelIDs, enableDiscordMessages)

	// start the checkers.
	go binanceListingsChecker.Start()

	runtime.Goexit() // Keep the program running.
}
