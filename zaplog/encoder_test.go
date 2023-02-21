package zaplog

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	testEntry = zapcore.Entry{
		Level:   zapcore.InfoLevel,
		Time:    time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC),
		Message: "hello, world",
		Caller: zapcore.EntryCaller{
			Defined: true,
			File:    "github.com/ironzhang/tlog/zaplog/encoder_test.go",
			Line:    16,
		},
	}
	testFields = []zapcore.Field{zap.String("k1", "v1")}
)

func testEncoder(t testing.TB, name string, cfg EncoderConfig, expect string) error {
	enc, err := newEncoder(name, cfg)
	if err != nil {
		return fmt.Errorf("new %s encoder: %w", name, err)
	}
	buf, err := enc.EncodeEntry(testEntry, testFields)
	if err != nil {
		return fmt.Errorf("encode entry: %w", err)
	}
	if got, want := buf.String(), expect; got != want {
		return fmt.Errorf("got %q, want %q", got, want)
	}
	t.Logf("buffer: %v, %s", cfg, buf.String())
	return nil
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		cfg     EncoderConfig
		json    string
		console string
	}{
		{
			cfg:     EncoderConfig{},
			json:    `{"M":"hello, world","k1":"v1"}` + "\n",
			console: "hello, world\t" + `{"k1": "v1"}` + "\n",
		},
		{
			cfg: NewConsoleEncoderConfig(),
			json: `{"L":"INFO","T":"1970-01-01T00:00:00.000Z","C":"zaplog/encoder_test.go:16",` +
				`"M":"hello, world","k1":"v1"}` + "\n",
			console: "1970-01-01T00:00:00.000Z\tINFO\tzaplog/encoder_test.go:16\thello, world\t" +
				`{"k1": "v1"}` + "\n",
		},
		{
			cfg: NewJSONEncoderConfig(),
			json: `{"level":"info","ts":0,"caller":"zaplog/encoder_test.go:16",` +
				`"msg":"hello, world","k1":"v1"}` + "\n",
			console: "0\tinfo\tzaplog/encoder_test.go:16\thello, world\t" +
				`{"k1": "v1"}` + "\n",
		},
	}
	for i, tt := range tests {
		if err := testEncoder(t, "json", tt.cfg, tt.json); err != nil {
			t.Errorf("%d: test json encoder: %v", i, err)
			continue
		}
		if err := testEncoder(t, "console", tt.cfg, tt.console); err != nil {
			t.Errorf("%d: test console encoder: %v", i, err)
			continue
		}
	}
}
