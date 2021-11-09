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

	"github.com/h2non/filetype"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/schollz/progressbar/v3"
	"github.com/viney-shih/goroutines"
)

func main() {
	/*
		imageflow_tool v1/build --json examples/export_4_sizes/export_4_sizes.json
		        --in waterhouse.jpg
		        --out 1 waterhouse_w1600.jpg
		              2 waterhouse_w1200.jpg
		              3 waterhouse_w800.jpg
		              4 waterhouse_w400.jpg
		        --response operation_result.json

	*/

	//valueCmd := map[int][]string{}

	//valueCmd[1] = []string{"v1/build", "--json", "./config.json"}
	//valueCmd[2] = []string{"--in", "./demoImg/3489", "./demoImg/3491"}
	//valueCmd[3] = []string{"--out", "1", "result/3489_50.jpg", "3", "result/3661_50.jpg"}
	//valueCmd[4] = []string{"--response", "operation_result.json"}

	//io := make([]Io, 0)

	//var (
	//cmd *exec.Cmd
	//sum []string
	//)

	//for i := 1; i < 5; i++ {
	//sum = append(sum, valueCmd[i]...)
	//}
	//cmd = exec.Command("imageflow_tool", sum...)
	//spew.Dump(cmd)

	//_, err := cmd.CombinedOutput()
	//spew.Dump(err)

	path, _ := os.Getwd()

	files, _ := ioutil.ReadDir(path + "/demoImg")
	WebpBar := NewBar(int64(len(files)))
	//MainBar := NewBar(int64(len(files)))

	wgA := new(sync.WaitGroup)
	wgB := new(sync.WaitGroup)

	//p := goroutines.NewPool(20)

	length := len(files)
	b := goroutines.NewBatch(2, goroutines.WithBatchSize(2))
	c := goroutines.NewBatch(5, goroutines.WithBatchSize(length))
	defer b.Close()

	fileA := files[:length/2]
	fileB := files[length/2:]

	b.Queue(func() (interface{}, error) {
		for i, fileInfo := range fileA {

			img, err := ImgDecode(fileInfo)
			if err != nil {
				log.Fatalln(err)
			}

			output, err := os.Create("./result/" + fileInfo.Name() + ".webp")
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, 75)
			if err != nil {
				log.Fatalln(err)
			}

			wgA.Add(1)
			c.Queue(func() (interface{}, error) {
				Img2Webp(output, img, options, wgA, WebpBar)
				return nil, nil
			})
			fmt.Println(i)
			//MainBar.Add(1)
		}
		wgA.Wait()
		return nil, nil
	})

	b.Queue(func() (interface{}, error) {
		for i, fileInfo := range fileB {

			img, err := ImgDecode(fileInfo)
			if err != nil {
				log.Fatalln(err)
			}

			output, err := os.Create("./result/" + fileInfo.Name() + ".webp")
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, 75)
			if err != nil {
				log.Fatalln(err)
			}

			wgB.Add(1)
			c.Queue(func() (interface{}, error) {
				Img2Webp(output, img, options, wgB, WebpBar)
				return nil, nil
			})
			//MainBar.Add(1)
			fmt.Println(i)
		}

		wgB.Wait()
		return nil, nil
	})

	b.QueueComplete()
	//wgA.Wait()
	//wgB.Wait()

}

func NewBar(length int64) *progressbar.ProgressBar {

	bar := progressbar.Default(length)
	return bar

}

func Img2Webp(w io.Writer, src image.Image, options *encoder.Options, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	defer bar.Add(1)
	if err := webp.Encode(w, src, options); err != nil {
		log.Fatalln(err)
	}

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
