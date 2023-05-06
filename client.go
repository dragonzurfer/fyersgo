package fyers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dragonzurfer/fyersgo/utils"

	log "github.com/sirupsen/logrus"
)

type client struct {
	host        string
	apiKey      string
	accessToken string
	debug       bool
	httpClient  *http.Client
}

func New(apiKey, accessToken string) *client {
	return &client{
		apiKey:      apiKey,
		accessToken: accessToken,
		debug:       false,
		httpClient: &http.Client{
			Timeout: time.Duration(20) * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Second * time.Duration(20),
			},
		},
	}
}

func (c *client) WithHttpClient(httpClient *http.Client) *client {
	c.httpClient = httpClient
	return c
}

func (c *client) EnableDebug() *client {
	c.debug = true
	return c
}

func (c *client) WithHost(host string) *client {
	c.host = host
	return c
}

func (c *client) invoke(method, url string, body interface{}) ([]byte, error) {
	headerMap := map[string]string{
		"Authorization": fmt.Sprintf("%s", c.apiKey+":"+c.accessToken),
		"Content-Type":  "application/json",
	}
	fmt.Printf("Headers: %v\n", headerMap)
	var bodyByte []byte
	if bodyByteArr, err := json.Marshal(body); err != nil {
		return nil, err
	} else {
		bodyByte = bodyByteArr
	}
	if resp, err := utils.DoHttpCall(method, url, bodyByte, headerMap); err != nil {
		log.Error("Failed to make http call", err)
		return nil, err
	} else {
		fmt.Println(string(resp))
		fmt.Println(url)

		return resp, nil
	}
}

func (c *client) toUri(version, path string, params ...string) string {
	host := Host
	if len(c.host) > 0 {
		host = c.host
	}
	return fmt.Sprintf("%s%s%s%s", host, version, path, strings.Join(params, ""))
}

func (c *client) toUriData(version, path string, params ...string) string {
	host := DataHost
	if len(c.host) > 0 {
		host = c.host
	}
	return fmt.Sprintf("%s%s%s%s", host, version, path, strings.Join(params, ""))
}