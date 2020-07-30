package goforward

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func withStubHTTPServer(port string, directory string, f func()) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.FileServer(http.Dir(directory)),
	}

	go server.ListenAndServe()

	f()

	server.Close()
}

type benchmark struct {
	Rate        int
	DurationMin float64
	DurationMax float64
}

func TestBenchmarks(t *testing.T) {
	benchmarks := []benchmark{
		{
			Rate:        2048 * 1024,
			DurationMin: 2,
			DurationMax: 4,
		},
		{
			Rate:        1024 * 1024,
			DurationMin: 6,
			DurationMax: 7,
		},
		{
			Rate:        512 * 1024,
			DurationMin: 12,
			DurationMax: 14,
		},

	}

	withStubHTTPServer("8080", ".", func() {
		shutdown := make(chan bool, 1)

		for _, b := range benchmarks {
			go Listen("8888", b.Rate, shutdown)

			proxyURL, _ := url.Parse("http://127.0.0.1:8888")
			client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			response, err := client.Get("http://127.0.0.1:8080/goforward.exe")

			if err != nil {
				t.Errorf(err.Error())
			}

			defer response.Body.Close()

			start := time.Now()

			_, err = ioutil.ReadAll(response.Body)

			duration := time.Since(start)

			if err != nil {
				t.Errorf(err.Error())
			}

			if duration.Seconds() > b.DurationMax {
				t.Errorf("Too slow.")
			}

			if duration.Seconds() < b.DurationMin {
				t.Errorf("Too fast.")
			}

			shutdown <- true

			time.Sleep(3)
		}
	})
}
