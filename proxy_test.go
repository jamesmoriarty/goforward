package main

import (
	"net/http"
	"testing"
	"net/url"
	"io/ioutil"
	"time"
)

func withStubHTTPServer(port string, directory string, f func()) {
	http.Handle("/", http.FileServer(http.Dir(directory)))
	
	go http.ListenAndServe(":"+port, nil)

	f()
}

type benchmark struct {
	Rate int
	DurationMin float64
	DurationMax float64
}

func TestBenchmarks(t *testing.T) {
	benchmarks := []benchmark {
		benchmark {
			Rate: 512*1024,
			DurationMin: 12,
			DurationMax: 14,
		},
	}

	withStubHTTPServer("8080", ".", func() {
		for _, b := range benchmarks {
		
			go proxy("http", "8888", b.Rate)

			proxyURL, err := url.Parse("http://127.0.0.1:8888")
			client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			response, err := client.Get("http://127.0.0.1:8080/goforward.exe")

			defer response.Body.Close()

			start := time.Now()

			_, err = ioutil.ReadAll(response.Body)

			duration := time.Since(start)

			if err != nil {
				t.Errorf(err.Error())
			}

			if duration.Seconds() > b.DurationMax {
				t.Errorf("Exceded duration.")
			}

			if duration.Seconds() < b.DurationMin {
				t.Errorf("Finished early.")
			}
		}
	})
}