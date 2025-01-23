package telegram

import (
	"log"
	"net/url"
	"storage-links-bot/clients/telegram"
	"storage-links-bot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
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
	default:
		return newSendMsg(chatId, p.tg)(msgUnknownCommand)
	}
}

func (p *Processor) savePage(path string, chatId int, username string) error {
	senderMsg := newSendMsg(chatId, p.tg)
	page := storage.Page{
		URL:      path,
		Username: username,
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
	err = senderMsg(msgSaved)
	if err != nil {
		return err
	}
	return nil
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

	return p.storage.Remove(page)
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
