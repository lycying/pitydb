package mut

import "github.com/lycying/pitydb/log"

var logger *log.Logger

func init() {
	logger, _ = log.New(log.DEBUG, "")
}
