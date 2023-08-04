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

	iter := cp.ZoneIter()
	for {
		_, zone, ok := iter()
		if !ok {
			break
		}

		if zone.Name[0] == 0xff || zone.Name[0] == 0 {
			break
		}

		fmt.Printf("%s\n",
			zone.Name,
		)

		for _, chnum := range zone.Channels {
			if chnum == 0 {
				continue
			}

			ch, err := cp.Channel(int(chnum))
			if err != nil {
				panic(err)
			}

			fmt.Printf("  [%03d]: %s\n", chnum, ch)
		}
	}
}
