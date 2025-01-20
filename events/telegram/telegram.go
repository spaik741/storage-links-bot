package telegram

import (
	"errors"
	"storage-links-bot/clients/telegram"
	"storage-links-bot/events"
	"storage-links-bot/storage"
)

const (
	eventErr = "unknown type of event"
	metaErr  = "unknown type of meta"
)

type (
	Processor struct {
		tg      *telegram.Client
		offset  int
		storage storage.Storage
	}

	Meta struct {
		ChatID   int
		Username string
	}
)

func New(tg *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{tg: tg, offset: 0, storage: storage}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return nil, nil
	}
	result := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		result = append(result, mapEvent(update))
	}
	p.offset = updates[len(updates)-1].Id + 1
	return result, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.New(eventErr)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := convertMeta(event)
	if err != nil {
		return err
	}
	return p.doCmd(event.Text, meta.ChatID, meta.Username)
}

func convertMeta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if !ok {
		return *new(Meta), errors.New(metaErr)
	}
	return meta, nil
}

func mapEvent(update telegram.Update) events.Event {
	updType := fetchType(update)
	event := events.Event{
		Type: updType,
		Text: fetchText(update),
	}
	if updType == events.Message {
		event.Meta = Meta{
			ChatID:   update.Message.Chat.Id,
			Username: update.Message.From.Username,
		}
	}
	return event
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
