package id

import "./autoinc"

func NewAutoIncrementId32(init uint32) *autoinc.AutoIncrementId32 {
	return autoinc.NewAutoIncrementId32(init)
}
func NewAutoIncrementId64(init uint64) *autoinc.AutoIncrementId64 {
	return autoinc.NewAutoIncrementId64(init)
}
