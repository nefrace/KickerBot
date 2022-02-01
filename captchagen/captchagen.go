package captchagen

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
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
var XPositions []int = []int{}

// Создаём пустое изображение с серым градиентом, грузим шрифт из /assets и возвращаем контекст вызвавшей функции
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

// Генерация капчи.
//
// На пустое изображение наносятся логотипы из списка, предварительно перемешанного.
// К изображениям также добавляются порядковые номера (начиная с 1 вместо 0),
// а правильный вариант возвращается вместе с итоговой картинкой
func GenCaptcha() Captcha {
	dc := initImage()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(Logos), func(i, j int) { Logos[i], Logos[j] = Logos[j], Logos[i] })                          // Перемешиваем логотипы
	rand.Shuffle(len(XPositions), func(i, j int) { XPositions[i], XPositions[j] = XPositions[j], XPositions[i] }) // И позиции
	correct_answer := 0
	for i, logo := range Logos {
		x := XPositions[i]
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

func (captcha *Captcha) ToReader() *bytes.Reader {
	buff := new(bytes.Buffer)
	err := png.Encode(buff, captcha.Image)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}
	reader := bytes.NewReader(buff.Bytes())
	return reader
}

// Инициализация списка логотипов.
//
// Логотипы читаются из папки /assets рядом с исполняемым файлом.
// Принимается формат .png, логотип, представляющий правильный ответ называется godot.png
func Init() error {
	files, err := ioutil.ReadDir("./assets")
	if err != nil {
		return err
	}
	for _, file := range files {

		name := file.Name()
		if !strings.HasSuffix(name, ".png") { // Грузим только .png
			continue
		}
		log.Printf("%s", name)                   // Для отладки выводим имена файлов с логотипами
		path := fmt.Sprintf("./assets/%s", name) // Составляем путь до файла
		im, err := gg.LoadPNG(path)              // Грузим png, возвращаем ошибку если что-то идёт не так
		if err != nil {
			log.Print(err)
			return err
		}
		is_correct := strings.HasPrefix(name, "godot") // Если грузимый файл -- godot*.png - помечаем его как правильный
		logo := Logo{Image: im, IsCorrect: is_correct} // Создаём в памяти структуру для капчи
		Logos = append(Logos, logo)                    // Заносим логотип в общий список
	}

	for i := range Logos {
		XPositions = append(XPositions, 50+i*600/len(Logos)) // Горизонтальное расположение не рандомно: чтобы логотипы не перемешались.
	}
	return nil
}