package teler

import (
	"github.com/teler/scheduler"
	"github.com/teler/downloader"
	"github.com/teler/pages"
	"github.com/teler/pipeline"
	"github.com/teler/log"
	"net/http"
	"time"
)

type Core struct {
	crawler		Crawler					// 爬虫
	threads		int						// 线程数
	scheduler	*scheduler.Scheduler	// 调度器
	downloader	downloader.Downloader	// 下载器
	middleware	HandlerChain			// 中间件
	pipelines	[]pipeline.Pipeline		// 输出管道
	sleep		time.Duration			// 任务间隔
}

func New(c Crawler) *Core {
	return &Core {
		crawler: 	c,
		threads: 	20,
		scheduler: 	scheduler.New(),
		downloader:	downloader.NewHttpDownload(),
		middleware: make(HandlerChain, 0),
		pipelines:  make([]pipeline.Pipeline, 0),
		sleep: 		0,
	}
}

func (self *Core) AddUrl(u string) *Core {
	self.scheduler.Push(NewRequest(http.MethodGet, u, nil))
	return self
}

func (self *Core) AddRequest(req *Request) *Core {
	self.scheduler.Push(req)
	return self
}

func (self *Core) AddPipeline(p ...pipeline.Pipeline) *Core {
	self.pipelines = append(self.pipelines, p...)
	return self
}

func (self *Core) SetThreads(n int) *Core {
	self.threads = n
	return self
}

func (self *Core) Use(middleware ...HandlerFunc) *Core {
	self.middleware = append(self.middleware, middleware...)
	return self
}

func (self *Core) Sleep(t time.Duration) *Core {
	self.sleep = t
	return self
}

func (self *Core) Run() {
	log.Info("========== Spider Start ==========")
	defer log.Info("========== Spider End ==========")

	var gCtx = NewContext(self)

	self.crawler.Init(gCtx)

	self.Use(func(ctx *Context) {
		resp, err := self.downloader.Do(ctx.GetRequest())
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		ctx.setResponse(resp)
	})

	self.scheduler.SetHandler(func(task scheduler.QueueElement) {
		log.Info("Task ID: " + task.Id())

		ctx := gCtx.Copy()

		ctx.setRequest(task.GetRequest())

		// Middleware
		ctx.Next()

		page := pages.NewPageForRes(ctx.GetResponse())
		defer page.Free()

		// Page Process
		self.crawler.Process(ctx, page)

		// Page Item
		items := page.GetItem()
		self.crawler.Output(items)

		// Pipeline
		if len(self.pipelines) > 0 {
			for _, v := range self.pipelines {
				v.Output(items)
			}
		}

		if self.sleep > 0 {
			time.Sleep(self.sleep)
		}
	})

	self.scheduler.Dispatch(self.threads)
}
