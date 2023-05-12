// Description: This package contains functions for creating Telegram messages.
package telegramMessages

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
)

// Returns a string containing a message for a new asset.
func newAssetMessage(asset string, url string, symbolInfo binance.Symbol) string {
	return fmt.Sprintf("💎 <u>Binance listed new asset (<a href='%s'>%s</a>)</u>\n\n", url, asset) +
		fmt.Sprintf("- <b>Base Asset:</b> %s\n", symbolInfo.BaseAsset) +
		fmt.Sprintf("- <b>Quota Asset:</b> %s\n", symbolInfo.QuoteAsset)
}

// Return a string containing a message for a removed asset.
func removedAssetMessage(asset string) string {
	return fmt.Sprintf("🗑 <u>Binance removed asset (%s)</u>\n", asset)
}

// Returns a string containing a Telegram message for a new or removed asset.
func AssetMessage(removed bool, asset string, url string, symbolInfo binance.Symbol) string {
	switch removed {
	case true:
		return removedAssetMessage(asset)
	default:
		return newAssetMessage(asset, url, symbolInfo)
	}
}
