package spa

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"time"
)

type modtimer interface {
	Mtime() time.Time
}

func NewBundledFS(fs fs.FS, prefix, mtime string) fs.FS {
	mtimeParsed, _ := time.Parse(time.RFC3339, mtime)
	return &spaBundledFs{FS: fs, prefix: prefix, mtime: mtimeParsed}
}

type spaBundledFs struct {
	fs.FS
	prefix string
	mtime  time.Time
}

func (fs *spaBundledFs) Open(name string) (fs.File, error) {

	if fs == nil {
		return nil, fmt.Errorf("bundled fs not loaded")
	}

	file, err := fs.FS.Open(path.Join(fs.prefix, name))
	if err != nil {
		return nil, err
	}

	return &spaBundledFile{File: file, fs: fs}, nil
}

func (fs *spaBundledFs) Mtime() time.Time {
	return fs.mtime
}

type spaBundledFile struct {
	fs.File
	fs *spaBundledFs
}

func (file *spaBundledFile) Stat() (fs.FileInfo, error) {

	if file == nil || file.fs == nil {
		return nil, fmt.Errorf("file disconnected from it's fs")
	}

	info, err := file.File.Stat()
	if err != nil {
		return nil, err
	}

	return &spaBundledFileInfo{FileInfo: info, fs: file.fs}, nil
}

func (file *spaBundledFile) Seek(offset int64, whence int) (int64, error) {

	seeker, ok := file.File.(io.Seeker)
	if !ok {
		return 0, fmt.Errorf("file is not seekable")
	}

	return seeker.Seek(offset, whence)
}

type spaBundledFileInfo struct {
	fs.FileInfo
	fs *spaBundledFs
}

func (info *spaBundledFileInfo) ModTime() time.Time {

	if info == nil || info.fs == nil {
		return time.Time{}
	}

	if mtime := info.fs.mtime; mtime.Year() > 2000 {
		return mtime
	}

	return time.Time{}
}
