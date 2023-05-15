# Crypto Listing Sniper

A bot that watches for new coin/token listings on Binance and posts a Telegram/Discord message.

## Features

- Posts a Discord/Telegram message when a new Binance listing is found.
- Posts a Discord/Telegram message when a new [Binance listings announcement](https://www.binance.com/en/support/announcement/new-cryptocurrency-listing?c=48) is published.
- Allows users to request the Telegram link using the Discord `/telegram-invite` slash command.

## How to use

1. Install the Golang dependencies using `go get`.
2. Build the bot using `go build`
3. Rename the `.env.tamplate` file to `.env` and insert the required environmental variables.
4. Run the bot using `go run crypto-listings-sniper`.
