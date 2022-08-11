package meltysynth

import (
	"encoding/binary"
	"fmt"
	"os"
)

func Sandbox() {

	srcFile, err := os.Open("timgm6mb.sf2")
	if err != nil {
		panic("OMG1")
	}

	soundFont, err := NewSoundFont(srcFile)
	if err != nil {
		panic("OMG2")
	}

	var preset *Preset
	for i := 0; i < len(soundFont.Presets); i++ {
		preset = soundFont.Presets[i]
		if preset.Name == "Overdrive Guitar" {
			break
		}
	}

	synthesizer := new(Synthesizer)
	synthesizer.SampleRate = 44100

	oscillator := newOscillator(synthesizer)

	var key int32 = 60
	var velocity int32 = 100

	var presetRegion *PresetRegion
	var instrumentRegion *InstrumentRegion
LOOP:
	for i := 0; i < len(preset.Regions); i++ {
		presetRegion = preset.Regions[i]
		if presetRegion.contains(key, velocity) {
			for j := 0; j < len(presetRegion.Instrument.Regions); j++ {
				instrumentRegion = presetRegion.Instrument.Regions[j]
				if instrumentRegion.contains(key, velocity) {
					break LOOP
				}
			}
		}
	}

	fmt.Println(presetRegion.Instrument.Name)
	fmt.Println(instrumentRegion.Sample.Name)

	pair := newRegionPair(presetRegion, instrumentRegion)

	oscillator.startByRegion(soundFont.WaveData, pair)

	buffer := make([]float32, 441000)
	oscillator.process(buffer, float32(key))

	dstFile, err := os.Create("out.pcm")
	if err != nil {
		panic("OMG3")
	}

	binary.Write(dstFile, binary.LittleEndian, buffer)
}
