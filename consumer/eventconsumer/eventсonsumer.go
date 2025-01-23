package eventсonsumer

import (
	"log"
	"storage-links-bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() {
	for {
		fetch, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer %s , %v", err.Error(), err)
			continue
		}
		if len(fetch) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		c.handleEvents(fetch)
	}
}

func (c *Consumer) handleEvents(events []events.Event) {
	for _, event := range events {
		err := c.processor.Process(event)
		log.Printf("start process event: %v", event)
		if err != nil {
			log.Printf("[ERR] can't handle event %s", err.Error())
			continue
		}
	}
}
