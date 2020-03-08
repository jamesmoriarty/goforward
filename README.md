# Goforward

[![Latest Tag][2]][3] [![Go Report Card][4]][5] [![GitHub Workflow Status][6]][7]

Go forward proxy with rate limiting. The code is based on [Michał Łowicki's][8] 100 LOC forward proxy.

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
goforward
```

```
goforward -h
Usage of goforward:
  -port string
        Proxy listen port (default "8888")g
  -rate int
        Proxy bandwidth ratelimit (default 524288)
```

use with `.exe` on windows.

# Build 

```
go build
```

# Test

```
go test
```

# Why

I needed a way to download 53GB without making my household internet unusable. In summary:

1. [Free games](https://www.pcgamer.com/au/faeria-is-the-next-free-epic-game-store-game-kingdom-come-deliverance-and-aztez-are-available-now/).
2. [Australia's terrible internet](https://en.wikipedia.org/wiki/List_of_countries_by_Internet_connection_speeds).
3. [Learning Go](https://golang.org/).

## First Solution

Shape the traffic in the application.

[![Application Bandwidth Shaping][9]][9]

## Second Solution

Shape the traffic in kernal space.

[![Windows Filtering Platform][10]][10]

## Third Solution

Shape the traffic in user space.

[![Forward Proxy][11]][11]

[1]: docs/screenshot.PNG
[2]: https://img.shields.io/github/v/tag/jamesmoriarty/goforward.svg?logo=github&label=latest
[3]: https://github.com/jamesmoriarty/goforward/releases
[4]: https://goreportcard.com/badge/github.com/jamesmoriarty/goforward
[5]: https://goreportcard.com/report/github.com/jamesmoriarty/goforward
[6]: https://img.shields.io/github/workflow/status/jamesmoriarty/goforward/Release
[7]: https://github.com/jamesmoriarty/goforward/actions?query=workflow%3ARelease
[8]: https://github.com/mlowicki
[9]: docs/diagram-1.png
[10]: docs/diagram-2.png
[11]: docs/diagram-3.png