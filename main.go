package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
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
	wTimeCurrent.SetText(fmt.Sprintf("Uhrzeit:             %s", timeFormatted))

	timeLeftT := time.Time{}.Add(timeLeft)
	timeFormatted = timeLeftT.Format("15:04:05")
	wTimeLeft.SetText(fmt.Sprintf("verbleibend:   %s", timeFormatted))

	timeEnd := timeNow.Add(timeLeft)
	timeFormatted = timeEnd.Format("15:04:05")
	wTimeEnd.SetText(fmt.Sprintf("Ende:                 %s", timeFormatted))
}

func addTime(duration time.Duration) {
	timeLeft += duration
	outputTime()
}

func updateTime() {
	if runs {
		if timeLeft <= 0 {
			runs = false
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
var wTimeCurrent, wTimeLeft, wTimeEnd *widget.Label
var bStart, bStop *widget.Button
var runs bool

func main() {
	app := app.New()
	window := app.NewWindow("Shutdown Timer")

	wTimeCurrent = widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)
	wTimeLeft = widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: true},
	)
	wTimeEnd = widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)

	// dummyStart
	timeLeft = 10 * time.Second

	outputTime()

	containerTop := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerTop.AddObject(widget.NewHBox(
		widget.NewVBox(
			wTimeCurrent,
			wTimeLeft,
			wTimeEnd,
		),
	))
	containerCenter := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerCenter.AddObject(widget.NewHBox(
		widget.NewButton("+5 min", func() {
			addTime(5 * time.Minute)
		}),
		widget.NewButton("+10 min", func() {
			addTime(10 * time.Minute)
		}),
		widget.NewButton("+30 min", func() {
			addTime(30 * time.Minute)
		}),
	))

	bStart = widget.NewButton("START", func() {
		runs = true
		bStart.Disable()
		bStop.Enable()
	})
	bStop = widget.NewButton("STOP", func() {
		runs = false
		bStart.Enable()
		bStop.Disable()
	})
	runs = false
	bStop.Disable()

	containerBottom := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerBottom.AddObject(widget.NewHBox(
		bStart,
		bStop,
	))

	containerMaster := fyne.NewContainerWithLayout(layout.NewBorderLayout(containerTop, containerBottom, nil, nil))
	containerMaster.AddObject(containerTop)
	containerMaster.AddObject(containerBottom)
	containerMaster.AddObject(containerCenter)

	window.SetContent(containerMaster)

	// settings for the window
	window.SetIcon(img.Icon)
	window.Resize(fyne.NewSize(400, 200))
	window.CenterOnScreen()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			updateTime()
			outputTime()
		}
	}()

	window.ShowAndRun()
}
