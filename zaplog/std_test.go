package zaplog

import (
	"context"
	"testing"
)

func TestStdContextHook(t *testing.T) {
	call := 0

	logger := StdLogger()
	logger.WithContext(context.Background())
	StdContextHook = func(ctx context.Context) []interface{} {
		call++
		return nil
	}
	logger.WithContext(context.Background())
	logger.WithContext(context.Background())

	if call != 2 {
		t.Errorf("call: got %v, want %v", call, 2)
	}
}
