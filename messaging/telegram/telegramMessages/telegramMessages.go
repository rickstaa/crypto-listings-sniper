// Description: The telegramMessages package contains functions for creating Telegram messages.
package telegramMessages

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
)

// newAssetMessage returns a new asset Telegram message.
func newAssetMessage(asset string, url string, symbolInfo binance.Symbol) string {
	return fmt.Sprintf("💎 <u>Binance listed new asset (<a href='%s'>%s</a>)</u>\n\n", url, asset) +
		fmt.Sprintf("- <b>Base Asset:</b> %s\n", symbolInfo.BaseAsset) +
		fmt.Sprintf("- <b>Quota Asset:</b> %s\n", symbolInfo.QuoteAsset)
}

// removedAssetMessage return a removed asset Telegram message.
func removedAssetMessage(asset string) string {
	return fmt.Sprintf("🗑 <u>Binance removed asset (%s)</u>\n", asset)
}

// AssetMessage returns a string containing new/removed asset Telegram message.
func AssetMessage(removed bool, asset string, url string, assetInfo binance.Symbol) string {
	if removed {
		return removedAssetMessage(asset)
	}
	return newAssetMessage(asset, url, assetInfo)
}

// Returns a string containing a message for a new announcement.
func AnnouncementMessage(url string, title string) string {
	return fmt.Sprintf("📢 <a href='%s'>%s</a>\n", url, title)
}
