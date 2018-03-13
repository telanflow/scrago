package pipeline

import (
	"github.com/ziiber/teler/pages"
)

type Pipeline interface {
	Output(items *pages.PageItem)
}