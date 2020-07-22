// go-iorelay/examples/fwd
// MIT License Copyright(c) 2019, 2020 Hiroshi Shimamoto
// vim:set sw=4 sts=4:

package main

import (
    "log"
    "net"
    "os"
    "time"

    "github.com/hshimamoto/go-iorelay"
)

func session(conn net.Conn, fwd string) {
    defer conn.Close()

    fconn, err := net.Dial("tcp", fwd)
    if err != nil {
	log.Printf("Dial %s %v\n", fwd, err)
	return
    }
    defer fconn.Close()

    iorelay.RelayWithTimeout(conn, fconn, time.Minute)
}

func main() {
    if len(os.Args) < 3 {
	log.Println("fwd listen dial")
	return
    }
    // start listening
    addr, err := net.ResolveTCPAddr("tcp", os.Args[1])
    if err != nil {
	log.Println("net.ResolveTCPAddr", err)
	return
    }
    l, err := net.ListenTCP("tcp", addr)
    if err != nil {
	log.Println("net.ListenTCP", err)
	return
    }
    defer l.Close()
    for {
	conn, err := l.AcceptTCP()
	if err != nil {
	    log.Println("AcceptTCP", err)
	    continue
	}
	go session(conn, os.Args[2])
    }
}
