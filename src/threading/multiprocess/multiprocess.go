package main
import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type data struct {
	millis int64
	goroutinename string
}

var datalist []data
var totalgoroutines int
var doneprocessing chan bool
var mutex = &sync.Mutex{}
var counter int64


func init() {
	doneprocessing = make(chan bool)
	totalgoroutines = 100
	counter = 100
}

/**
This sample is to demonstrate the following scenario

* multi process steps
    - the data is a list of millis and goroutine name
    - the number of gouroutines is the same with the number of data
    - each goroutine will process the same data by adding it's name and millis
    - on completing the last step the app will print the complete output of the data
 */
func main(){
	// the logic is when the counter hit 0 that means we are done
	// the logic act as a watcher and it need to be ran first before
	// data processing kicks in
	go func(mtx *sync.Mutex){
		for {
			mutex.Lock()
			if (counter<=0) {
				doneprocessing <- true // this will indicate to the main function that processing is complete
			}
			mutex.Unlock()
		}
	}(mutex)

	for  i:=0;i<totalgoroutines;i++ {
		// spin up routines
		go func(mtx *sync.Mutex, rname string){
			mutex.Lock()
			atomic.AddInt64(&counter, -1)
			datalist= append(datalist, data{millis: makeTimestamp(), goroutinename: rname})
			mutex.Unlock()
		}(mutex, "goroutine-"+strconv.Itoa(i) )
	}

	<-doneprocessing
	for _, d := range datalist {
		fmt.Println(d)
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
