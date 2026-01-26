package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type File struct {
	fileName string
	fileInfo os.FileInfo
}

func main() {
	//создаем канал для списка файлов
	filelist := make(chan File)
	//переменные для хранения изначльного пути и пути бэкапа
	var from, to string
	wg := &sync.WaitGroup{}
	fmt.Scan(&from, &to)

	wg.Add(1)
	go scanner(from, filelist, wg)

	go func() {
		wg.Wait()
		close(filelist)
	}()

	for v := range filelist {
		fmt.Println(v)
	}

}

func scanner(src string, ch chan<- File, wg *sync.WaitGroup) {
	defer wg.Done()
	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if path == src {
			return nil
		}
		if info.IsDir() {
			return nil
		} else {
			var n = File{
				fileName: info.Name(),
				fileInfo: info,
			}
			ch <- n
			return nil
		}

	})
}

func backuper(to string, ch <-chan File) {
	io.Copy(to, (<-ch).fileInfo)
}

func backuperPool() {

}

//scaner
//chan
