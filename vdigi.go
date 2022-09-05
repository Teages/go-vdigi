package vdigi

type Pointer struct {
	d PointerDevice
}

func CreatePointer() Pointer {
	p := Pointer{}
	p.d.Create()
	return p
}

func (p Pointer) Update(x, y int32, pressure uint32) error {
	return p.d.Update(x, y, pressure)
}

func (p Pointer) Destory() {
	p.d.Destory()
}
