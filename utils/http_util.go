package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	log "github.com/sirupsen/logrus"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

func DoHttpCall(method, url string, body []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	var req *http.Request
	var err error
	if method == "GET" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	}

	if err != nil {
		log.Error("Failed to get new http request", err)
		return nil, err
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	reqOut, _ := httputil.DumpRequest(req, true)
	log.Debugln(string(reqOut))
	resp, err := client.Do(req)
	respOut, _ := httputil.DumpResponse(resp, true)
	log.Debugln(string(respOut))

	// Print the request method, URL, headers, and body
	fmt.Printf("Request method: %s\n", req.Method)
	fmt.Printf("Request URL: %s\n", req.URL.String())
	fmt.Printf("Request headers: %v\n", req.Header)
	fmt.Printf("Request body: %s\n", string(body))
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}