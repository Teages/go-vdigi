package vdigi

import (
	"fmt"
	"syscall"
	"unsafe"
)

// Windows Platfrom
var (
	_PT_PEN = _POINTER_INPUT_TYPE(3)

	// Pointer Flags
	_POINTER_FLAG_NONE           = _POINTER_FLAGS(0x00000000) // Default
	_POINTER_FLAG_NEW            = _POINTER_FLAGS(0x00000001) // New pointer
	_POINTER_FLAG_INRANGE        = _POINTER_FLAGS(0x00000002) // Pointer has not departed
	_POINTER_FLAG_INCONTACT      = _POINTER_FLAGS(0x00000004) // Pointer is in contact
	_POINTER_FLAG_FIRSTBUTTON    = _POINTER_FLAGS(0x00000010) // Primary action
	_POINTER_FLAG_SECONDBUTTON   = _POINTER_FLAGS(0x00000020) // Secondary action
	_POINTER_FLAG_THIRDBUTTON    = _POINTER_FLAGS(0x00000040) // Third button
	_POINTER_FLAG_FOURTHBUTTON   = _POINTER_FLAGS(0x00000080) // Fourth button
	_POINTER_FLAG_FIFTHBUTTON    = _POINTER_FLAGS(0x00000100) // Fifth button
	_POINTER_FLAG_PRIMARY        = _POINTER_FLAGS(0x00002000) // Pointer is primary for system
	_POINTER_FLAG_CONFIDENCE     = _POINTER_FLAGS(0x00004000) // Pointer is considered unlikely to be accidental
	_POINTER_FLAG_CANCELED       = _POINTER_FLAGS(0x00008000) // Pointer is departing in an abnormal manner
	_POINTER_FLAG_DOWN           = _POINTER_FLAGS(0x00010000) // Pointer transitioned to down state (made contact)
	_POINTER_FLAG_UPDATE         = _POINTER_FLAGS(0x00020000) // Pointer update
	_POINTER_FLAG_UP             = _POINTER_FLAGS(0x00040000) // Pointer transitioned from down state (broke contact)
	_POINTER_FLAG_WHEEL          = _POINTER_FLAGS(0x00080000) // Vertical wheel
	_POINTER_FLAG_HWHEEL         = _POINTER_FLAGS(0x00100000) // Horizontal wheel
	_POINTER_FLAG_CAPTURECHANGED = _POINTER_FLAGS(0x00200000) // Lost capture
	_POINTER_FLAG_HASTRANSFORM   = _POINTER_FLAGS(0x00400000) // Input has a transform associated with it

	// PEN_FLAGS
	_PEN_FLAG_NONE     = _PEN_FLAGS(0x00000000) // Default
	_PEN_FLAG_BARREL   = _PEN_FLAGS(0x00000001) // The barrel button is pressed
	_PEN_FLAG_INVERTED = _PEN_FLAGS(0x00000002) // The pen is inverted
	_PEN_FLAG_ERASER   = _PEN_FLAGS(0x00000004) // The eraser button is pressed

	// PEN_MASK
	_PEN_MASK_NONE     = _PEN_MASK(0x00000000) // Default - none of the optional fields are valid
	_PEN_MASK_PRESSURE = _PEN_MASK(0x00000001) // The pressure field is valid
	_PEN_MASK_ROTATION = _PEN_MASK(0x00000002) // The rotation field is valid
	_PEN_MASK_TILT_X   = _PEN_MASK(0x00000004) // The tiltX field is valid
	_PEN_MASK_TILT_Y   = _PEN_MASK(0x00000008) // The tiltY field is valid
)

type _POINTER_TYPE_INFO struct {
	_type   _POINTER_INPUT_TYPE
	penInfo _POINTER_PEN_INFO
}

type _POINTER_PEN_INFO struct {
	pointerInfo _POINTER_INFO
	penFlags    _PEN_FLAGS
	penMask     _PEN_MASK
	pressure    uint32
	rotation    uint32
	tiltX       int32
	tiltY       int32
}

type _POINTER_INFO struct {
	pointerType        _POINTER_INPUT_TYPE
	pointerId          uint32
	frameId            uint32
	pointerFlags       _POINTER_FLAGS
	sourceDevice       _HANDLE
	hwndTarget         _HWND
	ptPixelLocation    _POINT
	_                  _POINT // ptHimetricLocation    POINT // for touch
	ptPixelLocationRaw _POINT
	_                  _POINT // ptHimetricLocationRaw POINT // for touch
	dwTime             _DWORD
	historyCount       uint32
	InputData          int32
	dwKeyStates        _DWORD
	PerformanceCount   uint64
	ButtonChangeType   _POINTER_BUTTON_CHANGE_TYPE
}

const (
	_SM_CMONITORS = 80
)

/* Copy from https://github.com/winlabs/gowin32/blob/master/wrappers/wingdi.go */
/* Apache License 2.0 */
/* START */
const (
	_CCHDEVICENAME                = 32
	_CCHFORMNAME                  = 32
	_ENUM_CURRENT_SETTINGS uint32 = 0xFFFFFFFF
)

type _DEVMODE struct {
	DeviceName    [_CCHDEVICENAME]uint16
	SpecVersion   uint16
	DriverVersion uint16
	Size          uint16
	DriverExtra   uint16
	Fields        uint32

	positionX          int32
	positionY          int32
	DisplayOrientation uint32
	DisplayFixedOutput uint32

	// Orientation      int16
	// PaperSize        int16
	// PaperLength      int16
	// PaperWidth       int16
	// Scale            int16
	// Copies           int16
	// DefaultSource    int16
	// PrintQuality     int16

	Color            int16
	Duplex           int16
	YResolution      int16
	TTOption         int16
	Collate          int16
	FormName         [_CCHFORMNAME]uint16
	LogPixels        uint16
	BitsPerPel       uint32
	PelsWidth        uint32
	PelsHeight       uint32
	DisplayFlags     uint32
	DisplayFrequency uint32
	ICMMethod        uint32
	ICMIntent        uint32
	MediaType        uint32
	DitherType       uint32
	Reserved1        uint32
	Reserved2        uint32
	PanningWidth     uint32
	PanningHeight    uint32
}

type _DISPLAY_DEVICE struct {
	Cb           uint32
	DeviceName   [32]uint16
	DeviceString [128]uint16
	StateFlags   uint32
	DeviceID     [128]uint16
	DeviceKey    [128]uint16
}

/* END */

type _HSYNTHETICPOINTERDEVICE _HANDLE

type _POINTER_INPUT_TYPE _DWORD
type _DWORD uint32
type _POINTER_FLAGS uint32
type _HANDLE uintptr
type _HWND uintptr
type _POINT struct {
	x, y int32
}
type _POINTER_BUTTON_CHANGE_TYPE int32

type _PEN_FLAGS uint32
type _PEN_MASK uint32

// Code
type _PointerDevice struct {
	device          _HSYNTHETICPOINTERDEVICE
	pointerTypeInfo _POINTER_TYPE_INFO
	_inContact      bool
	_lastContact    bool
}

var (
	user32dll   = syscall.NewLazyDLL("user32.dll")
	kernel32dll = syscall.NewLazyDLL("Kernel32.dll")

	createSyntheticPointerDevice  = user32dll.NewProc("CreateSyntheticPointerDevice")
	injectSyntheticPointerInput   = user32dll.NewProc("InjectSyntheticPointerInput")
	destroySyntheticPointerDevice = user32dll.NewProc("DestroySyntheticPointerDevice")
	getForegroundWindow           = user32dll.NewProc("GetForegroundWindow")

	getSystemMetrics      = user32dll.NewProc("GetSystemMetrics")
	enumDisplayDevices    = user32dll.NewProc("EnumDisplayDevicesW")
	enumDisplaySettingsEx = user32dll.NewProc("EnumDisplaySettingsExW")

	getLastErr = kernel32dll.NewProc("GetLastError")
)

func (device *_PointerDevice) Create() {
	d, _, _ := createSyntheticPointerDevice.Call(uintptr(_PT_PEN), uintptr(1), uintptr(2))
	// CreateSyntheticPointerDevice(PT_PEN, MAX_COUNT = 1, POINTER_FEEDBACK_INDIRECT);
	device.device = _HSYNTHETICPOINTERDEVICE(d)
	device.pointerTypeInfo = _POINTER_TYPE_INFO{
		_type: _PT_PEN,
	}
	setupDigiInfo(&device.pointerTypeInfo.penInfo)

	clearDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_NEW))
	sendDigiData(*device)
	clearDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_INRANGE|_POINTER_FLAG_PRIMARY))
}

func (device *_PointerDevice) Update(x, y int32, pressure uint32) error {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	if pressure > 0 {
		device.pointerTypeInfo.penInfo.pressure = pressure
		setDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_INCONTACT|_POINTER_FLAG_FIRSTBUTTON))
		device._inContact = true
	} else {
		device.pointerTypeInfo.penInfo.pressure = 1
		unsetDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_INCONTACT|_POINTER_FLAG_FIRSTBUTTON))
		device._inContact = false
	}

	device.pointerTypeInfo.penInfo.pointerInfo.ptPixelLocation = _POINT{x, y}
	device.pointerTypeInfo.penInfo.pointerInfo.ptPixelLocationRaw = _POINT{x, y}
	if device._inContact != device._lastContact {
		if device._inContact {
			unsetDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_UP|_POINTER_FLAG_UPDATE))
			setDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_DOWN))
		} else {
			unsetDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_DOWN|_POINTER_FLAG_UPDATE))
			setDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_UP))
		}
		device._lastContact = device._inContact
	} else {
		setDigiFlags(&device.pointerTypeInfo, _PEN_FLAGS(_POINTER_FLAG_UPDATE))
	}

	// if (subButton > 0 && device._inContact) {
	//     setPointerFlags(device.pointerTypeInfo, POINTER_FLAG_SECONDBUTTON);
	// } else {
	//     unsetPointerFlags(device.pointerTypeInfo, POINTER_FLAG_SECONDBUTTON);
	// }

	w, _, _ := getForegroundWindow.Call()
	device.pointerTypeInfo.penInfo.pointerInfo.hwndTarget = _HWND(w)
	return sendDigiData(*device)
}

func (device *_PointerDevice) Destroy() bool {
	ans, _, _ := destroySyntheticPointerDevice.Call(uintptr(device.device))
	return ans != 0
}

func setupDigiInfo(penInfo *_POINTER_PEN_INFO) {
	penInfo.pointerInfo.pointerType = _PT_PEN
	penInfo.pointerInfo.pointerId = 0
	penInfo.pointerInfo.frameId = 0
	penInfo.pointerInfo.pointerFlags = _POINTER_FLAG_NONE
	penInfo.pointerInfo.sourceDevice = 0
	w, _, _ := getForegroundWindow.Call()
	penInfo.pointerInfo.hwndTarget = _HWND(w)
	penInfo.pointerInfo.ptPixelLocation = _POINT{0, 0}
	penInfo.pointerInfo.ptPixelLocationRaw = _POINT{0, 0}
	penInfo.pointerInfo.dwTime = 0
	penInfo.pointerInfo.historyCount = 0
	penInfo.pointerInfo.dwKeyStates = 0
	penInfo.pointerInfo.PerformanceCount = 0
	penInfo.pointerInfo.ButtonChangeType = 0 // POINTER_CHANGE_NONE

	penInfo.penFlags = _PEN_FLAG_NONE
	penInfo.penMask = _PEN_MASK_PRESSURE
	penInfo.pressure = 512
	penInfo.rotation = 0
	penInfo.tiltX = 0
	penInfo.tiltY = 0
}

func clearDigiFlags(pointerTypeInfo *_POINTER_TYPE_INFO, flags _PEN_FLAGS) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags = _POINTER_FLAGS(flags)
}

// func cleanDigiFlags(pointerTypeInfo *POINTER_TYPE_INFO) {
// 	pointerTypeInfo.penInfo.pointerInfo.pointerFlags = 0
// }

func setDigiFlags(pointerTypeInfo *_POINTER_TYPE_INFO, flags _PEN_FLAGS) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags |= _POINTER_FLAGS(flags)
}
func unsetDigiFlags(pointerTypeInfo *_POINTER_TYPE_INFO, flags _PEN_FLAGS) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags &= ^_POINTER_FLAGS(flags)
}

func sendDigiData(device _PointerDevice) error {
	ans, _, _ := injectSyntheticPointerInput.Call(uintptr(device.device), uintptr(unsafe.Pointer(&device.pointerTypeInfo)), 1)

	if ans == 0 {
		err, _, _ := getLastErr.Call()
		return fmt.Errorf("injectSyntheticPointerInput: failed(%d)", err)
	}
	return nil
}

func nativeGetScreenCount() int {
	ans, _, _ := getSystemMetrics.Call(_SM_CMONITORS)
	return int(ans)
}

func nativeGetScreens() _VScreensData {
	count := nativeGetScreenCount()
	var screens []_Screen

	fixX := 0
	fixY := 0

	for id := 0; id < count; id++ {
		var device _DISPLAY_DEVICE
		device.Cb = uint32(uint64(unsafe.Sizeof(device)))
		enumDisplayDevices.Call(uintptr(0), uintptr(id), uintptr(unsafe.Pointer(&device)), uintptr(0))

		var deviceMode _DEVMODE
		deviceMode.Size = uint16(uint64(unsafe.Sizeof(deviceMode)))
		deviceMode.DriverExtra = 0
		enumDisplaySettingsEx.Call(uintptr(unsafe.Pointer(&device.DeviceName)), uintptr(_ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(&deviceMode)), uintptr(0))

		offsetX := int(deviceMode.positionX)
		offsetY := int(deviceMode.positionY)
		width := int(deviceMode.PelsWidth)
		height := int(deviceMode.PelsHeight)
		scaleX := 1.0
		scaleY := 1.0

		if offsetX < fixX {
			fixX = offsetX
		}
		if offsetY < fixY {
			fixY = offsetY
		}

		screen := _Screen{
			id, width, height, scaleX, scaleY, offsetX, offsetY,
		}
		screens = append(screens, screen)
	}

	totalWidth := 0
	totalHeight := 0

	for id := 0; id < count; id++ {
		screen := screens[id]

		screen.offsetX -= fixX
		screen.offsetY -= fixY

		if screen.offsetX+screen.width > totalWidth {
			totalWidth = screen.offsetX + screen.width
		}
		if screen.offsetY+screen.height > totalHeight {
			totalHeight = screen.offsetY + screen.height
		}

		screens[id] = screen
	}

	fmt.Printf("%d, %d\n", fixX, fixY)

	fmt.Printf("%d, %d\n", totalWidth, totalHeight)

	fmt.Printf("%+v\n\n", screens)

	return _VScreensData{
		count, screens, totalHeight, totalWidth,
	}
}
