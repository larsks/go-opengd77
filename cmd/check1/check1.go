package main

import (
	"encoding/binary"
	"fmt"
	"opengd77/pkg/opengd77"
	"os"
	"strings"
)

func main() {
	fd, err := os.Open("data.bin")
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	for _, addr := range opengd77.ChannelBlocks {
		fd.Seek(int64(addr)+16, 0)
		for i := 0; i < 128; i++ {
			ch := opengd77.NewChannel()
			if err := binary.Read(fd, binary.LittleEndian, ch); err != nil {
				panic(err)
			}

			if ch.Name[0] == 0xff {
				continue
			}

			if strings.Contains(string(ch.Name[:]), "MMRA") {
				fmt.Println(ch.Name)
				fmt.Printf("rx %f [t %s] tx %f [t %s]\n",
					ch.GetRxFreq(),
					ch.RxTone,
					ch.GetTxFreq(),
					ch.TxTone,
				)
			}
		}
	}
}
