package vdigi

import (
	"errors"
	"fmt"
	"strconv"
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
	Uid    string
	Width  int
	Height int
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
	return p.UpdateWithTilt(x, y, pressure, 0, 0)
}

func (p Pointer) UpdateWithTilt(x, y int32, pressure uint32, tiltX, tiltY int32) error {
	if p.screenId < 0 {
		return p.d.Update(x, y, pressure, tiltX, tiltY)
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
		tiltX,
		tiltY,
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
		return ScreenData{"", 0, 0}, errors.New("Wrong screenId")
	}
	return ScreenData{
		getScreenUid(s.screens[screenId].offsetX, s.screens[screenId].offsetY),
		s.screens[screenId].width,
		s.screens[screenId].height,
	}, nil
}

func (s _VScreensData) GetScreenIdByUid(screenUid string) (int, error) {
	for i, screen := range s.screens {
		uid := getScreenUid(screen.offsetX, screen.offsetY)
		if uid == screenUid {
			return i, nil
		}
	}
	return 0, errors.New("Wrong screenId")
}

func getScreenUid(x, y int) string {
	return fmt.Sprintf("%03s%03s",
		strconv.FormatInt(int64(x), 16),
		strconv.FormatInt(int64(y), 16),
	)
}
