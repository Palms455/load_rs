package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"load_rs/internal/load_file"
	"load_rs/configs"
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

	folder := configs.GetFolder()
	filech := make(chan string)
	files, err := load_file.GetFiles(folder)
	if err != nil {
		fmt.Println(err)
		return
	}
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(filech, &wg, i)
		fmt.Println("go worker ", i)
	}

	for _, filename := range files {
		filech <- filename
	}
	close(filech)

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Итоговое время: %s\n", elapsed)
	fmt.Println("Нажмите Enter для выхода из программы!")
	fmt.Scanln()
}
