package crawler

import (
	"github.com/telanflow/crawler/pages"
	"github.com/telanflow/crawler/pipeline"
)

type Crawler interface {
	Init(ctx *Context)
	Process(ctx *Context, page *pages.Page)
	pipeline.Pipeline
}