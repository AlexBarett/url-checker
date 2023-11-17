package fileread

import (
	"bufio"
	"io"
	"os"
	"sync"
)

func ReadFormFile(file *os.File, queue *chan string, wg *sync.WaitGroup) ([]*string, error) {
	printOrder := make([]*string, 0)
	reader := bufio.NewReader(file)
	readWg := sync.WaitGroup{}
	var readErr error
	readWg.Add(1)

	go func() {
		for {
			url, err := reader.ReadString('\n')
			if err != nil {
				close(*queue)
				readWg.Done()

				if err == io.EOF {
					return
				}

				readErr = err
				return
			}

			if url[len(url)-1] == '\n' {
				url = url[:len(url)-1]
			}
			wg.Add(1)
			*queue <- url
			printOrder = append(printOrder, &url)
		}
	}()

	readWg.Wait()
	return printOrder, readErr
}
