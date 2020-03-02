# Goforward

[![Go Report Card](https://goreportcard.com/badge/github.com/jamesmoriarty/goforward)](https://goreportcard.com/report/github.com/jamesmoriarty/goforward)

Go forward proxy with rate limiting.

![Screenshot][1]

# Install

```
go get -v github.com/jamesmoriarty/goforward
go install github.com/jamesmoriarty/goforward
```

# Usage

```
.\goforward.exe
```

```
.\goforward.exe -h
Usage of .\goforward.exe:
  -port string
        Proxy listen port (default "8888")
  -proto string
        Proxy protocol (default "http")
  -rate int
        Proxy bandwidth ratelimit (default 524288)
```

# Build 

```
go build
```

# Test

```
> go test -c
> .\goforward.test.exe
time="2020-03-02T22:48:13+11:00" level=info msg="Goforward listening on :8888 with ratelimit 512K"
time="2020-03-02T22:48:13+11:00" level=info msg="http://127.0.0.1:8080/goforward.exe" Method=GET RemoteAddr="127.0.0.1:63828"
time="2020-03-02T22:48:26+11:00" level=info msg="http://127.0.0.1:8080/goforward.exe" Duration=13.2064385s Rate=550.8K/s Size=7.1M
PASS
```

[1]: docs/screenshot.PNG
