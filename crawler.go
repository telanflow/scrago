package teler

import (
	"github.com/teler/pages"
	"github.com/teler/pipeline"
)

type Crawler interface {
	Init(ctx *Context)
	Process(ctx *Context, page *pages.Page)
	pipeline.Pipeline
}