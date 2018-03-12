package teler

import (
	"time"
	"net/http"
	"github.com/teler/pages"
	"github.com/teler/scheduler"
	"github.com/teler/downloader"
)

type HandlerFunc func(*Context)
type HandlerChain []HandlerFunc
func (c HandlerChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

type Context struct {
	Keys 		map[string]interface{}
	request		*http.Request
	response	*http.Response
	core		*Core
	index		int8
}

func NewContext(c *Core) *Context {
	return &Context{
		core: c,
		index: -1,
	}
}

func (c *Context) reset() {
	c.Keys = nil
}

func (c *Context) Copy() *Context {
	var cp = *c
	c.request = nil
	c.response = nil
	c.index = -1
	return &cp
}

func (c *Context) Next() {
	c.index++
	for s := int8(len(c.core.middleware)); c.index < s; c.index++ {
		c.core.middleware[c.index](c)
	}
}

func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.Keys[key]
	return
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (c *Context) AddUrl(u string) {
	c.AddQueue(NewRequest(http.MethodGet, u, nil))
}

func (c *Context) AddRequest(req *Request) {
	c.AddQueue(req)
}

func (c *Context) AddQueue(v scheduler.QueueElement) {
	c.core.scheduler.Push(v)
}

func (c *Context) GetDownloader() downloader.Downloader {
	return c.core.downloader
}

func (c *Context) GetRequest() *http.Request {
	return c.request
}

func (c *Context) setRequest(req *http.Request) {
	c.request = req
}

func (c *Context) GetResponse() *http.Response {
	return c.response
}

func (c *Context) setResponse(resp *http.Response) {
	c.response = resp
}

func (c *Context) HttpGet(u string) *pages.Page {
	req := NewRequest(http.MethodGet, u, nil).GetRequest()
	req.Header.Set("User-Agent", downloader.RandomUA())
	res, _ := c.core.downloader.Do(req)
	return pages.NewPageForRes(res)
}

func (c *Context) HttpPost(u string, params interface{}) *pages.Page {
	req := NewRequest(http.MethodPost, u, params).GetRequest()
	req.Header.Set("User-Agent", downloader.RandomUA())
	res, _ := c.core.downloader.Do(req)
	return pages.NewPageForRes(res)
}