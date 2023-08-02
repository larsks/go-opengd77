package main

import (
	"fmt"
	"opengd77/pkg/opengd77"
	"unsafe"
)

func main() {
	fmt.Printf("size: %d\n", unsafe.Sizeof(opengd77.Channel{}))
}
