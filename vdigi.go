package vdigi

import (
	"errors"
)

type Pointer struct {
	d        _PointerDevice
	screenId int
}

type _VScreensData struct {
	count       int
	screens     []_Screen
	totalWidth  int
	totalHeight int
}

type _Screen struct {
	id      int
	width   int
	height  int
	scaleX  float64
	scaleY  float64
	offsetX int
	offsetY int
}

type ScreenData struct {
	width  int
	height int
}

var (
	screensCache _VScreensData
)

func init() {
	screensCache = nativeGetScreens()
}

func CreatePointer() *Pointer {
	p := Pointer{}
	p.screenId = -1
	p.d.Create()
	return &p
}

func CreatePointerForScreen(screenId int) (*Pointer, error) {
	if screenId < 0 || screenId > screensCache.count {
		return nil, errors.New("wrong screenID")
	}
	p := Pointer{}
	p.screenId = screenId
	p.d.Create()
	return &p, nil
}

func CreatePointerForMainScreen() *Pointer {
	p, _ := CreatePointerForScreen(0)
	return p
}

func GetScreens() _VScreensData {
	return screensCache
}

func (p Pointer) Update(x, y int32, pressure uint32) error {
	if p.screenId < 0 {
		return p.d.Update(x, y, pressure)
	}

	screen := screensCache.screens[p.screenId]
	offsetX := screen.offsetX
	offsetY := screen.offsetY
	scaleX := screen.scaleX
	scaleY := screen.scaleY

	realX := int32(int(float64(x)*scaleX) + offsetX)
	realY := int32(int(float64(y)*scaleY) + offsetY)
	return p.d.Update(
		realX,
		realY,
		pressure,
	)
}

func (p Pointer) Destroy() {
	p.d.Destroy()
}

func (s _VScreensData) GetTotalSize() (int, int) {
	return s.totalWidth, s.totalHeight
}

func (s _VScreensData) GetScreenCount() int {
	return s.count
}

func (s _VScreensData) GetScreen(screenId int) (ScreenData, error) {
	if screenId < 0 || screenId > s.count {
		return ScreenData{0, 0}, errors.New("Wrong screenId")
	}
	return ScreenData{
		s.screens[screenId].width,
		s.screens[screenId].height,
	}, nil
}
