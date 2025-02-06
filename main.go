package main

import (
	"log"
	tgCLient "storage-links-bot/clients/telegram"
	"storage-links-bot/conf"
	"storage-links-bot/consumer/eventconsumer"
	"storage-links-bot/events/telegram"
	"storage-links-bot/storage"
	"storage-links-bot/storage/files"
	"storage-links-bot/storage/foreign"
	"strings"
)

var cfg *conf.Configuration

func init() {
	cfg = conf.Compile()
}

func main() {
	tg := tgCLient.New(cfg.ApiBot, cfg.Token)
	var storageImpl storage.Storage
	if strings.HasPrefix(cfg.StoragePath, "http") {
		storageImpl = foreign.NewRestTemplate(cfg.StoragePath)
	} else {
		storageImpl = files.New(cfg.StoragePath)
	}
	eventProcessor := telegram.New(tg, storageImpl)
	log.Printf("[INF] storage-link-bot started")
	consumer := eventсonsumer.New(eventProcessor, eventProcessor, cfg.BatchSize)
	consumer.Start()
}
