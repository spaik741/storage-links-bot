package main

import (
	"flag"
	"log"
	tgCLient "storage-links-bot/clients/telegram"
	"storage-links-bot/consumer/eventconsumer"
	"storage-links-bot/events/telegram"
	"storage-links-bot/storage/files"
)

const (
	tgBot     = "api.telegram.org"
	storage   = "files"
	batchSize = 100
)

func main() {
	tg := tgCLient.New(tgBot, mustToken())
	path := files.New(storage)
	eventProcessor := telegram.New(tg, path)
	log.Printf("[INF] storage-link-bot started")
	consumer := eventсonsumer.New(eventProcessor, eventProcessor, batchSize)
	consumer.Start()
}

func mustToken() string {
	token := flag.String("token-bot", "", "Give me your token for bot: ")
	flag.Parse()
	if *token == "" {
		log.Fatal("You need to provide your token")
	}
	return *token
}
