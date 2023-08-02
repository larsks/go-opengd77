package main

import (
	"log"
	"os"

	"opengd77/pkg/opengd77"
)

func main() {
	radio, err := opengd77.NewRadio("/dev/ttyACM0")
	if err != nil {
		panic(err)
	}
	defer radio.Close()

	radio.ScreenClear()
	radio.ScreenShow()
	radio.ScreenPrint(0, 16, 3, 1, 0, "Reading codeplug")
	radio.ScreenRender()

	fd, err := os.Create("data.bin")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for _, block := range opengd77.ConfigBlocks {
		log.Printf("reading %d bytes from %x", block.Length, block.RadioOffset)
		buf := make([]byte, block.Length)
		if err := radio.ReadMemory(block.Kind, block.RadioOffset, block.Length, buf); err != nil {
			panic(err)
		}

		if _, err := fd.Seek(int64(block.FileOffset), 0); err != nil {
			panic(err)
		}

		if _, err := fd.Write(buf); err != nil {
			panic(err)
		}
	}

	radio.ScreenClose()

	log.Printf("all done")
}
