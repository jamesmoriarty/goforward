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

[1]: docs/screenshot.PNG
