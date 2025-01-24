package main

import (
	"log"
	tgCLient "storage-links-bot/clients/telegram"
	"storage-links-bot/conf"
	"storage-links-bot/consumer/eventconsumer"
	"storage-links-bot/events/telegram"
	"storage-links-bot/storage/files"
)

var cfg *conf.Configuration

func init() {
	cfg = conf.Compile()
}

func main() {
	tg := tgCLient.New(cfg.ApiBot, cfg.Token)
	path := files.New(cfg.StoragePath)
	eventProcessor := telegram.New(tg, path)
	log.Printf("[INF] storage-link-bot started")
	consumer := eventсonsumer.New(eventProcessor, eventProcessor, cfg.BatchSize)
	consumer.Start()
}
