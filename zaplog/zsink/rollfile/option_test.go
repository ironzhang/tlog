package rollfile

import (
	"testing"
)

func TestOption(t *testing.T) {
	tests := []struct {
		opt Option
		chk func(f *File) bool
	}{
		{
			opt: SetCutFormat(DayCut),
			chk: func(f *File) bool { return f.cutFmt == DayCut },
		},
		{
			opt: SetMaxSeq(1),
			chk: func(f *File) bool { return f.maxSeq == 1 },
		},
		{
			opt: SetMaxSize(2),
			chk: func(f *File) bool { return f.maxSize == 2 },
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
