package main

import (
	"io/ioutil"
	"os"
	"sync"

	vips "github.com/davidbyttow/govips/v2"
	"github.com/schollz/progressbar/v3"
	"github.com/viney-shih/goroutines"
)

func main() {
	path, _ := os.Getwd()
	files, _ := ioutil.ReadDir(path + "/demoImg")
	length := len(files)
	MainBar := NewBar(int64(length))

	wgA := new(sync.WaitGroup)

	vips.Startup(nil)
	defer vips.Shutdown()

	c := goroutines.NewBatch(4, goroutines.WithBatchSize(length))
	wgA.Add(length)
	defer c.Close()

	ep := vips.NewDefaultWEBPExportParams()
	ep.Quality = 50
	ep.Compression = 6
	ep.Lossless = false
	ep.Effort = 0

	for _, fileInfo := range files {
		fileName := fileInfo.Name()

		vipImg, _ := vips.NewImageFromFile("./demoImg/" + fileName)
		vipImg.AutoRotate()

		c.Queue(func() (interface{}, error) {
			im, _, _ := vipImg.Export(ep)

			ioutil.WriteFile("./result/"+fileName+".webp", im, 0644)

			MainBar.Add(1)

			vipImg.Close()

			return nil, nil
		})

	}

	wgA.Wait()
	c.QueueComplete()
}

func NewBar(length int64) *progressbar.ProgressBar {
	bar := progressbar.Default(length)
	return bar
}
