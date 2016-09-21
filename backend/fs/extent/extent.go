package extent

import "../page"

const DEFAULT_EXTENT_SIZE = 64

type Extent struct {
	pages []*page.Page
}
