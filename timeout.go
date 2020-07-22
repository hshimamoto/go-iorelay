// go-iorelay
// MIT License Copyright(c) 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:

package iorelay

import (
    "errors"
    "io"
    "time"
)

type TimeoutReadWriter struct {
    rw io.ReadWriter
    timeout time.Duration
}

func NewTimeoutReadWriter(rw io.ReadWriter, timeout time.Duration) *TimeoutReadWriter {
    return &TimeoutReadWriter{rw, timeout}
}

type ReadResult struct {
    n int
    err error
}

func (rw *TimeoutReadWriter)Read(p []byte) (int, error) {
    ch := make(chan ReadResult, 1)
    go func() {
	n, err := rw.rw.Read(p)
	ch <- ReadResult{n, err}
    }()
    select {
    case res := <-ch:
	return res.n, res.err
    case <-time.After(rw.timeout):
	return 0, errors.New("iorelay Timeout")
    }
}

func (rw *TimeoutReadWriter)Write(p []byte) (int, error) {
    return rw.rw.Write(p)
}

func RelayWithTimeout(io1, io2 io.ReadWriter, timeout time.Duration) {
    tio1 := NewTimeoutReadWriter(io1, timeout)
    tio2 := NewTimeoutReadWriter(io2, timeout)
    Relay(tio1, tio2)
}
