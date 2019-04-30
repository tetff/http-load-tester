package loadtest

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestLoadWithoutRandom(t *testing.T) {
	loader := New(&http.Client{
		Timeout: 3000 * time.Millisecond,
	}, responseHandler)

	req, _ := http.NewRequest("GET", "https://google.com", nil)
	loader.Setup(req, &Config{
		Random:              false,
		IntervalMinMillisec: 1000,
		IntervalMaxMillisec: 2000,
	})

	var stop = make(chan bool)
	go loader.Send(stop)

	time.Sleep(10 * time.Second)
	close(stop)
}

func TestLoadWithRandom(t *testing.T) {
	loader := New(&http.Client{
		Timeout: 3000 * time.Millisecond,
	}, responseHandler)

	req, _ := http.NewRequest("GET", "https://google.com", nil)
	loader.Setup(req, &Config{
		Random:              true,
		IntervalMinMillisec: 500,
		IntervalMaxMillisec: 2000,
	})

	var stop = make(chan bool)
	go loader.Send(stop)

	time.Sleep(10 * time.Second)
	close(stop)
}

func responseHandler(resp *http.Response) {
	fmt.Println(resp.StatusCode)
}
