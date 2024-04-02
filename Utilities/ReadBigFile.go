package Utilities

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sync"
	"time"
)

var K int64 = 1024
var M int64 = 1024 * K
var G int64 = 1024 * M
var BufSize = 512 * K

func ReadBigFile(savePath string, fh *multipart.FileHeader) error {
	fmt.Println("Starting upload video...")
	var cntGs int
	var cntDones int

	//新建file文件
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	fmt.Println("Create newFile Success!")

	defer file.Close()
	var limit int = 50
	//计算需要的缓存块数
	//BufSize := fh.Size / 50
	//numsBuf := (fh.Size + BufSize - 1) / BufSize

	//switch {
	//case fh.Size < 1024*1024*500:
	//	BufSize = 512 * K
	//	limit = 25
	//case fh.Size < 2*G:
	//	BufSize = 512 * K
	//	limit = 12
	//case fh.Size < 10*G:
	//	BufSize = 256 * K
	//	limit = 5
	//default:
	//	BufSize = 10 * M
	//	limit = 2
	//}
	numsBuf := (fh.Size + BufSize - 1) / BufSize
	fileReader, err := fh.Open()
	defer fileReader.Close()
	if err != nil {
		return err
	}
	bufReader := bufio.NewReader(fileReader)
	var wg sync.WaitGroup
	var ret error

	//Chunks := make(chan *Chunk, numsBuf)
	Limit := make(chan struct{}, limit)
	defer close(Limit)
	for i := int64(0); i < numsBuf; i++ {
		wg.Add(1)
		cntGs++

		//获得一块buffer
		buf := make([]byte, BufSize)
		n, err := io.ReadFull(bufReader, buf)
		fmt.Printf("read %d complete\n", i)
		Limit <- struct{}{} //控制并发度
		go func(idx int64, buf []byte) {
			fmt.Printf("Goroutine %d start:\n", idx)
			since := time.Now()
			defer func() {
				cntDones++
				wg.Done()
				fmt.Printf("Goroutine%d Done!,time:%v\n", idx, time.Since(since))
				<-Limit
			}()
			_, err = file.WriteAt(buf, idx*BufSize)
			if err != nil {
				//fmt.Println("err in writing:", err.Error())
				ret = err
				return
			}
		}(i, buf[:n])
	}
	fmt.Printf("cntGoroutines: %d\n", cntGs)
	wg.Wait()
	fmt.Println("Ending Wait")
	fmt.Println("BufSize = ", BufSize)

	return ret
}
