package telegram

import (
	"log"
	"net/url"
	"storage-links-bot/clients/telegram"
	"storage-links-bot/lib/e"
	"storage-links-bot/storage"
	"strconv"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	AllCmd   = "/all"
	ClearCmd = "/clear"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	command := strings.TrimSpace(text)
	log.Printf("%s do command %s\n", username, text)
	if isAddCmd(command) {
		return p.savePage(text, chatId, username)
	}
	switch command {
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	case RndCmd:
		return p.sendRandom(chatId, username)
	case AllCmd:
		return p.sendAll(chatId, username)
	case ClearCmd:
		return p.sendClear(chatId, username)
	default:
		return newSendMsg(chatId, p.tg)(msgUnknownCommand)
	}
}

func (p *Processor) savePage(path string, chatId int, username string) error {
	senderMsg := newSendMsg(chatId, p.tg)
	page := storage.Page{
		URL:      path,
		Username: username,
		ChatId:   strconv.Itoa(chatId),
	}
	exist, err := p.storage.IsExist(&page)
	if err != nil {
		return err
	}
	if exist {
		return senderMsg(msgAlreadyExists)
	}
	if err := p.storage.Save(&page); err != nil {
		return err
	}
	return senderMsg(msgSaved)
}

func (p *Processor) sendRandom(chatId int, username string) error {
	senderMsg := newSendMsg(chatId, p.tg)
	page, err := p.storage.PickRandom(username)
	if err != nil {
		log.Printf("Error getting random page %s\n", err)
		return senderMsg(msgNoSavedPages)
	}
	if err := senderMsg(page.URL); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendAll(chatId int, username string) error {
	senderMsg := newSendMsg(chatId, p.tg)
	result, err := p.storage.PickAll(strconv.Itoa(chatId))
	if err != nil {
		log.Printf("Error getting random page %s\n", err)
		return senderMsg(msgNoSavedPages)
	}
	if err := senderMsg(strings.Join(result, ",")); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendClear(chatId int, username string) error {
	senderMsg := newSendMsg(chatId, p.tg)
	err := p.storage.Remove(&storage.Page{
		URL:      "",
		Username: username,
		ChatId:   strconv.Itoa(chatId),
	})
	if err != nil {
		return e.Wrap("Bad request to clean", err)
	}
	return senderMsg(msgClearStorage)
}

func (p *Processor) sendHelp(chatId int) error {
	return newSendMsg(chatId, p.tg)(msgHelp)
}

func (p *Processor) sendHello(chatId int) error {
	return newSendMsg(chatId, p.tg)(msgHello)
}

func newSendMsg(chatId int, tg *telegram.Client) func(string) error {
	return func(text string) error {
		return tg.SendMessage(chatId, text)
	}
}

func isAddCmd(cmd string) bool {
	return isUrl(cmd)
}

func isUrl(path string) bool {
	u, err := url.Parse(path)
	return err == nil && u.Host != ""
}
