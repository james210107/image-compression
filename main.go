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
		ConcurrencyLevel: 2,
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
	bar := progressbar.NewOptions64(length,
		// 是否要顯示顏色
		progressbar.OptionEnableColorCodes(true),
		// 設置bar條長度
		//progressbar.OptionSetWidth(50),
		// 設置title
		progressbar.OptionSetDescription("壓爆圖片中..."),
		// 設置bar條樣式
		progressbar.OptionSetTheme(progressbar.Theme{
			// 進度條樣式
			Saucer: "[yellow]-[reset]",
			// 進度條的頭
			SaucerHead: "[yellow]─=≡Σ((つ•̀ω•́)つ[reset]",
			// 進度條還沒到的地方的樣式
			SaucerPadding: " ",
			// 進度條左邊框框
			BarStart: "[",
			// 進度條右邊框框
			BarEnd: "]",
		}),
		//progressbar.OptionFullWidth(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("張"),
	)
	return bar
}
