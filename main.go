package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	vips "github.com/davidbyttow/govips/v2"
	"github.com/h2non/filetype"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/schollz/progressbar/v3"
	"github.com/viney-shih/goroutines"
)

func main() {
	path, _ := os.Getwd()

	files, _ := ioutil.ReadDir(path + "/demoImg")
	//WebpBar := NewBar(int64(len(files)))
	MainBar := NewBar(int64(len(files)))

	wgA := new(sync.WaitGroup)

	//p := goroutines.NewPool(20)

	vips.Startup(nil)
	defer vips.Shutdown()

	length := len(files)
	c := goroutines.NewBatch(8, goroutines.WithBatchSize(length))
	defer c.Close()

	for _, fileInfo := range files {

		tmp := fileInfo.Name()
		wgA.Add(1)
		vipImg, _ := vips.NewImageFromFile("./demoImg/" + fileInfo.Name())

		vipImg.AutoRotate()

		ep := vips.NewDefaultWEBPExportParams()
		ep.Quality = 50
		ep.Compression = 6
		ep.Lossless = false
		ep.Effort = 0

		c.Queue(func() (interface{}, error) {
			im, _, err := vipImg.Export(ep)

			//img, err := ImgDecode(fileInfo)
			//if err != nil {
			//log.Fatalln(err)
			//}

			output, err := os.Create("./result/" + tmp + ".webp")
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			//options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, 50)
			//if err != nil {
			//log.Fatalln(err)
			//}

			ioutil.WriteFile("./result/"+tmp+".webp", im, 0644)

			MainBar.Add(1)
			//Img2Webp(output, img, options, wgA, WebpBar)
			vipImg.Close()
			return nil, nil
		})

	}
	wgA.Wait()
	c.QueueComplete()

	//wgA.Wait()
	//wgB.Wait()

}

func NewBar(length int64) *progressbar.ProgressBar {

	bar := progressbar.Default(length)
	return bar

}

func Img2Webp(w io.Writer, src image.Image, options *encoder.Options, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	//defer bar.Add(1)
	start := time.Now()
	if err := webp.Encode(w, src, options); err != nil {
		log.Fatalln(err)
	}
	dur := time.Now().Sub(start).Seconds()
	fmt.Printf("%v秒", dur)
	fmt.Println("")
}

func ImgDecode(fileInfo os.FileInfo) (image.Image, error) {

	file, err := os.Open("./demoImg/" + fileInfo.Name())

	if err != nil {
		log.Fatalln(err)
	}

	kind, _ := filetype.MatchReader(file)

	file, err = os.Open("./demoImg/" + fileInfo.Name())

	if err != nil {
		log.Fatalln(err)
	}

	switch kind.Extension {
	case "jpg", "jpeg":
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, err
	case "png":
		return png.Decode(file)
	case "webp":
		return webp.Decode(file, nil)
	default:
		fmt.Println(kind.Extension)
	}

	return nil, nil
}

type Config struct {
	Io        []Io      `json:"io"`
	Framewise Framewise `json:"framewise"`
}

type Io struct {
	ID        string `json:"io_id"`
	Direction string `json:"direction"`
	Io        string `json:"io"`
}

type Decode struct {
	ID string `json:"io_id"`
}

type Framewise struct {
	Graph Graph `json:"graph"`
}

type Graph struct {
	Nodes map[string]Node `json:"nodes"`
	Edges []Edge          `json:"edges"`
}

type Node struct {
	Decode Decode `json:"decode,omitempty"`
	Encode Encode `json:"encode,omitempty"`
}

type Edge struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Kind string `json:"input"`
}

type Encode struct {
	ID     string `json:"io_id"`
	Preset Preset `json:"preset"`
}

type Preset struct {
	MozImage MozImage `json:"mozjpeg"`
}

type MozImage struct {
	Quality     int  `json:"quality"`
	Progressive bool `json:"progressive"`
}
