package rollfile

import (
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	tests := []struct {
		opt Option
		chk func(f *File) bool
	}{
		{
			opt: SetLayout(DayLayout),
			chk: func(f *File) bool { return f.layout == DayLayout },
		},
		{
			opt: SetPeriod(time.Hour),
			chk: func(f *File) bool { return f.period == time.Hour },
		},
		{
			opt: SetPeriod(tickInterval / 10),
			chk: func(f *File) bool { return f.period == tickInterval },
		},
		{
			opt: SetMaxSeq(1),
			chk: func(f *File) bool { return f.maxSeq == 1 },
		},
		{
			opt: SetMaxSize(2),
			chk: func(f *File) bool { return f.maxSize == 2 },
		},
		{
			opt: DisableCreateLog(),
			chk: func(f *File) bool { return f.disableCreateLog == true },
		},
	}
	for i, tt := range tests {
		f := &File{}
		tt.opt(f)
		if !tt.chk(f) {
			t.Errorf("%d: failed to check", i)
		}
	}
}
