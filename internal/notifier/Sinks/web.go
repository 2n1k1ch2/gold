package Sinks

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

const SendTimeOut = 5

type WebHook struct {
	client *http.Client
	url    string
}

func NewWebHook(url string) *WebHook {
	return &WebHook{
		url: url,
		client: &http.Client{
			Timeout: SendTimeOut * time.Second,
		},
	}
}

func (w *WebHook) Send(payload []byte) error {
	req, err := http.NewRequest(http.MethodPost, w.url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("webhook: non-2xx response")
	}
	return nil

}

func (w *WebHook) Close() {
	w.client.CloseIdleConnections()
}
