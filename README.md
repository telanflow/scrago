# Teler   [![Build Status](https://travis-ci.org/ziiber/teler.svg?branch=master)](https://travis-ci.org/ziiber/teler) [![GitHub stars](https://img.shields.io/github/stars/ziiber/teler.svg)](https://github.com/ziiber/teler/stargazers) [![Go version](https://img.shields.io/badge/Go-%3E1.7-brightgreen.svg)](https://github.com/ziiber/teler)



  A micro crawler framework. achieved by GOLANG.

## Quick Start

#### Download and install

    go get github.com/ziiber/teler
    
#### Create file `my_spider.go`
```go
package main

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/ziiber/teler"
	"github.com/ziiber/teler/pages"
	"github.com/ziiber/teler/downloader"
)


type MySpider struct{
	jar http.CookieJar
}

// Init
func (m *MySpider) Init(ctx *teler.Context) {
	// Set the persistent cookie.
	m.jar, _ = cookiejar.New(nil)
	ctx.GetDownloader().UseOptions(downloader.WithCookieJar(m.jar))

	// Add Target Url
	//ctx.AddUrl("https://www.baidu.com")
}

// Page Process
func (m *MySpider) Process(ctx *teler.Context, page *pages.Page) {

}

// Pipeline Output
func (m *MySpider) Output(items *pages.PageItem) {

}


func main() {
	// Start Spider
	teler.New(&MySpider{}).AddUrl("https://www.baidu.com").Run()
}

```

#### Build and run

    go build my_spider.go
    ./my_spider
    
## Documentation
[中文文档](https://github.com/ziiber/teler/wiki/%E6%A1%86%E6%9E%B6%E7%AE%80%E4%BB%8B)
    
## License

teler licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
