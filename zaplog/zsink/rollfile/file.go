package rollfile

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	BufferSize     = 256 * 1024
	PrintCreateLog = false
	LayoutPID      = false
	SeqLimit       = 1000
	pid            = os.Getpid()
)

const (
	dayLayout  = "20060102"
	hourLayout = "2006010215"
)

var (
	tickInterval  = 1 * time.Second
	flushInterval = 5 * time.Second
	touchInterval = 10 * time.Second
)

// 日志切割模式
type CutFormat string

// 日志切割模式常量定义
const (
	SizeCut CutFormat = "Size"
	HourCut CutFormat = "Hour"
	DayCut  CutFormat = "Day"
)

func isValidCutFormat(format CutFormat) bool {
	switch format {
	case SizeCut, HourCut, DayCut:
		return true
	}
	return false
}

type File struct {
	mu        sync.Mutex
	file      *os.File
	writer    *bufio.Writer
	seq       int
	size      int
	createdAt time.Time
	flushedAt time.Time
	touchedAt time.Time
	closed    bool
	done      chan struct{}

	dir     string
	name    string
	cutFmt  CutFormat
	maxSeq  int
	maxSize int
}

func Open(name string, opts ...Option) (*File, error) {
	f := &File{
		dir:     filepath.Dir(name),
		name:    filepath.Base(name),
		cutFmt:  SizeCut,
		maxSeq:  0,
		maxSize: 0,
	}
	for _, opt := range opts {
		opt(f)
	}
	if !isValidCutFormat(f.cutFmt) {
		return nil, &os.PathError{Op: "open", Path: name, Err: fmt.Errorf("invalid cut format %q", f.cutFmt)}
	}
	if f.maxSeq > SeqLimit {
		f.maxSeq = SeqLimit
	}

	if err := f.init(); err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	return f, nil
}

func (f *File) init() (err error) {
	// 1. 创建目录
	if err = createDir(f.dir); err != nil {
		return err
	}

	// 2. 打开文件
	if err = f.open(time.Now()); err != nil {
		return err
	}

	// 启动定时协程
	f.done = make(chan struct{})
	go f.running()

	return nil
}

func (f *File) baseName(t time.Time) string {
	switch f.cutFmt {
	case SizeCut:
		return sizeCutFileName(f.name, f.seq)
	case HourCut:
		return timeCutFileName(f.name, t, hourLayout)
	case DayCut:
		return timeCutFileName(f.name, t, dayLayout)
	}
	return sizeCutFileName(f.name, f.seq)
}

func (f *File) open(t time.Time) error {
	// 1. 读取 seq
	f.seq = readLinkSeq(f.dir, f.name)
	if f.seq < 0 || f.seq >= f.maxSeq {
		f.seq = 0
	}

	// 2. 打开文件
	base := f.baseName(t)
	file, err := openFile(f.dir, base, f.name)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	f.file = file
	f.writer = bufio.NewWriterSize(file, BufferSize)
	f.size = int(fi.Size())
	f.createdAt = t
	f.flushedAt = t
	f.touchedAt = t

	// 3. 输出文件打开日志
	if PrintCreateLog && f.size <= 0 {
		f.size, err = fmt.Fprintf(f.file, "Log file created at: %s\n", t.Format(time.RFC3339Nano))
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) create(t time.Time) error {
	// 1. 创建目标文件
	filename := f.baseName(t)
	file, err := createFile(f.dir, filename, f.name)
	if err != nil {
		return err
	}

	// 2. 关闭原文件
	if f.file != nil {
		f.writer.Flush()
		f.file.Close()
	}

	f.file = file
	f.writer = bufio.NewWriterSize(file, BufferSize)
	f.size = 0
	f.createdAt = t
	f.flushedAt = t
	f.touchedAt = t

	// 3. 输出文件打开日志
	if PrintCreateLog {
		f.size, err = fmt.Fprintf(f.file, "Log file created at: %s\n", t.Format(time.RFC3339Nano))
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) shouldFlush(t time.Time) bool {
	if t.Sub(f.flushedAt) < flushInterval {
		return false
	}
	return true
}

func (f *File) flush(t time.Time) error {
	f.flushedAt = t
	return f.writer.Flush()
}

func (f *File) shouldRotate(t time.Time) bool {
	switch f.cutFmt {
	case SizeCut:
		if f.maxSize > 0 && f.size >= f.maxSize {
			return true
		}
		return false
	case HourCut:
		if isSamePeriod(t, f.createdAt, time.Hour) {
			return false
		}
		return true
	case DayCut:
		if isSamePeriod(t, f.createdAt, 24*time.Hour) {
			return false
		}
		return true
	}
	return false
}

func (f *File) rotate(t time.Time) error {
	err := f.create(t)
	if err != nil {
		return err
	}

	f.seq++
	if f.seq < 0 || f.seq >= f.maxSeq {
		f.seq = 0
	}
	return nil
}

func (f *File) shouldTouch(t time.Time) bool {
	if t.Sub(f.touchedAt) < touchInterval {
		return false
	}
	return true
}

func (f *File) touch(t time.Time) error {
	f.touchedAt = t

	filename := f.baseName(t)
	file := filepath.Join(f.dir, filename)
	if !fileExist(file) {
		err := f.create(t)
		if err != nil {
			return err
		}
	}

	// 创建 link 文件
	link := filepath.Join(f.dir, f.name)
	if !fileExist(link) {
		os.Symlink(filename, link)
	}
	return nil
}

func (f *File) tick() {
	f.mu.Lock()
	defer f.mu.Unlock()

	var err error
	now := time.Now()

	// 1. 刷新缓冲
	if f.shouldFlush(now) {
		if err = f.flush(now); err != nil {
			fmt.Fprintf(os.Stderr, "rollfile.File: flush file: %v\n", err)
		}
	}

	// 2. 滚动文件
	if f.shouldRotate(now) {
		if err = f.rotate(now); err != nil {
			fmt.Fprintf(os.Stderr, "rollfile.File: rotate file: %v\n", err)
		}
	}

	// 3. 检测文件是否存在
	if f.shouldTouch(now) {
		if err = f.touch(now); err != nil {
			fmt.Fprintf(os.Stderr, "rollfile.File: touch file: %v\n", err)
		}
	}
}

func (f *File) running() {
	t := time.NewTicker(tickInterval)
	defer t.Stop()

	for {
		select {
		case <-f.done:
			return
		case <-t.C:
			f.tick()
		}
	}
}

func (f *File) wrapErr(op string, err error) error {
	return &os.PathError{Op: op, Path: f.name, Err: err}
}

func (f *File) Write(p []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 1. 是否已关闭
	if f.closed {
		return 0, f.wrapErr("write", os.ErrClosed)
	}

	// 2. 是否滚动文件
	now := time.Now()
	if f.shouldRotate(now) {
		if err = f.rotate(now); err != nil {
			return 0, f.wrapErr("write", fmt.Errorf("rotate: %w", err))
		}
	}

	// 3. 写入数据
	n, err = f.writer.Write(p)
	f.size += n
	if err != nil {
		return n, f.wrapErr("write", err)
	}
	return n, nil
}

func (f *File) Flush() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return f.wrapErr("flush", os.ErrClosed)
	}
	if err = f.flush(time.Now()); err != nil {
		return f.wrapErr("flush", err)
	}
	return nil
}

func (f *File) Sync() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return f.wrapErr("sync", os.ErrClosed)
	}
	f.flush(time.Now())
	if err = f.file.Sync(); err != nil {
		return f.wrapErr("sync", err)
	}
	return nil
}

func (f *File) Close() (err error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return f.wrapErr("close", os.ErrClosed)
	}
	f.closed = true
	close(f.done)
	f.writer.Flush()
	if err = f.file.Close(); err != nil {
		return f.wrapErr("close", err)
	}
	return nil
}
