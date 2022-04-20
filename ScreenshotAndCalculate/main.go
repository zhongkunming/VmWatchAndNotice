package main

import (
	"ScreenshotAndCalculate/sys"
	"fmt"
	"github.com/0xAX/notificator"
	"github.com/corona10/goimagehash"
	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
	"image"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

// 截图后的像素（真实尺寸）
var pixelW = 0
var pixelY = 0

// 屏幕的大小
var sizeX = 0
var sizeY = 0

var plus = 0.0
var saveImg = true

// rtx 位置
var rtxX = 0
var rtxY = 0

var x1 = 0
var x2 = 0
var y1 = 0
var y2 = 0
var notify = notificator.New(notificator.Options{
	DefaultIcon: "./icon/default.png",
	AppName:     "系统通知",
})

func main() {
	getScreenPixelWidthAndHeight()
	getScreenSizeWidthAndHeight()
	getRTXCoordinates()
	go keepTheVirtualMachineAlive()
	go takeScreenshotAndCalculate()
	select {}
}

func getScreenPixelWidthAndHeight() {
	fullScreen, _ := screenshot.CaptureDisplay(0)
	pixelW = fullScreen.Bounds().Max.X
	pixelY = fullScreen.Bounds().Max.Y
	fmt.Printf("获取到屏幕像素级 宽度: %v, 高度: %v \n", pixelW, pixelY)
	if saveImg {
		savePicture(fullScreen, "./img/full.png")
	}
}

func getScreenSizeWidthAndHeight() {
	if runtime.GOOS != "windows" {
		notify.Push("点击屏幕最左侧以获取屏幕的宽度", "", "./icon/default.png", notificator.UR_CRITICAL)
	}
	sizeX, sizeY = sys.GetSystemMetrics()

	fmt.Printf("获取到屏幕尺寸级 宽度: %v, 高度: %v \n", sizeX, sizeY)
	plus = float64(pixelW) / float64(sizeX)
	fmt.Printf("像素和屏幕的比例是： %v \n", plus)
}

func getRTXCoordinates() {
	notify.Push("点击RTX图标以获取RTX图标的位置", "", "./icon/default.png", notificator.UR_CRITICAL)

	if robotgo.AddEvent("mleft") {
		rtxX, rtxY = robotgo.GetMousePos()
		fmt.Printf("RTX位置是： %v %v\n", rtxX, rtxY)
	}

	x1 = int(float64(rtxX-20) * plus)
	x2 = int(float64(rtxX+20) * plus)
	y1 = int(float64(rtxY-20) * plus)
	y2 = int(float64(rtxY+20) * plus)
}

func keepTheVirtualMachineAlive() {
	notify.Push("把鼠标焦点放在虚拟机内", "", "./icon/default.png", notificator.UR_CRITICAL)
	time.Sleep(10 * time.Second)
	for {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		if time.Now().UnixMilli()/2 == 0 {
			robotgo.ScrollSmooth(-10)
		} else {
			robotgo.ScrollSmooth(-10, 6, 200, -10)
		}
	}
}

func rtxScreenshot(index int) *image.RGBA {
	time.Sleep(time.Duration(index) * time.Second)
	rect := image.Rect(x1, y1, x2, y2)
	img, _ := screenshot.CaptureRect(rect)
	if saveImg {
		savePicture(img, "./img/rtx"+strconv.Itoa(index)+".png")
	}
	return img
}

func savePicture(img *image.RGBA, filePath string) {
	file, _ := os.Create(filePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_ = png.Encode(file, img)
}

func takeScreenshotAndCalculate() {

	for {
		img1 := rtxScreenshot(1)
		img2 := rtxScreenshot(2)
		img3 := rtxScreenshot(3)

		hash1, _ := goimagehash.AverageHash(img1)
		hash2, _ := goimagehash.AverageHash(img2)
		hash3, _ := goimagehash.AverageHash(img3)

		distance1, _ := hash1.Distance(hash2)
		distance2, _ := hash1.Distance(hash3)

		log, _ := os.Create("./log/rtx_msg.txt")
		if distance2 != distance1 {
			fmt.Printf("distance1 = %v, distance2 = %v \n", distance1, distance2)
			robotgo.Move(rtxX, rtxY)
			robotgo.Click("left", true)
			fullScreen, _ := screenshot.CaptureDisplay(0)
			savePicture(fullScreen, "./img/rtx_msg.png")
			log.WriteString("1")
			robotgo.MoveSmooth(sizeX/2, sizeY/2, 1.0, 2.0)
			robotgo.Click("left", true)
		}
	}

}
