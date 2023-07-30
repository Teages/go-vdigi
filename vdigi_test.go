package vdigi_test

import (
	"fmt"
	"testing"

	"github.com/Teages/go-vdigi"
)

func TestPlatfrom(t *testing.T) {
	// test old api
	d := vdigi.CreatePointer()
	err := d.Update(20, 267, 0)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	d.Destroy()

	// multi screen api
	s := vdigi.GetScreens()
	t.Log(s.GetTotalSize())
	fmt.Printf("have %d screen(s)\n", s.GetScreenCount())
	for i := 0; i < s.GetScreenCount(); i++ {
		screen, _ := s.GetScreen(i)
		fmt.Printf("Screen %d: %v\n", i, screen)
	}
	dd, _ := vdigi.CreatePointerForScreen(0)
	err = dd.Update(100, 100, 0)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
	dd.Destroy()

	// screen uid
	mainScreen, _ := vdigi.GetScreens().GetScreen(0)
	mainScreenId, err := vdigi.GetScreens().GetScreenIdByUid(mainScreen.Uid)
	if err != nil || mainScreenId != 0 {
		t.Log(err.Error())
		t.Fail()
	}
}
