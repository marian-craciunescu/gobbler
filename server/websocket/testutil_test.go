package websocket

import (
	"github.com/golang/mock/gomock"
	"github.com/smancke/guble/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var ctrl *gomock.Controller

func init() {
	// disable error output while testing
	// because also negative tests are tested
	protocol.LogLevel = protocol.LEVEL_ERR
}

func initCtrl(t *testing.T) func() {
	ctrl = gomock.NewController(t)
	return func() { ctrl.Finish() }
}

func enableDebugForMethod() func() {
	reset := protocol.LogLevel
	protocol.LogLevel = protocol.LEVEL_DEBUG
	return func() { protocol.LogLevel = reset }
}

func expectDone(a *assert.Assertions, doneChannel chan bool) {
	select {
	case <-doneChannel:
		return
	case <-time.After(time.Second):
		a.Fail("timeout in expectDone")
	}
}
