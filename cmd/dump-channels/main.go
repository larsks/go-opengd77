package main

import (
	"fmt"
	"opengd77/pkg/opengd77"
	"os"
)

func main() {
	fd, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	cp, err := opengd77.NewCodeplug(fd)
	if err != nil {
		panic(err)
	}

	iter := cp.ChannelIter()
	for {
		_, ch, ok := iter()
		if !ok {
			break
		}

		fmt.Printf("%s\n", ch)
		fmt.Printf("  libredmrflag1: %+v\n", ch.LibreDMRFlag1.Unpack())
		fmt.Printf("  flag4: %+v\n", ch.Flag4.Unpack())
	}
}
