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
go test
```

```
time="2020-03-03T22:51:58+11:00" level=info msg="Goforward listening on :8888 with ratelimit 512K"
time="2020-03-03T22:51:59+11:00" level=info msg="http://127.0.0.1:8080/goforward.exe" Method=GET RemoteAddr="127.0.0.1:63286"
time="2020-03-03T22:52:12+11:00" level=info msg="http://127.0.0.1:8080/goforward.exe" Duration=13.2338348s Rate=550.7K/s Size=7.1M
PASS
ok      github.com/jamesmoriarty/goforward      14.042s
```

[1]: docs/screenshot.PNG
