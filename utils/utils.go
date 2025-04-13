package utils

import (
	"fmt"
	"io"
	"sync"
)

func NewErr(err error, format string, a ...any) error {
	return fmt.Errorf("%w; %s", err, fmt.Sprintf(format, a...))
}

var IoBuffPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 32*1024*2) // Two times of size in io package
	},
}

func CopyWithBuffer(dst io.Writer, src io.Reader) (written int64, err error) {
	buff := IoBuffPool.Get().([]byte)
	defer IoBuffPool.Put(buff)
	written, err = io.CopyBuffer(dst, src, buff)
	if err != nil {
		return
	}
	return written, nil
}

func CopyWithBufferN(dst io.Writer, src io.Reader, n int64) (written int64, err error) {
	written, err = CopyWithBuffer(dst, io.LimitReader(src, n))
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		err = io.EOF
	}
	return
}
