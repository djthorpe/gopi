package log_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/log"
)

type app struct {
	gopi.Unit
	*log.Log
}

func Test_Log_000(t *testing.T) {
	t.Log("Test_Log_000")
}
