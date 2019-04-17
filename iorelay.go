// go-iorelay
// MIT License Copyright(c) 2019 Hiroshi Shimamoto
// vim:set sw=4 sts=4:

package iorelay

import (
    "io"
    "time"
)

func fwding(dst, src io.ReadWriter) chan struct{} {
    done := make(chan struct{})
    go func() {
	io.Copy(dst, src)
	done <- struct{}{}
    }()
    return done
}

func wait(done chan struct{}, timeout time.Duration) {
    select {
    case <-done:
    case <-time.After(timeout):
    }
}

func Relay(io1, io2 io.ReadWriter) {
    done1 := fwding(io1, io2)
    done2 := fwding(io2, io1)
    timeout := 5 * time.Second
    select {
    case <-done1: wait(done2, timeout)
    case <-done2: wait(done1, timeout)
    }
}
