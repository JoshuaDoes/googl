package googl

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type Googl struct {
	Key string
}

type ShortMsg struct {
	Kind    string `json:"kind"`
	Id      string `json:"id"`
	LongUrl string `json:"longUrl"`
}

type LongMsg struct {
	Kind    string `json:"kind"`
	Id      string `json:"id"`
	LongUrl string `json:"longUrl"`
	Status  string `json:"status"`
}

func NewClient(key string) *Googl {
	return &Googl{Key: key}
}

func (c *Googl) Shorten(url string) (*ShortMsg, error) {
	request := gorequest.New()
	gUrl := "https://www.googleapis.com/urlshortener/v1/url?key=" + c.Key

	if c.Key == "" {
		return nil, errors.New("You must set the Google URL Shortener API key")
	} else if url == "" {
		return nil, errors.New("You must set the URL to be shortened")
	} else {
		resp, _, _ := request.Post(gUrl).
			Set("Accept", "application/json").
			Set("Content-Type", "application/json").
			Send(`{"longUrl":"` + url + `"}`).End()
		if resp.Status == "200 OK" {
			shortMsg := &ShortMsg{}
			err := unmarshal(resp, shortMsg)
			return shortMsg, err
		} else {
			return nil, errors.New("Unknown error")
		}
	}
}

func (c *Googl) Expand(shortUrl string) (*LongMsg, error) {
	request := gorequest.New()
	gUrl := "https://www.googleapis.com/urlshortener/v1/url?key=" + c.Key + "&shortUrl=" + shortUrl

	if c.Key == "" {
		return nil, errors.New("You must set the Google URL Shortener API key")
	} else if shortUrl == "" {
		return nil, errors.New("You must set the URL to be shortened")
	} else {
		resp, _, _ := request.Get(gUrl).
			Set("Accept", "application/json").
			Set("Content-Type", "application/json").End()
		if resp.Status == "200 OK" {
			longMsg := &LongMsg{}
			err := unmarshal(resp, longMsg)
			return longMsg, err
		} else {
			return nil, errors.New("Unknown error")
		}
	}
}

func unmarshal(body *http.Response, target interface{}) error {
	defer body.Body.Close()
	return json.NewDecoder(body.Body).Decode(target)
}
