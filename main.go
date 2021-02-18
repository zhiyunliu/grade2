package main

import (
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/golang/freetype"
)

var cnt = flag.Int("cnt", 120, "总共的计算题数量")

//go:embed arabtype.ttf
var fontBytes []byte

func main() {
	flag.Parse()
	createImage()
}
func createImage() {
	//图片的宽度
	dx := 1920
	//图片的高度
	dy := 1080
	imgfile, err := os.Create(time.Now().Format("2006-01-02_150405") + ".png")
	if err != nil {
		log.Println("1", err)
		return
	}
	defer imgfile.Close()

	img := image.NewNRGBA(image.Rect(0, 0, dx, dy))

	//读取字体数据
	// fontBytes, err := ioutil.ReadFile("arabtype.ttf")
	// if err != nil {
	// 	log.Println("2", err)
	// }
	//载入字体数据
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("3.load front fail", err)
	}
	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(120)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(26)
	f.SetClip(img.Bounds())
	//设置输出的图片
	f.SetDst(img)
	//设置字体颜色(红色)
	f.SetSrc(image.NewUniform(color.Black))

	writeContent(f)

	//以png 格式写入文件
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal("4", err)
	}
}

func writeContent(f *freetype.Context) {

	calcs := buildFormula(*cnt)
	w := 300
	cyc := 6
	adapter := 20
	var h, v int

	for i, clc := range calcs {
		if i%30 == 0 && i != 0 {
			pt := freetype.Pt(0, adapter+100+v*50+int(f.PointToFixed(26))>>8)
			f.DrawString(strings.Repeat("-", 200), pt)
		}
		h = i % cyc
		v = i / cyc
		//设置字体的位置
		pt := freetype.Pt(40+h*w, 60+v*50+int(f.PointToFixed(26))>>8)

		_, err := f.DrawString(clc, pt)
		if err != nil {
			log.Fatal("5", err)
		}
	}
}

//CC CC
type CC func() (int, int)

func buildFormula(cnt int) []string {
	symbols := []string{"+", "-", "×", "÷"}

	md := map[string]CC{
		"×": func() (int, int) {
			a := rand.Intn(100)
			time.Sleep(time.Millisecond * 20)
			b := rand.Intn(9) + 1
			return a, b
		},
		"÷": func() (int, int) {
			a := rand.Intn(20) + 1
			time.Sleep(time.Millisecond * 20)
			b := rand.Intn(9) + 1
			return a * (b + 1), b
		},
		"+": func() (int, int) {
			a := rand.Intn(100)
			time.Sleep(time.Millisecond * 20)

			b := rand.Intn(100)
			return a, b
		},
		"-": func() (int, int) {
			a := rand.Intn(100)
			time.Sleep(time.Millisecond * 20)
			b := rand.Intn(100)
			if a > b {
				return a, b
			}
			return b, a
		},
	}
	rows := []string{}
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := 0; i < cnt; i++ {
		p := rand.Intn(4)
		a, b := md[symbols[p]]()

		rows = append(rows, fmt.Sprintf("%d %s %d =    ", a, symbols[p], b))
	}
	return rows
}
