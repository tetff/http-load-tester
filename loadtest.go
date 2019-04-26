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

type request struct {
	method  string
	address string
	request []byte
	headers *map[string]string
	ctx     context.Context
}

func New(client *http.Client, responseHandler func(*http.Response)) load {
	return load{
		client:          client,
		responseHandler: responseHandler,
	}
}

func (l *load) Setup(req request) error {
	request, err := http.NewRequest(req.method, req.address, bytes.NewBuffer(req.request))
	if err != nil {
		return err
	}

	for k, v := range *req.headers {
		request.Header.Set(k, v)
	}

	request.WithContext(req.ctx)

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

func NewRequest(method, address string, req []byte, headers *map[string]string, ctx context.Context) request {
	return request{
		method:  method,
		address: address,
		request: req,
		headers: headers,
		ctx:     ctx,
	}
}
