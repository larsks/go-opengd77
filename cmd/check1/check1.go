package main

import (
	"encoding/binary"
	"fmt"
	"opengd77/pkg/opengd77"
	"os"
	"strings"
)

func ternary(condition bool, valIfTrue, valIfFalse int) int {
	if condition {
		return valIfTrue
	} else {
		return valIfFalse
	}
}

func main() {
	fd, err := os.Open("data.bin")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for _, addr := range opengd77.ChannelBlocks {
		fd.Seek(int64(addr)+16, 0)
		for i := 0; i < 128; i++ {
			ch := opengd77.Channel{}
			if err := binary.Read(fd, binary.LittleEndian, &ch); err != nil {
				panic(err)
			}

			if ch.Name[0] == 0xff {
				continue
			}

			fmt.Println(ch.Name)
			if strings.Contains(string(ch.Name[:]), "MMRA") {
				fmt.Printf("rx %s [t %s] tx %s [t %s]\n",
					ch.RxFreq,
					ch.RxTone,
					ch.TxFreq,
					ch.TxTone,
				)
			}
		}
	}
}
