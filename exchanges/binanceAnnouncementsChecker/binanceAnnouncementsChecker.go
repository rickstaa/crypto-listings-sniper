// Description: Package binanceAnnouncementsChecker contains a class that when started checks the unofficial Binance announcements api for new announcements.
// NOTE: Retrieved from https://stackoverflow.com/a/69673063/8135687.
package binanceAnnouncementsChecker

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

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

// BinanceAnnouncementsEndpoint returns the binance announcements endpoint while making sure that the endpoint is not cached.
func BinanceAnnouncementsEndpoint() string {
	queries := map[string]string{
		"catalogId": "48",
		"pageNo":    "1",
		"pageSize":  "10", // NOTE: This is to prevent the endpoint from being cached.
	}
	var url strings.Builder
	url.WriteString("https://www.binance.com/bapi/composite/v1/public/cms/article/catalog/list/query?")
	for key, value := range queries {
		url.WriteString(key + "=" + value + "&")
	}
	return url.String()
}

// Retrieves the Latest Binance announcement from the hidden Binance announcements endpoint.
// NOTE: Endpoint https://stackoverflow.com/a/69673063/8135687.
func RetrieveLatestBinanceAnnouncements() (binanceAnnouncements map[string]string) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)

	// Make request.
	request.SetRequestURI(BinanceAnnouncementsEndpoint())
	request.Header.SetMethod("GET")
	request.Header.Set("Content-Type", "application/json")

	// Retrieve response.
	err := fasthttp.Do(request, response)
	if err != nil {
		log.Fatalf("Error scraping binance announcements endpoint: %v", err)
	}

	// Unmarshal response.
	var announcements BinanceAnnouncements
	err = json.Unmarshal(response.Body(), &announcements)
	if err != nil {
		log.Fatalf("Error unmarshalling binance announcements response: %v", err)
	}

	// Store last 10 announcement in map.
	binanceAnnouncements = make(map[string]string)
	for _, article := range announcements.Data.Articles[:10] {
		binanceAnnouncements[article.Code] = article.Title
	}

	return binanceAnnouncements
}
