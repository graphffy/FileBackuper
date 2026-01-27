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
	filePath string
}

func main() {
	filelist := make(chan File)
	var from, to string
	wg := &sync.WaitGroup{}
	Bwg := &sync.WaitGroup{}
	fmt.Scan(&from, &to)

	wg.Add(1)
	go scanner(from, filelist, wg)

	for i := 0; i < 3; i++ {
		Bwg.Add(1)
		go backuper(to, filelist, Bwg)
	}

	go func() {
		wg.Wait()
		close(filelist)
	}()
	Bwg.Wait()

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
				filePath: path,
			}
			ch <- n
			return nil
		}

	})
}

func backuper(to string, ch <-chan File, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range ch {

		src, err := os.Open(file.filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer src.Close()

		dst, err := os.Create(to + "/" + file.fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer dst.Close()
		io.Copy(dst, src)
	}
}

///home/graphffy/Desktop/folder/folderSrc
///home/graphffy/Desktop/folder/folderDst
