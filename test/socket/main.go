// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
    "net"
	"os"
	"io/ioutil"
	"time"
	"strconv"
)

func main() {
	number := 1
	if len(os.Args[1:]) > 0 {
		for _, arg := range os.Args[1:] {
			number, _ = strconv.Atoi(arg)
		}
	}
	fmt.Println(os.Args)
	fmt.Printf("%d requests\n", number)

	start := time.Now()
	ch := make(chan string)

	for i := 0; i < number; i++ {
		go fetch2(i, ch) // start a goroutine
	}
	for i := 0; i < number; i++ {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch2(index int, ch chan<- string) {
	start := time.Now()

    service := "10.224.22.186:80";
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf("index:%5d %.2fs err1:%v", index, secs, err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf("index:%5d %.2fs err2:%v", index, secs, err)
		return
	}

    _, err = conn.Write([]byte("GET /__query__ HTTP/1.0\r\n\r\n"))
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf("index:%5d %.2fs err3:%v", index, secs, err)
		return
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf("index:%5d %.2fs err4:%v", index, secs, err)
		return
	}

	secs := time.Since(start).Seconds()
	nbytes := len(string(result))
	ch <- fmt.Sprintf("index:%5d %.2fs nbytes:%d", index, secs, nbytes)
}

//!-
