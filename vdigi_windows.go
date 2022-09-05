package vdigi

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows Platfrom
var (
	PT_PEN = POINTER_INPUT_TYPE(3)

	// Pointer Flags
	POINTER_FLAG_NONE           = POINTER_FLAGS(0x00000000) // Default
	POINTER_FLAG_NEW            = POINTER_FLAGS(0x00000001) // New pointer
	POINTER_FLAG_INRANGE        = POINTER_FLAGS(0x00000002) // Pointer has not departed
	POINTER_FLAG_INCONTACT      = POINTER_FLAGS(0x00000004) // Pointer is in contact
	POINTER_FLAG_FIRSTBUTTON    = POINTER_FLAGS(0x00000010) // Primary action
	POINTER_FLAG_SECONDBUTTON   = POINTER_FLAGS(0x00000020) // Secondary action
	POINTER_FLAG_THIRDBUTTON    = POINTER_FLAGS(0x00000040) // Third button
	POINTER_FLAG_FOURTHBUTTON   = POINTER_FLAGS(0x00000080) // Fourth button
	POINTER_FLAG_FIFTHBUTTON    = POINTER_FLAGS(0x00000100) // Fifth button
	POINTER_FLAG_PRIMARY        = POINTER_FLAGS(0x00002000) // Pointer is primary for system
	POINTER_FLAG_CONFIDENCE     = POINTER_FLAGS(0x00004000) // Pointer is considered unlikely to be accidental
	POINTER_FLAG_CANCELED       = POINTER_FLAGS(0x00008000) // Pointer is departing in an abnormal manner
	POINTER_FLAG_DOWN           = POINTER_FLAGS(0x00010000) // Pointer transitioned to down state (made contact)
	POINTER_FLAG_UPDATE         = POINTER_FLAGS(0x00020000) // Pointer update
	POINTER_FLAG_UP             = POINTER_FLAGS(0x00040000) // Pointer transitioned from down state (broke contact)
	POINTER_FLAG_WHEEL          = POINTER_FLAGS(0x00080000) // Vertical wheel
	POINTER_FLAG_HWHEEL         = POINTER_FLAGS(0x00100000) // Horizontal wheel
	POINTER_FLAG_CAPTURECHANGED = POINTER_FLAGS(0x00200000) // Lost capture
	POINTER_FLAG_HASTRANSFORM   = POINTER_FLAGS(0x00400000) // Input has a transform associated with it

	// PEN_FLAGS
	PEN_FLAG_NONE     = PEN_FLAGS(0x00000000) // Default
	PEN_FLAG_BARREL   = PEN_FLAGS(0x00000001) // The barrel button is pressed
	PEN_FLAG_INVERTED = PEN_FLAGS(0x00000002) // The pen is inverted
	PEN_FLAG_ERASER   = PEN_FLAGS(0x00000004) // The eraser button is pressed

	// PEN_MASK
	PEN_MASK_NONE     = PEN_MASK(0x00000000) // Default - none of the optional fields are valid
	PEN_MASK_PRESSURE = PEN_MASK(0x00000001) // The pressure field is valid
	PEN_MASK_ROTATION = PEN_MASK(0x00000002) // The rotation field is valid
	PEN_MASK_TILT_X   = PEN_MASK(0x00000004) // The tiltX field is valid
	PEN_MASK_TILT_Y   = PEN_MASK(0x00000008) // The tiltY field is valid
)

type POINTER_TYPE_INFO struct {
	_type   POINTER_INPUT_TYPE
	penInfo POINTER_PEN_INFO
}

type POINTER_PEN_INFO struct {
	pointerInfo POINTER_INFO
	penFlags    PEN_FLAGS
	penMask     PEN_MASK
	pressure    uint32
	rotation    uint32
	tiltX       int32
	tiltY       int32
}

type POINTER_INFO struct {
	pointerType        POINTER_INPUT_TYPE
	pointerId          uint32
	frameId            uint32
	pointerFlags       POINTER_FLAGS
	sourceDevice       HANDLE
	hwndTarget         HWND
	ptPixelLocation    POINT
	_                  POINT // ptHimetricLocation    POINT // for touch
	ptPixelLocationRaw POINT
	_                  POINT // ptHimetricLocationRaw POINT // for touch
	dwTime             DWORD
	historyCount       uint32
	InputData          int32
	dwKeyStates        DWORD
	PerformanceCount   uint64
	ButtonChangeType   POINTER_BUTTON_CHANGE_TYPE
}

type HSYNTHETICPOINTERDEVICE HANDLE

type POINTER_INPUT_TYPE DWORD
type DWORD uint32
type POINTER_FLAGS uint32
type HANDLE windows.Handle
type HWND windows.HWND
type POINT struct {
	x, y int32
}
type POINTER_BUTTON_CHANGE_TYPE int32

type PEN_FLAGS uint32
type PEN_MASK uint32

// Code
type PointerDevice struct {
	device          HSYNTHETICPOINTERDEVICE
	pointerTypeInfo POINTER_TYPE_INFO
	_inContact      bool
	_lastContact    bool
}

var (
	// Library
	libuser32   *windows.LazyDLL
	libkernel32 *windows.LazyDLL
	libtest     *windows.LazyDLL

	// Functions
	createSyntheticPointerDevice  *windows.LazyProc
	injectSyntheticPointerInput   *windows.LazyProc
	destroySyntheticPointerDevice *windows.LazyProc
	getForegroundWindow           *windows.LazyProc

	// Error Function
	getLastErr *windows.LazyProc
	tstStru    *windows.LazyProc
)

func init() {
	libuser32 = windows.NewLazySystemDLL("user32.dll")
	libkernel32 = windows.NewLazySystemDLL("Kernel32.dll")

	createSyntheticPointerDevice = libuser32.NewProc("CreateSyntheticPointerDevice")
	injectSyntheticPointerInput = libuser32.NewProc("InjectSyntheticPointerInput")
	destroySyntheticPointerDevice = libuser32.NewProc("DestroySyntheticPointerDevice")
	getForegroundWindow = libuser32.NewProc("GetForegroundWindow")

	getLastErr = libkernel32.NewProc("GetLastError")
}

func (device *PointerDevice) Create() {
	d, _, _ := createSyntheticPointerDevice.Call(uintptr(PT_PEN), uintptr(1), uintptr(2))
	// CreateSyntheticPointerDevice(PT_PEN, MAX_COUNT = 1, POINTER_FEEDBACK_INDIRECT);
	device.device = HSYNTHETICPOINTERDEVICE(d)
	device.pointerTypeInfo = POINTER_TYPE_INFO{
		_type: PT_PEN,
	}
	setupDigiInfo(&device.pointerTypeInfo.penInfo)

	clearDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_NEW))
	sendDigiData(*device)
	clearDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_INRANGE|POINTER_FLAG_PRIMARY))
}

func (device *PointerDevice) Update(x, y int32, pressure uint32) error {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	if pressure > 0 {
		device.pointerTypeInfo.penInfo.pressure = pressure
		setDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_INCONTACT|POINTER_FLAG_FIRSTBUTTON))
		device._inContact = true
	} else {
		device.pointerTypeInfo.penInfo.pressure = 1
		unsetDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_INCONTACT|POINTER_FLAG_FIRSTBUTTON))
		device._inContact = false
	}

	device.pointerTypeInfo.penInfo.pointerInfo.ptPixelLocation = POINT{x, y}
	device.pointerTypeInfo.penInfo.pointerInfo.ptPixelLocationRaw = POINT{x, y}
	if device._inContact != device._lastContact {
		if device._inContact {
			unsetDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_UP|POINTER_FLAG_UPDATE))
			setDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_DOWN))
		} else {
			unsetDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_DOWN|POINTER_FLAG_UPDATE))
			setDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_UP))
		}
		device._lastContact = device._inContact
	} else {
		setDigiFlags(&device.pointerTypeInfo, PEN_FLAGS(POINTER_FLAG_UPDATE))
	}

	// if (subButton > 0 && device._inContact) {
	//     setPointerFlags(device.pointerTypeInfo, POINTER_FLAG_SECONDBUTTON);
	// } else {
	//     unsetPointerFlags(device.pointerTypeInfo, POINTER_FLAG_SECONDBUTTON);
	// }

	w, _, _ := getForegroundWindow.Call()
	device.pointerTypeInfo.penInfo.pointerInfo.hwndTarget = HWND(w)
	return sendDigiData(*device)
}

func (device *PointerDevice) Destory() bool {
	ans, _, _ := destroySyntheticPointerDevice.Call(uintptr(device.device))
	return ans != 0
}

func setupDigiInfo(penInfo *POINTER_PEN_INFO) {
	penInfo.pointerInfo.pointerType = PT_PEN
	penInfo.pointerInfo.pointerId = 0
	penInfo.pointerInfo.frameId = 0
	penInfo.pointerInfo.pointerFlags = POINTER_FLAG_NONE
	penInfo.pointerInfo.sourceDevice = 0
	w, _, _ := getForegroundWindow.Call()
	penInfo.pointerInfo.hwndTarget = HWND(w)
	penInfo.pointerInfo.ptPixelLocation = POINT{0, 0}
	penInfo.pointerInfo.ptPixelLocationRaw = POINT{0, 0}
	penInfo.pointerInfo.dwTime = 0
	penInfo.pointerInfo.historyCount = 0
	penInfo.pointerInfo.dwKeyStates = 0
	penInfo.pointerInfo.PerformanceCount = 0
	penInfo.pointerInfo.ButtonChangeType = 0 // POINTER_CHANGE_NONE

	penInfo.penFlags = PEN_FLAG_NONE
	penInfo.penMask = PEN_MASK_PRESSURE
	penInfo.pressure = 512
	penInfo.rotation = 0
	penInfo.tiltX = 0
	penInfo.tiltY = 0
}

func clearDigiFlags(pointerTypeInfo *POINTER_TYPE_INFO, flags PEN_FLAGS) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags = POINTER_FLAGS(flags)
}

func cleanDigiFlags(pointerTypeInfo *POINTER_TYPE_INFO) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags = 0
}

func setDigiFlags(pointerTypeInfo *POINTER_TYPE_INFO, flags PEN_FLAGS) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags |= POINTER_FLAGS(flags)
}
func unsetDigiFlags(pointerTypeInfo *POINTER_TYPE_INFO, flags PEN_FLAGS) {
	pointerTypeInfo.penInfo.pointerInfo.pointerFlags &= ^POINTER_FLAGS(flags)
}

func sendDigiData(device PointerDevice) error {
	ans, _, _ := injectSyntheticPointerInput.Call(uintptr(device.device), uintptr(unsafe.Pointer(&device.pointerTypeInfo)), 1)

	if ans == 0 {
		err, _, _ := getLastErr.Call()
		return fmt.Errorf("injectSyntheticPointerInput: failed(%d)", err)
	}
	return nil
}
