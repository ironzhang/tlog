package rollfile

type Option func(*File)

func SetCutFormat(format CutFormat) Option {
	return func(f *File) {
		f.cutFmt = format
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
