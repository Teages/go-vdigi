package vdigi_test

import (
	"fmt"
	"testing"

	"github.com/Teages/go-vdigi"
)

func TestPlatfrom(t *testing.T) {
	// Create a digi devices
	d := vdigi.CreatePointer()

	// Update the position and pressure
	err := d.Update(20, 267, 0)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	// Destroy the digi
	d.Destroy()

	s := vdigi.GetScreens()

	t.Log(s.GetTotalSize())

	fmt.Printf("have %d screen(s)\n", s.GetScreenCount())
	for i := 0; i < s.GetScreenCount(); i++ {
		screen, _ := s.GetScreen(i)
		fmt.Printf("Screen %d: %v\n", i, screen)
	}

	dd := vdigi.CreatePointerForMainScreen()
	err = dd.Update(20, 20, 0)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	dd.Destroy()
}
