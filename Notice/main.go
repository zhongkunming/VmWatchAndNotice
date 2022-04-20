package main

import (
	"github.com/0xAX/notificator"
	"github.com/fsnotify/fsnotify"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
)

var notify = notificator.New(notificator.Options{
	DefaultIcon: "./icon/default.png",
	AppName:     "系统消息",
})

const logFile = "/Users/zhongkunming/VmFiles/ScreenshotAndCalculate/log/rtx_msg.txt"

const imgFile = "/Users/zhongkunming/VmFiles/ScreenshotAndCalculate/img/rtx_msg.png"

func main() {
	var watch, _ = fsnotify.NewWatcher()
	defer watch.Close()
	watch.Add(logFile)

	go func() {
		for {
			select {
			case ev := <-watch.Events:
				if ev.Op&fsnotify.Write == fsnotify.Write {
					notify.Push("RTX有消息", "", "./icon/default.png", notificator.UR_CRITICAL)
					go sendEmail()
				}
			}
		}
	}()

	select {}
}

func sendEmail() {
	from := "通知 <123@qq.com>"
	to := []string{"123@163.com"}
	username := "123@qq.com"
	password := ""

	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = from
	// 设置接收方的邮箱
	e.To = to
	//设置主题
	e.Subject = "RTX通知"
	//设置文件发送的内容
	e.Text = []byte("")
	e.AttachFile(imgFile)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", username, password, "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}
