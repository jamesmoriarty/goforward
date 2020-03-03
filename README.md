# Goforward

[![Latest Tag][2]][3] [![Go Report Card][4]][5] [![GitHub Workflow Status][6]][7]

Go forward proxy with rate limiting.

![Screenshot][1]

# Download

Releases can be downloaded from [here][3].

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
[2]: https://img.shields.io/github/v/tag/jamesmoriarty/goforward.svg?logo=github&label=latest
[3]: https://github.com/jamesmoriarty/goforward/releases
[4]: https://goreportcard.com/badge/github.com/jamesmoriarty/goforward
[5]: https://goreportcard.com/report/github.com/jamesmoriarty/goforward
[6]: https://img.shields.io/github/workflow/status/jamesmoriarty/goforward/Release
[7]: https://github.com/jamesmoriarty/goforward/actions?query=workflow%3ARelease
