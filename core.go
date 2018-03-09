package teler

import (
	"github.com/teler/middleware"
	"github.com/teler/scheduler"
	"github.com/teler/downloader"
	"github.com/teler/pages"
	"github.com/teler/pipeline"
	"github.com/teler/log"
	"net/http"
	"time"
)

type DownloadBeforeFunc func(r scheduler.QueueElement, next func())
type DownloadAfterFunc  func(page *pages.Page, next func())

type Core struct {
	crawler		Crawler					// 爬虫
	threads		uint					// 线程数
	scheduler	*scheduler.Scheduler	// 调度器
	downloader	downloader.Downloader	// 下载器
	middleware	middleware.Middleware	// 中间件
	pipelines	[]pipeline.Pipeline		// 输出管道
	sleep		time.Duration			// 任务间隔
}

func New(c Crawler) *Core {
	return &Core {
		crawler: 	c,
		threads: 	10,
		scheduler: 	scheduler.New(),
		downloader:	downloader.NewHttpDownload(),
		middleware: middleware.NewHandler(),
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

func (self *Core) Sleep(t time.Duration) {
	self.sleep = t
}

func (self *Core) Use(k string, v interface{}) *Core {
	switch k {
	case MIDDLEWARE_DownloadBefore:
		if hand, ok := v.(DownloadBeforeFunc); ok {
			self.middleware.Use(MIDDLEWARE_DownloadBefore, hand)
		} else {
			log.Panic("type func not DownloadBeforeFunc")
		}
	case MIDDLEWARE_DownloadAfter:
		if hand, ok := v.(DownloadAfterFunc); ok {
			self.middleware.Use(MIDDLEWARE_DownloadAfter, hand)
		} else {
			log.Panic("type func not DownloadAfterFunc")
		}
	}

	return self
}

func (self *Core) Run() {
	log.Info("========== Spider Start ==========")
	defer log.Info("========== Spider End ==========")

	// 初始化
	var c = NewContext(self.scheduler, self.downloader)
	self.crawler.Init(c)

	self.scheduler.Use(scheduler.WithHandler( func(task scheduler.QueueElement) {
		ctx := c.Copy()

		log.Info("Task ID: " + task.Id())

		// 下载器中间件 - Before
		self.middleware.Exec(MIDDLEWARE_DownloadBefore, func(v interface{}, next func()) {
			v.(DownloadBeforeFunc)(task, next)
		})

		req := task.GetRequest()
		resp, err := self.downloader.Do(req)
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		page := pages.NewPageForRes(resp)
		defer page.Free()

		// 下载器中间件 - After
		self.middleware.Exec(MIDDLEWARE_DownloadAfter, func(v interface{}, next func()) {
			v.(DownloadAfterFunc)(page, next)
		})

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
	}))

	self.scheduler.Dispatch()
}
