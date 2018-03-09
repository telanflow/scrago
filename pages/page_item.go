package pages

import "sync"

type PageItem struct {
	items map[string]string
}

var (
	itemPool = sync.Pool{
		New: func() interface{} {
			return &PageItem{
				items: make(map[string]string),
			}
		},
	}
)

func NewItem() *PageItem {
	return itemPool.Get().(*PageItem)
}

func (i *PageItem) Set(k string, v string) {
	i.items[k] = v
}

func (i *PageItem) Get(k string) string {
	if v, ok := i.items[k]; ok {
		return v
	}

	return ""
}

func (i *PageItem) GetAll() map[string]string {
	return i.items
}

func (i *PageItem) Free() {
	for k := range i.items {
		delete(i.items, k)
	}
	itemPool.Put(i)
}