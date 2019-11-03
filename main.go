package main

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/termnml/shutdowntimer/img"
)

type counter struct {
	out *widget.Label
	add *widget.Button
}

func newCounter() *counter {
	c := &counter{}
	c.out = widget.NewLabel("0")
	c.add = widget.NewButton("Add", func() {})
	return c
}

func setTime(clockCurrent *widget.Label) {
	timeFormatted := time.Now().Format("15:04:05")
	clockCurrent.SetText(timeFormatted)
}

func setTimeLeft(clockLeft *widget.Label) {
	timeLeftDuration := timeEnd.Sub(time.Now())
	timeLeftTime := time.Time{}.Add(timeLeftDuration)
	timeFormatted := timeLeftTime.Format("15:04:05")
	clockLeft.SetText(timeFormatted)
}

func setTimeEnd(clockEnd *widget.Label) {
	timeFormatted := timeEnd.Format("15:04:05")
	clockEnd.SetText(timeFormatted)
}

func addTimeToEnd(addTime time.Duration, valTimeLeft *widget.Label, valTimeShutdown *widget.Label) {
	timeEnd = timeEnd.Add(addTime)
	setTimeLeft(valTimeLeft)
	setTimeEnd(valTimeShutdown)
}

var timeLeft time.Duration
var timeEnd time.Time

func main() {
	app := app.New()
	w := app.NewWindow("Shutdown Timer")

	lblTimeCurrent := widget.NewLabelWithStyle("aktuelle Uhrzeit:",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
	)
	lblTimeLeft := widget.NewLabelWithStyle("verbleibende Zeit:",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
	)
	lblTimeShutdown := widget.NewLabelWithStyle("Shutdown Uhrzeit:",
		fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
	)
	valTimeCurrent := widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)
	valTimeLeft := widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)
	valTimeShutdown := widget.NewLabelWithStyle("",
		fyne.TextAlignTrailing, fyne.TextStyle{Bold: false},
	)

	setTime(valTimeCurrent)

	// dummy
	timeLeft = 120
	timeEnd = time.Now().Add(timeLeft * time.Second)

	setTimeEnd(valTimeShutdown)

	containerTop := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerTop.AddObject(widget.NewHBox(
		widget.NewVBox(
			lblTimeCurrent,
			lblTimeLeft,
			lblTimeShutdown,
		),
		widget.NewVBox(
			valTimeCurrent,
			valTimeLeft,
			valTimeShutdown,
		),
	))
	containerCenter := fyne.NewContainerWithLayout(layout.NewCenterLayout())
	containerCenter.AddObject(widget.NewHBox(
		widget.NewButton("+5 min", func() {
			addTimeToEnd(5*time.Minute, valTimeLeft, valTimeShutdown)
		}),
		widget.NewButton("+10 min", func() {
			addTimeToEnd(10*time.Minute, valTimeLeft, valTimeShutdown)
		}),
		widget.NewButton("+30 min", func() {
			addTimeToEnd(30*time.Minute, valTimeLeft, valTimeShutdown)
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

	w.SetContent(containerMaster)

	// settings for the window
	w.SetIcon(img.Icon)
	w.Resize(fyne.NewSize(400, 200))
	w.CenterOnScreen()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			setTime(valTimeCurrent)
			setTimeLeft(valTimeLeft)
		}
	}()

	w.ShowAndRun()
}
