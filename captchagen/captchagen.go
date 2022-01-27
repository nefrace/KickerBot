package captchagen

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type Logo struct {
	Image     image.Image
	IsCorrect bool
}

type Captcha struct {
	Image         image.Image
	CorrectAnswer int
}

var Logos []Logo = []Logo{}

func initImage() *gg.Context {
	dc := gg.NewContext(600, 400)
	if err := dc.LoadFontFace("./assets/font.ttf", 24); err != nil {
		panic(err)
	}
	grad := gg.NewLinearGradient(0, 0, 600, 400)
	grad.AddColorStop(0, color.RGBA{71, 100, 106, 255})
	grad.AddColorStop(1, color.RGBA{44, 43, 51, 255})

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, 600, 400)
	dc.Fill()
	return dc
}

func GenCaptcha() Captcha {
	dc := initImage()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(Logos), func(i, j int) { Logos[i], Logos[j] = Logos[j], Logos[i] })
	correct_answer := 0
	count := len(Logos)
	for i, logo := range Logos {
		x := i*600/count + 50
		y := rand.Intn(400 - logo.Image.Bounds().Dy())
		dc.DrawImage(logo.Image, x, y)
		if logo.IsCorrect {
			correct_answer = i + 1
		}
		var tx float64 = float64(x) + 50.0
		var ty float64 = float64(y) + 120.0
		text := fmt.Sprintf("%d", i+1)
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(text, tx, ty, 0.5, 0.5)
	}
	img := dc.Image()
	captcha := Captcha{
		Image:         img,
		CorrectAnswer: correct_answer,
	}

	return captcha
}

func InitImages() error {
	images := []Logo{}
	files, err := ioutil.ReadDir("./assets")
	if err != nil {
		return err
	}
	for _, file := range files {

		name := file.Name()
		if !strings.HasSuffix(name, ".png") {
			continue
		}
		log.Printf("%s", name)
		path := fmt.Sprintf("./assets/%s", name)
		im, err := gg.LoadPNG(path)
		if err != nil {
			log.Print(err)
			return err
		}
		is_correct := strings.HasPrefix(name, "godot")
		logo := Logo{Image: im, IsCorrect: is_correct}
		images = append(images, logo)
	}
	Logos = images
	log.Printf("%v", Logos)
	return nil
}
