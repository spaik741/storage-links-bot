package main

import (
	"flag"
	"log"
	"storage-links-bot/clients/telegram"
)

const (
	tgBot = "api.telegram.org"
)

func main() {
	telegram.New(tgBot, mustToken())
}

func mustToken() string {
	token := flag.String("token-bot", "", "Give me your token for bot: ")
	flag.Parse()
	if *token == "" {
		log.Fatal("You need to provide your token")
	}
	return *token
}
