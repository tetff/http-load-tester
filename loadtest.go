package loadtest

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Config struct {
	Random              bool
	IntervalMinMillisec int
	IntervalMaxMillisec int
}

type Load struct {
	client          *http.Client
	request         *http.Request
	responseHandler func(*http.Response)

	c *Config
}

func New(client *http.Client, responseHandler func(*http.Response)) *Load {
	return &Load{
		client:          client,
		responseHandler: responseHandler,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Random:              false,
		IntervalMinMillisec: 4000,
		IntervalMaxMillisec: 5000,
	}
}

func DefaultRandomConfig() *Config {
	return &Config{
		Random:              true,
		IntervalMinMillisec: 4000,
		IntervalMaxMillisec: 5000,
	}
}

func (l *Load) Setup(req *http.Request, c *Config) {
	l.request = req
	l.c = c
}

func (l *Load) Send(stop chan bool) {
	var ticker *time.Ticker
	if l.c.Random {
		if l.c.IntervalMinMillisec > l.c.IntervalMaxMillisec {
			log.Fatal(errors.New("Minimum interval bigger than maximum"))
		}
		ticker = time.NewTicker(time.Duration(l.Random()) * time.Millisecond)
	} else {
		ticker = time.NewTicker(time.Duration(l.c.IntervalMinMillisec) * time.Millisecond)
	}
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			resp, err := l.client.Do(l.request)
			if err != nil {
				log.Fatal(err)
			}
			l.responseHandler(resp)
		}
	}

}

func (l *Load) Random() int {
	rand.Seed(time.Now().UnixNano())
	return l.c.IntervalMinMillisec + rand.Intn(l.c.IntervalMaxMillisec-l.c.IntervalMinMillisec)
}
