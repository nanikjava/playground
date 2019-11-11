package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var str []string
var c chan bool
var readcomplete chan bool

var mutex = &sync.Mutex{}

func init() {
	str = make([]string, 1000)
	c = make(chan bool)
	readcomplete = make(chan bool)
}

func main(){
	go func(){
		l := len(str)
		for i := 0; i < l; i++ {
			str[i] = strconv.Itoa(i)
		}
		fmt.Println("Done poulating")
		c <- true

		// wait until read is done'
		select {
		case <-readcomplete:
			fmt.Println("Reading list is done")
		}
	}()

	select {
	case <-c:
		l := len(str)
		for i := 0; i < l; i++ {
			fmt.Println(" ", str[i])
		}
		readcomplete <- true

		// give a bit of time for the go routine
		// to complete
		time.Sleep(100*time.Millisecond)
	}
}
