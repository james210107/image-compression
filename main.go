package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	vips "github.com/davidbyttow/govips/v2"
	"github.com/schollz/progressbar/v3"
)

func main() {
	quality := flag.Int("q", 70, "Quality(int): default 70")
	worker := flag.Int("w", 4, "Number of Goroutine run at the same time(int): default 4")
	flag.Parse()
	fmt.Printf("%d %d", *quality, *worker)

	path, _ := os.Getwd()
	files, _ := ioutil.ReadDir(path + "/demoImg")
	length := len(files)
	MainBar := NewBar(int64(length))

	wgA := new(sync.WaitGroup)

	vips.Startup(&vips.Config{
		ConcurrencyLevel: 1,
		MaxCacheFiles:    0,
		MaxCacheMem:      512 * 1024,
		MaxCacheSize:     100,
		ReportLeaks:      false,
		CacheTrace:       false,
		CollectStats:     false,
	})
	defer vips.Shutdown()

	//c := goroutines.NewBatch(4, goroutines.WithBatchSize(length))
	wgA.Add(length)
	//defer c.Close()

	ep := vips.NewDefaultWEBPExportParams()
	ep.Quality = *quality
	ep.Lossless = false
	ep.Effort = 0

	limit := make(chan struct{}, *worker)

	for _, fileInfo := range files {
		limit <- struct{}{}
		fileName := fileInfo.Name()

		go func() (interface{}, error) {
			defer func() {
				<-limit
				wgA.Done()
			}()

			vipImg, _ := vips.NewImageFromFile("./demoImg/" + fileName)
			vipImg.AutoRotate()

			im, _, _ := vipImg.Export(ep)

			ioutil.WriteFile("./result/"+fileName+".webp", im, 0644)

			vipImg.Close()

			return nil, nil
		}()
		MainBar.Add(1)
	}

	wgA.Wait()
	//c.QueueComplete()
}

func NewBar(length int64) *progressbar.ProgressBar {
	bar := progressbar.Default(length)
	return bar
}
