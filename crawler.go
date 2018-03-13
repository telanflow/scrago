package teler

import (
	"github.com/ziiber/teler/pages"
	"github.com/ziiber/teler/pipeline"
)

type Crawler interface {
	Init(ctx *Context)
	Process(ctx *Context, page *pages.Page)
	pipeline.Pipeline
}