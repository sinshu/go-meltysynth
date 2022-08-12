package main

import (
	"fmt"
	"os"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

func main() {

	srcFile, err := os.Open("PROGROCK.MID")
	if err != nil {
		panic(err)
	}

	mf, err := meltysynth.NewMidiFile(srcFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(mf)

	/*
		srcFile, err := os.Open("timgm6mb.sf2")
		if err != nil {
			panic("OMG1")
		}

		soundFont, err := meltysynth.NewSoundFont(srcFile)
		if err != nil {
			panic("OMG2")
		}

		var sampleRate int32 = 44100

		settings := meltysynth.NewSynthesizerSettings(sampleRate)

		synthesizer, err := meltysynth.NewSynthesizer(soundFont, settings)
		if err != nil {
			panic("OMG3")
		}

		synthesizer.NoteOn(0, 60, 100)
		synthesizer.NoteOn(0, 64, 100)
		synthesizer.NoteOn(0, 67, 100)

		left := make([]float32, 3*sampleRate)
		right := make([]float32, 3*sampleRate)

		synthesizer.Render(left, right)

		buf := make([]float32, 2*len(left))
		for i := 0; i < len(left); i++ {
			buf[2*i] = left[i]
			buf[2*i+1] = right[i]
		}

		dstFile, err := os.Create("out.pcm")
		if err != nil {
			panic("OMG4")
		}

		binary.Write(dstFile, binary.LittleEndian, buf)

		fmt.Println("DONE!")
	*/
}
