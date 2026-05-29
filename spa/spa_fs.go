package spa

import (
	"fmt"
	"io"
	"io/fs"
	"time"
)

func NewBundledFS(fs fs.FS, rfc3339 string) fs.FS {
	date, _ := time.Parse(time.RFC3339, rfc3339)
	return &spaBundledFs{FS: fs, mtime: date}
}

type spaBundledFs struct {
	fs.FS
	mtime time.Time
}

func (fs *spaBundledFs) Open(name string) (fs.File, error) {

	if fs == nil {
		return nil, fmt.Errorf("bundled fs not loaded")
	}

	file, err := fs.FS.Open(name)
	if err != nil {
		return nil, err
	}

	return &spaBundledFile{File: file, fs: fs}, nil
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
