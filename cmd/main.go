package main

import (
	"fmt"
	"github.com/rjeczalik/notify"
	"load_rs/internal/load_file"
	"load_rs/internal/storage"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func flc_worker(ch chan int, wg *sync.WaitGroup) {
	for file_id := range ch {
		storage.FLC(file_id)
		runtime.Gosched()
	}
	wg.Done()
}

func worker(ch chan string, wg *sync.WaitGroup, i int) {

	for filename := range ch {
		time.Sleep(100*time.Millisecond)
		load_file.LoadRS(filename, i)
		runtime.Gosched()
	}
	wg.Done()
}


func main() {
	runtime.GOMAXPROCS(8)

	notify_chan := make(chan notify.EventInfo, 1000)
	folder := load_file.GetFolder()
	file_chan := make(chan string, 1000)

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
	}

	for i:= 0; i < 6; i++ {
		wg.Add(1)
		go flc_worker(storage.Flc_chan, &wg)
	}

	for _, filename := range files {
		file := filepath.Join(folder, filename)
		file_chan <- file
	}

	for res := range notify_chan {
			file_chan <- res.Path()
		}

	wg.Wait()
}