package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/lGreenLightl/link-saver-bot/lib/consts"
	"github.com/lGreenLightl/link-saver-bot/lib/e"
)

type Client struct {
	host       string
	basePath   string
	httpClient http.Client
}

func NewClient(host string, token string) Client {
	return Client{
		host:       host,
		basePath:   basePath(token),
		httpClient: http.Client{},
	}
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.makeRequest(consts.SendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.makeRequest(consts.GetUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	response := UpdatesResponse{}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (c *Client) makeRequest(method string, query url.Values) (data []byte, err error) {
	currentURL := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	request, err := http.NewRequest(http.MethodGet, currentURL.String(), nil)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}
	request.URL.RawQuery = query.Encode()

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, e.Wrap("can't do request", err)
	}

	return responseBody, nil
}

func basePath(token string) string {
	return "bot" + token
}
