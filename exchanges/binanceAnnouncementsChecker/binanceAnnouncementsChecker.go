// Description: Package binanceAnnouncementsChecker contains a class that when started checks the unofficial Binance announcements api for new announcements.
package binanceAnnouncementsChecker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/mymmrac/telego"
	"github.com/rickstaa/crypto-listings-sniper/messaging"
	"github.com/rickstaa/crypto-listings-sniper/utils"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/maps"
	"golang.org/x/time/rate"
)

// GetBinanceAnnouncementsEndpoint returns the (unofficial) binance announcements endpoint.
// NOTE: Retrieved from https://stackoverflow.com/a/69673063/8135687.
func GetBinanceAnnouncementsEndpoint() string {
	queries := map[string]string{
		"catalogId": "48",
		"pageNo":    "1",
		"pageSize":  fmt.Sprintf("%d", rand.Intn(50-10)+10),
	}
	var url strings.Builder
	url.WriteString("https://www.binance.com/bapi/composite/v1/public/cms/article/catalog/list/query?")

	// Add queries to url.
	for key, value := range queries {
		url.WriteString(key + "=" + value + "&")
	}
	urlString := url.String()
	urlString = urlString[:len(urlString)-1]

	return urlString
}

// BinanceAnnouncement represents the Binance announcement.
type BinanceAnnouncements struct {
	Code          string                   `json:"code"`
	Message       string                   `json:"message"`
	MessageDetail string                   `json:"messageDetail"`
	Data          BinanceAnnouncementsData `json:"data"`
	Success       bool                     `json:"success"`
}

// BinanceData represents the Binance announcements data.
type BinanceAnnouncementsData struct {
	Articles []BinanceArticle `json:"articles"`
	Total    int64            `json:"total"`
}

// BinanceArticle represents a Binance article.
type BinanceArticle struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Title       string `json:"title"`
	ImageLink   string `json:"imageLink"`
	ShortLink   string `json:"shortLink"`
	Body        string `json:"body"`
	Type        string `json:"type"`
	CatalogId   int64  `json:"catalogId"`
	CatalogName string `json:"catalogName"`
	PublishDate string `json:"publishDate"`
	Footer      string `json:"footer"`
}

// BinanceAnnouncementsChecker is a class that when started checks Binance for new announcements and posts a message in set message channels
type BinanceAnnouncementsChecker struct {
	binanceClient         *binance.Client
	telegramBot           *telego.Bot
	telegramChatID        int64
	enableTelegramMessage bool
	discordBot            *discordgo.Session
	discordChannelIDs     []string
	enableDiscordMessages bool
	lastCheckTime         time.Time
}

// newBinanceAnnouncementsChecker creates a new BinanceAnnouncementsChecker.
func NewBinanceAnnouncementsChecker(binanceClient *binance.Client, telegramBot *telego.Bot, telegramChatID int64, enableTelegramMessage bool, discordBot *discordgo.Session, discordChannelIDs []string, enableDiscordMessages bool) *BinanceAnnouncementsChecker {
	return &BinanceAnnouncementsChecker{
		binanceClient:         binanceClient,
		telegramBot:           telegramBot,
		telegramChatID:        telegramChatID,
		enableTelegramMessage: enableTelegramMessage,
		discordBot:            discordBot,
		discordChannelIDs:     discordChannelIDs,
		enableDiscordMessages: enableDiscordMessages,
	}
}

// Retrieves the Binance announcements from the Binance announcements endpoint.
func (blc *BinanceAnnouncementsChecker) retrieveBinanceAnnouncements() (binanceAnnouncements map[string]string) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	// Make request.
	request.SetRequestURI(GetBinanceAnnouncementsEndpoint())
	request.Header.SetMethod("GET")
	request.Header.Set("Content-Type", "application/json")
	err := fasthttp.Do(request, response)
	if err != nil {
		log.Fatalf("Error scraping binance announcements endpoint: %v", err)
	}
	if response.StatusCode() != 200 {
		if time.Since(blc.lastCheckTime) > time.Minute { // Only log every minute.
			log.Printf("WARNING: Announcement API endpoint not responding.")
			blc.lastCheckTime = time.Now()
		}

		return binanceAnnouncements
	}
	blc.lastCheckTime = time.Now()

	// Unmarshal response.
	var announcements BinanceAnnouncements
	err = json.Unmarshal(response.Body(), &announcements)
	if err != nil {
		log.Fatalf("Error unmarshalling binance announcements response: %v", err)
	}

	// Return last 10 announcements.
	binanceAnnouncements = make(map[string]string)
	for _, article := range announcements.Data.Articles[:10] {
		binanceAnnouncements[article.Code] = article.Title
	}
	return binanceAnnouncements
}

// changedListings checks whether new announcements have been published on Binance.
func (blc *BinanceAnnouncementsChecker) binanceAnnouncementsCheck(oldAnnouncementsCodes *[]string) (newAnnouncementsCodes []string, newAnnouncements map[string]string) {
	announcements := blc.retrieveBinanceAnnouncements()
	if len(announcements) == 0 {
		return nil, nil
	}

	// Check if new announcements have been published.
	announcementsCodes := maps.Keys(announcements)
	_, newAnnouncementsCodes = utils.CompareLists(*oldAnnouncementsCodes, announcementsCodes)
	*oldAnnouncementsCodes = announcementsCodes

	// Create new announcements map.
	newAnnouncements = make(map[string]string)
	for _, code := range newAnnouncementsCodes {
		newAnnouncements[code] = announcements[code]
	}

	return newAnnouncementsCodes, newAnnouncements
}

// Start starts the BinanceAnnouncementsChecker.
func (blc *BinanceAnnouncementsChecker) Start(maxRate float64) {
	// Retrieve (old) Binance announcements.
	oldAnnouncements := utils.RetrieveOldAnnouncements()
	if len(oldAnnouncements) == 0 { // Get from Binance if no old announcements are stored.
		binanceAnnouncements := blc.retrieveBinanceAnnouncements()
		oldAnnouncements = maps.Keys(binanceAnnouncements)
		utils.StoreOldAnnouncements(oldAnnouncements)
	}

	// Check binance for new announcements and post Telegram/Discord message.
	limiter := rate.NewLimiter(rate.Limit(maxRate), 1)
	for {
		limiter.Wait(context.Background()) // NOTE: This is to prevent binance from blocking the IP address.

		// Check for new announcements.
		newAnnouncementsCodes, newAnnouncements := blc.binanceAnnouncementsCheck(&oldAnnouncements)

		// Post messages.
		for _, announcementCode := range newAnnouncementsCodes {
			// Log announcement.
			log.Printf("New Binance announcement: %s", newAnnouncements[announcementCode])

			// Post telegram and discord messages.
			go messaging.SendAnnouncementMessage(blc.telegramBot, blc.telegramChatID, blc.enableTelegramMessage, blc.discordBot, blc.discordChannelIDs, blc.enableDiscordMessages, announcementCode, newAnnouncements[announcementCode])

			utils.StoreOldAnnouncements(oldAnnouncements)
		}
	}
}
