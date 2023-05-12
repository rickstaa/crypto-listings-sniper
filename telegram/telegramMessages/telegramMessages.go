// Description: This package contains functions for creating Telegram messages.
package telegramMessages

import (
	"fmt"
)

// Returns a string containing a message for a new asset.
func newAssetMessage(asset string, url string) string {
	return fmt.Sprintf("💎 <u>Binance listed new asset (<a href='%s'>%s</a>)</u>", url, asset)
}

// Return a string containing a message for a removed asset.
func removedAssetMessage(asset string) string {
	return fmt.Sprintf("🗑 <u>Binance removed asset (%s)</u>\n", asset)
}

// Returns a string containing a Telegram message for a new or removed asset.
func AssetMessage(removed bool, asset string, url string) string {
	switch removed {
	case true:
		return removedAssetMessage(asset)
	default:
		return newAssetMessage(asset, url)
	}
}
