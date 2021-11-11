package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	vips "github.com/davidbyttow/govips/v2"
	"github.com/pterm/pterm"
	"github.com/schollz/progressbar/v3"
)

func main() {
	start := time.Now()
	newHeader := pterm.HeaderPrinter{
		TextStyle:       pterm.NewStyle(pterm.BgCyan),
		BackgroundStyle: pterm.NewStyle(pterm.FgBlack),
		Margin:          20,
	}

	// Print header.
	newHeader.Println("開始執行")
	//pterm.DefaultBigText.WithLetters(
	//pterm.NewLettersFromStringWithStyle("Tyr", pterm.NewStyle(pterm.FgLightBlue)),
	//pterm.NewLettersFromStringWithStyle("-ImageBroker", pterm.NewStyle(pterm.FgRed))).
	//Render()
	quality := flag.Int("q", 70, "Quality(int): default 70")
	worker := flag.Int("w", 2, "Number of Goroutine run at the same time(int): default 2")
	lossless := flag.Bool("l", false, "lossless")
	effort := flag.Int("e", 0, "effort")
	flag.Parse()
	fmt.Printf("quality: %d\nworker: %d\n", *quality, *worker)

	path, _ := os.Getwd()
	files, _ := ioutil.ReadDir(path + "/in")
	length := len(files)
	MainBar := NewBar(int64(length))

	wgA := new(sync.WaitGroup)

	vips.Startup(nil)
	defer vips.Shutdown()

	wgA.Add(length)

	ep := vips.NewDefaultWEBPExportParams()
	ep.Quality = *quality
	ep.Lossless = *lossless
	ep.Effort = *effort

	limit := make(chan struct{}, *worker)

	for _, fileInfo := range files {
		limit <- struct{}{}
		fileName := fileInfo.Name()

		go func() (interface{}, error) {
			defer func() {
				<-limit
				wgA.Done()
			}()

			vipImg, _ := vips.NewImageFromFile("./in/" + fileName)
			vipImg.AutoRotate()

			im, _, _ := vipImg.Export(ep)

			ioutil.WriteFile("./out/"+fileName+".webp", im, 0644)

			vipImg.Close()

			return nil, nil
		}()
		MainBar.Add(1)
	}

	wgA.Wait()
	end := time.Now()

	newHeader.Printfln(fmt.Sprintf("開始時間:  %s\n結束時間:  %v\n總花費秒數:%v秒", start.Local().Format(time.RFC3339), end.Local().Format(time.RFC3339), end.Sub(start).Seconds()))
}

func NewBar(length int64) *progressbar.ProgressBar {
	bar := progressbar.NewOptions64(length,
		// 是否要顯示顏色
		progressbar.OptionEnableColorCodes(true),
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
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("張"),
		progressbar.OptionSetPredictTime(true),
	)
	//bar := progressbar.Default(length)
	return bar
}
