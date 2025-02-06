package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	"storage-links-bot/lib/e"
)

const (
	errHashField = "hash error from field"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(id string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
	PickAll(id string) ([]string, error)
}

type Page struct {
	URL      string
	Username string
	ChatId   string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	_, errUrl := io.WriteString(h, p.URL)
	if errUrl != nil {
		return "", e.Wrap(errHashField, errUrl)
	}
	_, errName := io.WriteString(h, p.Username)
	if errName != nil {
		return "", e.Wrap(errHashField, errName)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
