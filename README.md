# go-vdigi

Create a virtual tablet input with golang.

Windows only for now, need help in support other platform.

## Example 

```golang
package main

import (
	"github.com/Teages/go-vdigi"
)

func main() {
	// Create a vdigi devices
	d := vdigi.CreatePointer()

	// Update the position and pressure
	err := d.Update(100, 100, 0)
  // Check error
	if err != nil {
		println(err.Error())
	}

	// Destroy the vdigi, or will close with the program
	d.Destory()
}
```

## TODO 

Linux & max support, but I don't have mac