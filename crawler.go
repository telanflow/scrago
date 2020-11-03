package scrago

import (
	"github.com/telanflow/scrago/pages"
	"github.com/telanflow/scrago/pipeline"
)

type Crawler interface {
	Init(ctx *Context)
	Process(ctx *Context, page *pages.Page)
	pipeline.Pipeline
}