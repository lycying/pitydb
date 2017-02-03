package yard

type PageManagement struct {
	pageMap    map[uint32]*Page
	nextPageID uint32
}

func NewPageMgr() *PageManagement {
	return &PageManagement{
		pageMap:    make(map[uint32]*Page),
		nextPageID: uint32(0),
	}
}

func (mgr *PageManagement) AddPage(pg *Page) {
	key := pg.pgID.GetValue().(uint32)
	mgr.pageMap[key] = pg
}

func (mgr *PageManagement) GetPage(pageId uint32) *Page {
	v, ok := mgr.pageMap[pageId]
	if ok {
		return v
	} else {
		println("FUCK.................", pageId)
	}
	//TODO read it from disk
	return v
}

func (mgr *PageManagement) RemovePage(pageId uint32) {
	mgr.pageMap[pageId] = nil
}

func (mgr *PageManagement) NextPageID() uint32 {
	mgr.nextPageID++
	return mgr.nextPageID
}
