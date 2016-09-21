package tspace

import "os"

type FileSpace struct {
	files   []*os.File //the real files
	path    string     //the space
	max     int        //in byte
	extend  bool       //if can be extend
	spaceId int        //can be ref
}
type TableSpace struct {
	spaces []*FileSpace
}
