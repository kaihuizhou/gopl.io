// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 17.
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	// "io"
	"os"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

func main() {
	number := 1
	for _, arg := range os.Args[1:] {
		number, _ = strconv.Atoi(arg)
	}
	fmt.Println(os.Args)
	fmt.Printf("%d requests\n", number)

	start := time.Now()
	ch := make(chan string)
	// url := "http://10.224.22.186/__query__"
	url := "http://10.224.22.186/__ping__?action=create_meeting&site_id=700284272&user_id=488467417&meeting_id=156036662504279148&meeting_name=a&meeting_key=217743798&est_num=10&GDM_solution=0"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for i := 0; i < number; i++ {
		go fetch(i, client, url, ch) // start a goroutine
	}
	for i := 0; i < number; i++ {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(index int, client *http.Client, url string, ch chan<- string) {
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("index:%d err:%v", index, err) // send to channel ch
		return
	}

	// nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	responseData, err := ioutil.ReadAll(resp.Body)
	httpcode := resp.StatusCode
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("index:%d httpcode:%d while reading %s: err:%v", index, httpcode, url, err)
		return
	}

	responseString := string(responseData)
	nbytes := len(responseString)

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("index:%d httpcode:%d %.2fs  %7d  %s\n        body:%s", index, httpcode, secs, nbytes, url, responseString)
}

//!-
