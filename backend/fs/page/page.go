package page

import "../row"

const DEFAULT_PAGE_SIZE = 1024 * 16

type Page struct {
	offset   int //64TiB max
	previous int
	next     int
	typ      int
	rows     []*row.Row
}

func NewPage() *Page {
	return &Page{
	}
}

//After this,the buf became a bean
func (this *Page) ToBean(buf []byte) {
}
//After this,the bean became an byte array
func (this *Page) ToBytes() {
}
