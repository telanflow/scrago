# Teler   

  A micro crawler framework. achieved by GOLANG.

[![Build Status](https://travis-ci.org/telanflow/scrago.svg?branch=master)](https://travis-ci.org/telanflow/scrago) [![GitHub stars](https://img.shields.io/github/stars/telanflow/scrago.svg)](https://github.com/telanflow/scrago/stargazers) [![Go version](https://img.shields.io/badge/Go-%3E1.7-brightgreen.svg)](https://github.com/telanflow/scrago)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
[![LICENSE](https://img.shields.io/badge/license-NPL%20(The%20996%20Prohibited%20License)-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

## Quick Start

#### Download and install

    go get github.com/telanflow/scrago
    
#### Create file `my_spider.go`
```go
package main

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/telanflow/scrago"
	"github.com/telanflow/scrago/pages"
	"github.com/telanflow/scrago/downloader"
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
	scrago.New(&MySpider{}).AddUrl("https://www.baidu.com").Run()
}

```

#### Build and run

    go build my_spider.go
    ./my_spider
    
## Documentation
[中文文档](https://github.com/telanflow/scrago/wiki/%E6%A1%86%E6%9E%B6%E7%AE%80%E4%BB%8B)
    
## License

teler licensed under the Apache Licence, Version 2.0
(http://www.apache.org/licenses/LICENSE-2.0.html).
