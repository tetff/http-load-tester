package loadtest

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"time"
)

type load struct {
	client          *http.Client
	request         *http.Request
	responseHandler func(*http.Response)
}

type Request struct {
	Method  string
	Address string
	Request []byte
	Headers map[string]string
	Ctx     *context.Context
}

func New(client *http.Client, responseHandler func(*http.Response)) load {
	return load{
		client:          client,
		responseHandler: responseHandler,
	}
}

func (l *load) Setup(req Request) error {
	request, err := http.NewRequest(req.Method, req.Address, bytes.NewBuffer(req.Request))
	if err != nil {
		return err
	}

	for k, v := range req.Headers {
		request.Header.Set(k, v)
	}

	request.WithContext(*req.Ctx)

	l.request = request
	return nil
}

func (l *load) Send(intervalInSeconds uint) {
	ticker := time.NewTicker(time.Duration(intervalInSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			resp, err := l.client.Do(l.request)
			if err != nil {
				log.Fatal(err)
			}
			l.responseHandler(resp)
		}
	}
}
