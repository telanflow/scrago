package pages

import (
	"net/http"
	"sync"
	"io"
)

type Page struct {
	Request	*http.Request
	Header	http.Header
	Cookies	[]*http.Cookie
	Body	[]byte
	item	*PageItem
}

var (
	pagePool = sync.Pool{
		New: func() interface{} {
			return &Page{
				Cookies: make([]*http.Cookie, 0),
				Body: make([]byte, 0),
				item: NewItem(),
			}
		},
	}
)

func NewPage() *Page {
	return pagePool.Get().(*Page)
}

func NewPageForRes(res *http.Response) *Page {
	page := pagePool.Get().(*Page)
	page.Request = res.Request
	page.Header = res.Header
	io.Copy(page, res.Body)

	cookies := res.Cookies()
	page.Cookies = make([]*http.Cookie, len(cookies))
	copy(page.Cookies, cookies)

	return page
}

func (self *Page) AddField(k, v string) {
	self.item.Set(k, v)
}

func (self *Page) Field(k string) string {
	return self.item.Get(k)
}

func (self *Page) GetItem() *PageItem {
	return self.item
}

func (self *Page) Write(p []byte) (n int, err error) {
	self.Body = make([]byte, len(p))
	n = copy(self.Body, p)
	return
}

func (self *Page) Read(p []byte) (n int, err error) {
	p = make([]byte, len(self.Body))
	n = copy(p, self.Body)
	return
}

func (self *Page) Free() {
	self.Request = nil
	self.Header = make(http.Header)
	self.Cookies = make([]*http.Cookie, 0)
	self.Body = make([]byte, 0)
	pagePool.Put(self)
}