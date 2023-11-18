package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"url-checker/cmd/app/internal/fileread"
	"url-checker/cmd/app/internal/filewrite"
	"url-checker/cmd/app/internal/request"
	"url-checker/cmd/app/internal/workerpool"
)

func main() {
	inputFileName := flag.String("input", "urls.txt", "Input file name")
	outputFileName := flag.String("output", "output.txt", "Output file name")
	retries := flag.Int("retries", 3, "Number of retries")
	timeout := flag.Int64("timeout", 2000, "Request timeout in ms")
	connectionLimit := flag.Int("connectionLimit", 0, "Number of max connections")
	flag.Parse()

	inputFile, err := os.Open(*inputFileName)
	if err != nil {
		panic(fmt.Errorf("Could not open file: %v; %v", inputFileName, err))
	}
	defer inputFile.Close()

	outputFile, err := os.Create(*outputFileName)
	if err != nil {
		panic(fmt.Errorf("Could not create file: %v; %v", outputFileName, err))
	}
	defer outputFile.Close()

	maxGorutines := *connectionLimit
	if maxGorutines == 0 {
		maxGorutines = runtime.NumCPU()
	}
	urlQueue := make(chan string, maxGorutines*2)

	wp := workerpool.New(maxGorutines)
	wg := sync.WaitGroup{}
	mutex := sync.RWMutex{}

	urlCompleteCount := 0

	client := request.New(*timeout, *retries)
	urlInfoMap := make(map[string]request.ResponseData)

	go func() {
		for url, ok := <-urlQueue; ok; url, ok = <-urlQueue {
			go func(url string) {
				wp.Exec(
					func() {
						info := client.GetRequestInfo(url)
						mutex.Lock()
						defer mutex.Unlock()
						urlInfoMap[url] = info
						urlCompleteCount++
						fmt.Printf("\rUrls complete: %v", urlCompleteCount)
						wg.Done()
					},
				)
			}(url)
		}
	}()

	urlsOrder, err := fileread.ReadFormFile(inputFile, &urlQueue, &wg)
	if err != nil {
		panic(fmt.Errorf("Read file error: %v", err))
	}

	wg.Wait()
	fmt.Println()

	filewrite.Print(outputFile, urlsOrder, urlInfoMap)
	fmt.Println("Done")
}
