package vdigi_test

import (
	"testing"

	"github.com/Teages/go-vdigi"
)

func TestPlatfrom(t *testing.T) {
	// Create a digi devices
	d := vdigi.CreatePointer()

	// Update the position and pressure
	err := d.Update(100, 100, 0)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	// Destroy the digi
	d.Destroy()
}
