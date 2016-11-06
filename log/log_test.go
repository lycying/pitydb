package log_test

import (
	"fmt"
	"github.com/lycying/pitydb/log"
	"testing"
)

func TestLog(t *testing.T) {
	logger, err0 := log.New(log.INFO, "")
	if err0 != nil {
		println(err0.Error())
		return
	}
	defer logger.Close()

	logger.Debug("will not print")
	logger.Info("test %v", " on the fly")
	logger.Warn("test %v", " on the fly")
	logger.Error("test %v", " on the fly")

	err := fmt.Errorf("it's the end of the world")
	logger.Err(err, "error stack")

}
