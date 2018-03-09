package downloader

import (
	"time"
	"sync"
	"errors"
	"net/url"
	"strconv"
	"net/http"
)

type redirectHandler func(req *http.Request, via []*http.Request) error

type HttpDownloader struct {
	jar 				http.CookieJar		// Cookie管理器
	transport 			*http.Transport		// 单次Http请求设置
	redirect			redirectHandler		// 自定义重定向handler
	redirectTimes		int					// 重定向次数  默认10次
	timeout 			time.Duration		// 超时限制，默认0无限制    该超时限制包括连接时间、重定向和读取回复主体的时间。
	connectTimeout		time.Duration		// 连接超时限制。
	readWriteTimeout	time.Duration		// 读取超时限制。
}

var (
	clientPool = sync.Pool{
		New: func() interface{} {
			return &http.Client{}
		},
	}
)

func NewHttpDownload() Downloader {
	return &HttpDownloader{
		timeout: 0,
		redirectTimes: 10,
		transport: &http.Transport{},
	}
}

func (self *HttpDownloader) Do(req *http.Request) (*http.Response, error) {

	client := clientPool.Get().(*http.Client)
	defer clientPool.Put(client)

	// 超时限制
	client.Timeout = self.timeout

	// Cookie管理器
	client.Jar = self.jar

	// 重定向
	client.CheckRedirect = self.checkRedirect()

	if self.transport != nil {

		//self.transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		//	conn, err := net.DialTimeout(network, addr, self.connectTimeout)
		//	if err != nil {
		//		return nil, err
		//	}
		//	err = conn.SetDeadline(time.Now().Add(self.readWriteTimeout))
		//	return conn, err
		//}

		client.Transport = self.transport
	}

	return client.Do(req)
}

func (self *HttpDownloader) UseOptions(options ...Option) {
	for _, v := range options {
		v(self)
	}
}

func (self *HttpDownloader) checkRedirect() redirectHandler {
	fn := func(req *http.Request, via []*http.Request) error {
		if len(via) >= self.redirectTimes {
			return errors.New("stopped after " + strconv.Itoa(self.redirectTimes) + " redirects")
		}
		return nil
	}

	if self.redirect != nil {
		fn = self.redirect
	}

	return fn
}


func WithProxyStr(u string) Option {
	return func(download Downloader) {
		proxy, err := url.Parse(u)
		if err != nil {
			panic(err.Error())
		}

		download.(*HttpDownloader).transport.Proxy = http.ProxyURL(proxy)
	}
}

func WithProxy(u *url.URL) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).transport.Proxy = http.ProxyURL(u)
	}
}

func WithTimeout(t time.Duration) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).timeout = t
	}
}

func WithConnTimeout(t time.Duration) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).connectTimeout = t
	}
}

func WithReadWriteTimeout(t time.Duration) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).readWriteTimeout = t
	}
}

func WithCookieJar(jar http.CookieJar) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).jar = jar
	}
}

func WithRedirectTimes(n int) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).redirectTimes = n
	}
}

func WithCheckRedirect(fn func(req *http.Request, via []*http.Request) error) Option {
	return func(download Downloader) {
		download.(*HttpDownloader).redirect = fn
	}
}