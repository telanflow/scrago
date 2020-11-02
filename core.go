package crawler

import (
	"github.com/telanflow/crawler/downloader"
	"github.com/telanflow/crawler/log"
	"github.com/telanflow/crawler/pages"
	"github.com/telanflow/crawler/pipeline"
	"github.com/telanflow/crawler/scheduler"
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

func (c *Core) AddUrl(u string) *Core {
	c.scheduler.Push(NewRequest(http.MethodGet, u, nil))
	return c
}

func (c *Core) AddRequest(req *Request) *Core {
	c.scheduler.Push(req)
	return c
}

func (c *Core) AddPipeline(p ...pipeline.Pipeline) *Core {
	c.pipelines = append(c.pipelines, p...)
	return c
}

func (c *Core) SetThreads(n int) *Core {
	c.threads = n
	return c
}

func (c *Core) Use(middleware ...HandlerFunc) *Core {
	c.middleware = append(c.middleware, middleware...)
	return c
}

func (c *Core) Sleep(t time.Duration) *Core {
	c.sleep = t
	return c
}

func (c *Core) Run() {
	log.Info("========== Spider Start ==========")
	defer log.Info("========== Spider End ==========")

	var gCtx = NewContext(c)

	c.crawler.Init(gCtx)

	c.Use(func(ctx *Context) {
		resp, err := c.downloader.Do(ctx.GetRequest())
		if err != nil {
			log.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		ctx.setResponse(resp)
	})

	c.scheduler.SetHandler(func(task scheduler.QueueElement) {
		log.Info("Task ID: " + task.Id())

		ctx := gCtx.Copy()

		ctx.setRequest(task.GetRequest())

		// Middleware
		ctx.Next()

		page := pages.NewPageForRes(ctx.GetResponse())
		defer page.Free()

		// Page Process
		c.crawler.Process(ctx, page)

		// Page Item
		items := page.GetItem()
		c.crawler.Output(items)

		// Pipeline
		if len(c.pipelines) > 0 {
			for _, v := range c.pipelines {
				v.Output(items)
			}
		}

		if c.sleep > 0 {
			time.Sleep(c.sleep)
		}
	})

	c.scheduler.Dispatch(c.threads)
}
