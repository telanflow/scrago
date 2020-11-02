package pipeline

import (
	"github.com/telanflow/crawler/pages"
)

type Pipeline interface {
	Output(items *pages.PageItem)
}