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
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/h2non/filetype"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/viney-shih/goroutines"
)

func main() {
	start := time.Now()
	spew.Dump(start)
	path, _ := os.Getwd()

	files, _ := ioutil.ReadDir(path + "/demoImg")

	length := len(files)
	batch := goroutines.NewBatch(10, goroutines.WithBatchSize(length))
	defer func() {
		batch.Close()
		end := time.Now()
		spew.Dump(end.Sub(start))
	}()

	for i, fileInfo := range files {
		tmp := fileInfo
		idx := i
		batch.Queue(func() (interface{}, error) {
			spew.Dump(tmp.Name())

			img, err := ImgDecode(tmp)
			if err != nil {
				log.Fatalln(err)
			}

			output, err := os.Create("./result/" + tmp.Name() + ".webp")
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, 75)
			if err != nil {
				log.Fatalln(err)
			}

			Img2Webp(output, img, options)

			fmt.Println(idx)
			return nil, nil
		})
	}

	batch.QueueComplete()

}

func Img2Webp(w io.Writer, src image.Image, options *encoder.Options) {
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
