package main

import (
	"fmt"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/termnml/shutdowntimer/img"
)

func printTime() {
	timeFormatted := time.Now().Format("15:04:05")
	wTimeCurrent.SetText(fmt.Sprintf("Uhrzeit:             %s", timeFormatted))

	timeLeftDuration := timeEnd.Sub(time.Now())
	timeLeft := time.Time{}.Add(timeLeftDuration)
	timeFormatted = timeLeft.Format("15:04:05")
	wTimeLeft.SetText(fmt.Sprintf("verbleibend:   %s", timeFormatted))

	timeFormatted = timeEnd.Format("15:04:05")
	wTimeEnd.SetText(fmt.Sprintf("Ende:                 %s", timeFormatted))
}

func addTime(addTime time.Duration) {
	timeEnd = timeEnd.Add(addTime)
	printTime()
}

var timeLeft time.Duration
var timeEnd time.Time
var wTimeCurrent, wTimeLeft, wTimeEnd *widget.Label

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
	timeLeft = 120
	timeEnd = time.Now().Add(timeLeft * time.Second)

	printTime()

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

	containerBottom := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerBottom.AddObject(widget.NewHBox(
		widget.NewButton("START", func() {}),
		widget.NewButton("STOP", func() {}),
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
			printTime()
		}
	}()

	window.ShowAndRun()
}
