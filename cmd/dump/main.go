package main

import (
	"flag"
	"log"
	"os"

	"opengd77/pkg/opengd77"
)

var (
	serialDevice string
	outputFile   string
)

func init() {
	flag.StringVar(&serialDevice, "port", "/dev/ttyACM0", "set serial port")
	flag.StringVar(&outputFile, "output", "data.bin", "set output filename")
}

func main() {
	flag.Parse()

	radio, err := opengd77.NewRadio(serialDevice)
	if err != nil {
		panic(err)
	}
	if err := radio.Open(); err != nil {
		panic(err)
	}
	defer radio.Close()

	fd, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	if err := radio.ShowMessage(3, "Reading codeplug"); err != nil {
		panic(err)
	}

	log.Printf("reading codeplug")
	if err := radio.ReadCodePlug(fd); err != nil {
		panic(err)
	}

	radio.ScreenClose()
	log.Printf("all done")
}
