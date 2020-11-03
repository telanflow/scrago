package pipeline

import (
	"github.com/telanflow/scrago/pages"
)

type Pipeline interface {
	Output(items *pages.PageItem)
}