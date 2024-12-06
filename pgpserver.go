package main

import (
	"flag"
	"pgpserver/internal/bot"
)

var BotToken string

func init() {
    flag.StringVar(&BotToken, "bottoken", "", "discord bot token")

    flag.Parse()
}

func main() {
    bot.Run(BotToken)
}
