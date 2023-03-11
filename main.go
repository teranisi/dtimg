package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/fogleman/gg"
)

// {"prob":0.8951893,"name":"bear","bounding-box":{"top":307,"left":237,"bottom":528,"right":492}},{"prob":0.7820245,"name":"bear","bounding-box":{"top":95,"left":60,"bottom":395,"right":210}}]
type Detection struct {
	Prob float64 `json:"prob"`
	Name string  `json:"name"`
	Box  BB      `json:"bounding-box"`
}

type BB struct {
	Top    int `json:"top"`
	Left   int `json:"left"`
	Bottom int `json:"bottom"`
	Right  int `json:"right"`
}

func main() {

	src := flag.String("s", "bear.jpg", "source file")
	filename := flag.String("j", "", "detected result json")

	flag.Parse()

	var r io.Reader
	switch *filename {
	case "", "-":
		r = os.Stdin
	default:
		f, err := os.Open(*filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		r = f
	}

	var det []Detection

	dat, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &det)

	if err != nil {
		panic(err)
	}
	//fmt.Printf("json: %v\n", det)

	im, err := gg.LoadImage(*src)
	if err != nil {
		panic(err)
	}

	m := im.Bounds().Max
	dc := gg.NewContext(m.X, m.Y)

	dc.Clear()
	dc.DrawImage(im, 0, 0)

	// font
	if err := dc.LoadFontFace("/System/Library/Fonts/Monaco.ttf", 18); err != nil {
		panic(err)
	}
	for _, e := range det {

		dc.SetRGBA(1.0, 0, 0, 1.0)
		dc.SetLineWidth(2.0)
		dc.DrawRoundedRectangle(float64(e.Box.Left), float64(e.Box.Top), float64(e.Box.Right-e.Box.Left), float64(e.Box.Bottom-e.Box.Top), 20)
		dc.Stroke()
		dc.SetRGBA(0.0, 0, 0, .2)
		dc.DrawRoundedRectangle(float64(e.Box.Left), float64(e.Box.Top), float64(e.Box.Right-e.Box.Left), float64(e.Box.Bottom-e.Box.Top), 15)
		dc.Fill()
		dc.SetRGB(1.0, 1.0, 1.0)
		//dc.DrawStringAnchored(fmt.Sprintf("%s (%.2f)", e.Name, e.Prob), float64(e.Box.Left), float64(e.Box.Top), -0.5, 1.0)
		dc.DrawString(fmt.Sprintf("%s (%.2f)", e.Name, e.Prob), float64(e.Box.Left)+10, float64(e.Box.Top)+20)
	}
	png.Encode(os.Stdout, dc.Image())

}
