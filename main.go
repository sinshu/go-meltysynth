package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/sinshu/go-meltysynth/meltysynth"
)

func main() {

	sf2, _ := os.Open("TimGM6mb.sf2")
	soundFont, _ := meltysynth.NewSoundFont(sf2)
	sf2.Close()

	settings := meltysynth.NewSynthesizerSettings(44100)
	synthesizer, _ := meltysynth.NewSynthesizer(soundFont, settings)

	mid, _ := os.Open("C:\\Windows\\Media\\flourish.mid")
	midiFile, _ := meltysynth.NewMidiFile(mid)
	mid.Close()

	sequencer := meltysynth.NewMidiFileSequencer(synthesizer)
	sequencer.Play(midiFile, true)

	length := int(float64(settings.SampleRate) * float64(midiFile.GetLength()) / float64(time.Second))
	left := make([]float32, length)
	right := make([]float32, length)
	sequencer.Render(left, right)

	interleaved := make([]float32, 2*length)
	for i := 0; i < length; i++ {
		interleaved[2*i] = left[i]
		interleaved[2*i+1] = right[i]
	}

	pcm, _ := os.Create("out.pcm")
	binary.Write(pcm, binary.LittleEndian, interleaved)
	pcm.Close()

	fmt.Println("DONE!")
}
