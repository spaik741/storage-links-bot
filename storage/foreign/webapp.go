package foreign

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"storage-links-bot/clients/dto"
	"storage-links-bot/lib/e"
	"storage-links-bot/storage"
	"time"
)

const (
	basePath = "links"
)

type RestTemplate struct {
	Client http.Client
	Host   string
}

func NewRestTemplate(host string) *RestTemplate {
	return &RestTemplate{
		Client: http.Client{Timeout: time.Second * 3},
		Host:   host,
	}
}

func (r *RestTemplate) Save(p *storage.Page) error {
	var buf bytes.Buffer
	linkRequest := dto.LinksRequest{
		ChatId: p.ChatId,
		Link:   p.URL,
	}
	if err := json.NewEncoder(&buf).Encode(linkRequest); err != nil {
		return err
	}
	u, _ := url.ParseRequestURI(r.Host)
	u.Path = path.Join(basePath)
	req, err := http.NewRequest(http.MethodPost, u.String(), &buf)
	if err != nil {
		return e.Wrap("Error create request", err)
	}
	response, err := r.Client.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		return e.Wrap("Fail request to storage", err)
	}
	return nil
}

func (r *RestTemplate) PickRandom(id string) (*storage.Page, error) {
	urls, err := r.PickAll(id)
	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(urls))
	result := urls[n]
	return &storage.Page{
		URL:      result,
		Username: "",
		ChatId:   id,
	}, nil
}

func (r *RestTemplate) Remove(p *storage.Page) error {
	u, _ := url.ParseRequestURI(r.Host)
	params := url.Values{}
	params.Add("chatId", p.ChatId)
	req, err := http.NewRequest(http.MethodDelete, u.String(), nil)
	req.URL.RawQuery = params.Encode()
	if err != nil {
		return e.Wrap("Error create request", err)
	}
	response, err := r.Client.Do(req)
	if err != nil || response.StatusCode != http.StatusOK {
		return e.Wrap("Fail request to storage", err)
	}
	return nil
}
func (r *RestTemplate) IsExist(p *storage.Page) (bool, error) {
	return false, nil
}

func (r *RestTemplate) PickAll(id string) ([]string, error) {
	u, _ := url.ParseRequestURI(r.Host)
	u.Path = path.Join(basePath, id)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	var result dto.LinksResponse
	response, err := r.Client.Do(req)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Links, nil
}
