# Crypto Listing Sniper

A bot that watches several exchanges for new coin/token listings and posts a Telegram/Discord message. Currently supports the following exchanges:

- [Binance](https://www.binance.com/en)

## Features

- Posts a Discord/Telegram message when a new exchange listing is found.
- Posts a Discord/Telegram message when a new exchange announcement is published.
- Allows users to request the Telegram link using the Discord `/telegram-invite` slash command.
- Allows users to request the GitHub repo link using the Discord `/github-repo` slash command.

## How to use

1. Setup a discord application (see [this guide](https://discordjs.guide/preparations/setting-up-a-bot-application.html#what-is-a-token-anyway)).
2. Setup a telegram bot (see [this guide](https://telegrambots.github.io/book/1/quickstart.html)).
3. Install the Golang dependencies using `go get`.
4. Build the bot using `go build`
5. Rename the `.env.tamplate` file to `.env` and insert the required environmental variables.
6. Run the bot using `go run crypto-listings-sniper`.

## Contributing

Feel free to open an issue if you have ideas on how to make this repository better or if you want to report a bug! All contributions are welcome. :rocket: Please consult the [contribution guidelines](CONTRIBUTING.md) for more information.
