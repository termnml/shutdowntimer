package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/termnml/shutdowntimer/img"
)

func outputTime() {
	timeNow := time.Now()
	timeFormatted := timeNow.Format("15:04:05")
	widTimeCurrent.SetText(fmt.Sprintf("Uhrzeit:             %s", timeFormatted))

	timeLeftT := time.Time{}.Add(timeLeft)
	timeFormatted = timeLeftT.Format("15:04:05")
	widTimeLeft.SetText(fmt.Sprintf("verbleibend:   %s", timeFormatted))

	timeEnd := timeNow.Add(timeLeft)
	timeFormatted = timeEnd.Format("15:04:05")
	widTimeEnd.SetText(fmt.Sprintf("Ende:                 %s", timeFormatted))
}

func addTime(duration time.Duration) {
	timeLeft += duration
	outputTime()
}

func setTime() {
	hour, err1 := strconv.Atoi(inTimeH.Text)
	min, err2 := strconv.Atoi(inTimeM.Text)
	if err1 != nil || err2 != nil {
		inTimeH.SetText("00")
		inTimeM.SetText("00")
		return
	}

	timeLeft = time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute
	outputTime()
}

func updateTime() {
	if timerIsRunning {
		if timeLeft <= 0 {
			timerIsRunning = false
			triggerShutdown()
		} else {
			timeLeft -= time.Second
		}
	}
}

func triggerShutdown() {
	if runtime.GOOS == "windows" {
		if err := exec.Command("cmd", "/C", "shutdown", "-s", "-f", "-t", "00").Run(); err != nil {
			fmt.Println("Failed to initiate shutdown:", err)
			log.Fatal(err)
		}
	}
	if runtime.GOOS == "linux" {
		if err := exec.Command("bash", "-c", "shutdown", "now").Run(); err != nil {
			fmt.Println("Failed to initiate shutdown:", err)
			log.Fatal(err)
		}
	}
}

var timeLeft time.Duration

var widTimeCurrent, widTimeLeft, widTimeEnd *widget.Label
var btnStart, btnStop, btnSet *widget.Button
var inTimeH, inTimeM *widget.Entry

var timerIsRunning bool

func main() {
	app := app.New()
	window := app.NewWindow("Shutdown Timer")

	// dummyStart
	timeLeft = 0

	widTimeCurrent = widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)
	widTimeLeft = widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: true},
	)
	widTimeEnd = widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)
	containerTop := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerTop.AddObject(widget.NewHBox(
		widget.NewVBox(
			widTimeCurrent,
			widTimeLeft,
			widTimeEnd,
		),
	))

	inTimeH = widget.NewEntry()
	inTimeH.SetText("00")
	inTimeM = widget.NewEntry()
	inTimeM.SetText("00")

	btnSet = widget.NewButton("Zeit setzen", func() {
		setTime()
	})

	containerCenter := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerCenter.AddObject(widget.NewVBox(
		widget.NewHBox(
			widget.NewButton("+5 min", func() {
				addTime(5 * time.Minute)
			}),
			widget.NewButton("+10 min", func() {
				addTime(10 * time.Minute)
			}),
			widget.NewButton("+30 min", func() {
				addTime(30 * time.Minute)
			}),
		),
		widget.NewHBox(
			inTimeH, widget.NewLabel("stunden"),
			inTimeM, widget.NewLabel("minuten"),
			btnSet,
		),
	))

	btnStart = widget.NewButton("START", func() {
		if timeLeft > 0 {
			timerIsRunning = true
			btnStart.Disable()
			btnStop.Enable()
		}
	})
	btnStop = widget.NewButton("STOP", func() {
		timerIsRunning = false
		btnStart.Enable()
		btnStop.Disable()
	})
	timerIsRunning = false
	btnStop.Disable()

	containerBottom := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerBottom.AddObject(widget.NewHBox(
		btnStart,
		btnStop,
	))

	containerMaster := fyne.NewContainerWithLayout(layout.NewBorderLayout(containerTop, containerBottom, nil, nil))
	containerMaster.AddObject(containerTop)
	containerMaster.AddObject(containerBottom)
	containerMaster.AddObject(containerCenter)

	window.SetContent(containerMaster)

	// settings for the window
	window.SetIcon(img.Icon)
	window.Resize(fyne.NewSize(400, 150))
	window.CenterOnScreen()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			updateTime()
			outputTime()
		}
	}()

	outputTime()
	window.ShowAndRun()
}
