package zaplog

import (
	"context"
	"testing"
)

type tContextHook struct {
	call int
}

func (p *tContextHook) WithContext(ctx context.Context) (args []interface{}) {
	p.call++
	return nil
}

func TestSetContextHook(t *testing.T) {
	var logger Logger
	var hook tContextHook
	SetContextHook(&hook)(&logger)
	if got, want := logger.hook, &hook; got != want {
		t.Fatalf("hook: got %v, want %v", got, want)
	}
}
