package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	KB = 1 << (10 * (iota + 1))
	MB
	GB
	TB
)

var (
	DataDir string
)

// CopyFile copy file
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = CreateNestedFile(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = CopyWithBuffer(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copy dir
func CopyDir(src, dst string) error {
	var err error
	var fds []os.DirEntry
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	if fds, err = os.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// CreateNestedFile create nested file
func CreateNestedFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if err := os.MkdirAll(basePath, 0700); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func CreateTempFile(r io.Reader, size int64) (*os.File, error) {
	if f, ok := r.(*os.File); ok {
		return f, nil
	}
	f, err := os.CreateTemp(filepath.Join(DataDir, "temp"), "file-*")
	if err != nil {
		return nil, err
	}
	readBytes, err := CopyWithBuffer(f, r)
	if err != nil {
		_ = os.Remove(f.Name())
		return nil, NewErr(err, "CreateTempFile failed")
	}
	if size > 0 && readBytes != size {
		_ = os.Remove(f.Name())
		return nil, NewErr(err, "CreateTempFile failed, incoming stream actual size= %d, expect = %d ", readBytes, size)
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		_ = os.Remove(f.Name())
		return nil, NewErr(err, "CreateTempFile failed, can't seek to 0 ")
	}
	return f, nil
}

func GetFileSize(path string) (int64, error) {
	f, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

func GetFileMode(path string) (os.FileMode, error) {
	f, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return f.Mode(), nil
}

func GetFileModTime(path string) (time.Time, error) {
	f, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return f.ModTime(), nil
}
