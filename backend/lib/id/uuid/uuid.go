package uuid

import (
	"github.com/google/uuid"
)

func GetNext() string {
	id,_:=uuid.NewUUID()
	return id.String()
}
