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

	radio.ReadCodePlug(fd)

	radio.ScreenClose()

	log.Printf("all done")
}
