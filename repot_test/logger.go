package repot_test

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	//"github.com/spf13/afero"
)

func InitLogger(t *testing.T) {
	file, err := os.OpenFile(fmt.Sprintf("/tmp/repot.%s.log", "test"), os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.SetLevel(log.DebugLevel)
}
