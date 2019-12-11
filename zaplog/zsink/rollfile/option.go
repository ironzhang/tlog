package rollfile

import "time"

type Option func(*File)

func SetLayout(layout string) Option {
	return func(f *File) {
		f.layout = layout
	}
}

func SetPeriod(period time.Duration) Option {
	return func(f *File) {
		if period < tickInterval {
			period = tickInterval
		}
		f.period = period
	}
}

func SetMaxSeq(maxSeq int) Option {
	return func(f *File) {
		f.maxSeq = maxSeq
	}
}

func SetMaxSize(maxSize int) Option {
	return func(f *File) {
		f.maxSize = maxSize
	}
}

func DisableCreateLog() Option {
	return func(f *File) {
		f.disableCreateLog = true
	}
}
