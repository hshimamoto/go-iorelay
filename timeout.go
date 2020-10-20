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
    timer *time.Timer
}

func NewTimeoutReadWriter(rw io.ReadWriter, timeout time.Duration) *TimeoutReadWriter {
    trw := &TimeoutReadWriter{rw: rw, timeout: timeout}
    trw.timer = time.NewTimer(time.Hour)
    if !trw.timer.Stop() {
	<-trw.timer.C
    }
    return trw
}

type ReadResult struct {
    n int
    err error
}

func (rw *TimeoutReadWriter)Read(p []byte) (int, error) {
    stopped := false
    defer func() {
	if !stopped && !rw.timer.Stop() {
	    <-rw.timer.C
	}
    }()
    rw.timer.Reset(rw.timeout)
    ch := make(chan ReadResult, 1)
    go func() {
	n, err := rw.rw.Read(p)
	ch <- ReadResult{n, err}
    }()
    for {
	select {
	case res := <-ch:
	    return res.n, res.err
	case <-rw.timer.C:
	    stopped = true
	    return 0, errors.New("iorelay Timeout")
	}
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
