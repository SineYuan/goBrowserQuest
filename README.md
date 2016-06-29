goBrowserQuest server documentation
================================

go implementation for [BrowserQuest](https://github.com/mozilla/BrowserQuest) server

Installation
-------------

```
go get github.com/SineYuan/goBrowserQuest
```

Configuration
-------------

```
  -config string
        configuration file path (default "./config.json")
  -client string
        BrowserQuest root directory to serve if provided
  -prefix string
        request url prefix when client is provided, cannot be '/'  (default "/game")

```

Deployment
----------

```
cd $GOPATH/src/github.com/SineYuan/goBrowserQuest
go build main.go
./main -config /path/to/config.json -client /path/to/BrowserQuest 
```

then you can play game at `http://{HOST}:{PORT}/game/client/index.html`

TODO
----------
goBrowserQuest have yet to implement all the function of BrowserQuest server. welcome to forks