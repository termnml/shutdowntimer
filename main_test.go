package main

import (
	"testing"

	"fyne.io/fyne/test"
)

func TestCounterStartsAtZero(t *testing.T) {
	c := newCounter()

	if c.out.Text != "0" {
		t.Error("Counter doesnt start at 0")
	}
}

func TestCounterIncrement(t *testing.T) {
	c := newCounter()

	test.Tap(c.add)

	if c.out.Text != "1" {
		t.Error("Counter increment incorrect")
	}
}
