package main

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/exp/maps"
	"golang.org/x/time/rate"

	"github.com/rickstaa/crypto-listings-sniper/exchanges/binanceAnnouncementsChecker"
	"github.com/rickstaa/crypto-listings-sniper/exchanges/binanceListingsChecker"
	"github.com/rickstaa/crypto-listings-sniper/messaging"
	dc "github.com/rickstaa/crypto-listings-sniper/messaging/discord"
	"github.com/rickstaa/crypto-listings-sniper/utils"

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
	binanceClient.SetApiEndpoint("https://api4.binance.com")
	log.Printf("Binance API endpoint: %s", binanceClient.BaseURL)

	// Initialize checkers.
	binanceListingsChecker := binanceListingsChecker.NewBinanceListingsChecker(binanceClient, telegramBot, telegramChatID, enableTelegramMessage, discordBot, discordChannelIDs, enableDiscordMessages)

	// start the checkers.
	go binanceListingsChecker.Start()

	// Check Binance for new announcements and post Telegram/Discord message.
	go BinanceAnnouncementsChecker(binanceClient, telegramBot, telegramChatID, enableTelegramMessage, discordBot, discordChannelIDs, enableDiscordMessages)

	runtime.Goexit() // Keep the program running.
}

// BinanceAnnouncementsChecker periodically checks Binance for new announcements and post Telegram/Discord message.
func BinanceAnnouncementsChecker(binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) {

	// Retrieve latest Binance announcements codes.
	oldAnnouncements := utils.RetrieveOldAnnouncements()
	if len(oldAnnouncements) == 0 {
		// Retrieve latest Binance announcement and store them if no old announcements are stored.
		binanceAnnouncements := binanceAnnouncementsChecker.RetrieveLatestBinanceAnnouncements()
		oldAnnouncements = maps.Keys(binanceAnnouncements)
		utils.StoreOldAnnouncements(oldAnnouncements)
	}

	r := rate.Every(1 * time.Second)
	limiter := rate.NewLimiter(r, 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.

		// Retrieve latest Binance announcement.
		announcements := binanceAnnouncementsChecker.RetrieveLatestBinanceAnnouncements()

		// Retrieve codes from announcement map
		announcementCodes := maps.Keys(announcements)

		_, newAnnouncements := utils.CompareLists(oldAnnouncements, announcementCodes)

		if len(newAnnouncements) > 0 {
			oldAnnouncements = announcementCodes

			// Send Telegram/Discord message.
			for _, announcementCode := range newAnnouncements {
				go messaging.SendAnnouncementMessage(telegramBot, telegramChatID, enableTelegramMessage, discordBot, discordChannelIDs, enableDiscordMessages, announcementCode, announcements[announcementCode])
			}
		}
	}
}
