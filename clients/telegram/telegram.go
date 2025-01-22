package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"storage-links-bot/lib/e"
	"strconv"
)

const (
	prefixBp      string = "bot"
	reqErr               = "can't do request"
	readErr              = "can't read response"
	parseErr             = "can't parse response"
	getUpdatesMtd        = "getUpdates"
	sendMsgMtd           = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: createBasePath(token),
		client:   http.Client{},
	}
}

func (c *Client) Updates(offset, limit int) ([]Update, error) {
	params := url.Values{}
	params.Add("offset", strconv.Itoa(offset))
	params.Add("limit", strconv.Itoa(limit))
	request, err := c.doRequest(getUpdatesMtd, params)
	if err != nil {
		return nil, err
	}
	var updates []Update
	err = json.Unmarshal(request, &updates)
	if err != nil {
		return nil, e.Wrap(parseErr, err)
	}
	return updates, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	params := url.Values{}
	params.Add("chat_id", strconv.Itoa(chatId))
	params.Add("text", text)
	_, err := c.doRequest(sendMsgMtd, params)
	if err != nil {
		return e.Wrap(reqErr, err)
	}
	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(reqErr, err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(reqErr, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(readErr, err)
	}
	return body, nil
}

func createBasePath(bp string) string {
	return fmt.Sprintf("%s%s", prefixBp, bp)
}
