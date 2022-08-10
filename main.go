package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

func main() {

	fs, err := os.Open("timgm6mb.sf2")

	if err != nil {
		panic("OMG")
	}

	sf2, _ := meltysynth.NewSoundFont(fs)

	fmt.Println(sf2.Info.Auther)

	file, err2 := os.Create("test.pcm")
	if err2 != nil {
		panic(err2.Error())
	}
	binary.Write(file, binary.LittleEndian, sf2.WaveData)
	file.Close()

	presetCount := len(sf2.Presets)
	for i := 0; i < presetCount; i++ {
		pre := sf2.Presets[i]
		fmt.Println(pre.Name)

		regCount := len(pre.Regions)
		for r := 0; r < regCount; r++ {
			reg := pre.Regions[r]
			fmt.Println("    " + reg.Instrument.Name)
		}
	}
}
