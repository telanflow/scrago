package pipeline

import (
	"github.com/teler/pages"
)

type Pipeline interface {
	Output(items *pages.PageItem)
}