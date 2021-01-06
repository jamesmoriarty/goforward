package goforward

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"testing"
	"time"
)

func bytes(path string) int {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	return (int)(info.Size())
}

func with(port string, directory string, f func()) {
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
			Rate: 512 * 1024,
		},
		{
			Rate: 256 * 1024,
		},
	}

	with("8080", ".", func() {
		shutdown := make(chan os.Signal, 1)

		for _, b := range benchmarks {
			go Listen("8888", b.Rate, shutdown)

			fileName := "goforward.exe"
			proxyURL, _ := url.Parse("http://127.0.0.1:8888")
			client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
			response, err := client.Get("http://127.0.0.1:8080/" + fileName)

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

			DurationExpected := (float64)(bytes("./"+fileName) / b.Rate)

			fmt.Printf("for %v@%v took %v expected %v\n", bytes("./"+fileName)/1024, b.Rate/1024, duration.Seconds(), DurationExpected)

			if duration.Seconds() > (DurationExpected * 1.2) {
				t.Errorf("for %v@%v took %v expected <%v", bytes("./"+fileName)/1024, b.Rate/1024, duration.Seconds(), DurationExpected)
			}

			if duration.Seconds() < (DurationExpected * 0.8) {
				t.Errorf("for %v@%v took %v expected >%v", bytes("./"+fileName)/1024, b.Rate/1024, duration.Seconds(), DurationExpected)
			}

			shutdown <- syscall.SIGKILL

			time.Sleep(3)
		}
	})
}
