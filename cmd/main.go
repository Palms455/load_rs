package main

import (
	"fmt"
	"github.com/rjeczalik/notify"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"load_rs/internal/load_file"
	"load_rs/configs"
	"log"
)


func worker(ch chan string, wg *sync.WaitGroup, i int) {

	for filename := range ch {
		time.Sleep(100*time.Millisecond)
		load_file.LoadRS(filename, i)
	}
	wg.Done()
}


func main() {
	runtime.GOMAXPROCS(8)

	notify_chan := make(chan notify.EventInfo, 1000)
	folder := configs.GetFolder()
	file_chan := make(chan string)

	if err := notify.Watch(folder, notify_chan, notify.Create); err != nil {
		log.Fatal(err)
	}
	defer close(file_chan)
	defer notify.Stop(notify_chan)

	files, err := load_file.GetFiles(folder)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(file_chan, &wg, i)
		fmt.Println("go worker ", i)
	}

	for _, filename := range files {
		file := filepath.Join(folder, filename)
		file_chan <- file
	}

	for {
		select {
		case res := <-notify_chan:
			file_chan <- res.Path()
		}
	}

	wg.Wait()
}