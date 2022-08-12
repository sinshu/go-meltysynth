package meltysynth

import (
	"math"
)

const (
	synth_channelCount      int32 = 16
	synth_percussionChannel int32 = 9
)

type Synthesizer struct {
	SoundFont             *SoundFont
	SampleRate            int32
	BlockSize             int32
	MaximumPolyphony      int32
	EnableReverbAndChorus bool

	minimumVoiceDuration int32

	presetLookup  map[int32]*Preset
	defaultPreset *Preset

	channels []*channel

	voices *voiceCollection

	blockLeft  []float32
	blockRight []float32

	inverseBlockSize float32

	blockRead int32

	MasterVolume float32
}

func NewSynthesizer(soundFont *SoundFont, settings *SynthesizerSettings) (*Synthesizer, error) {

	err := settings.validate()
	if err != nil {
		return nil, err
	}

	result := new(Synthesizer)

	result.SoundFont = soundFont
	result.SampleRate = settings.SampleRate
	result.BlockSize = settings.BlockSize
	result.MaximumPolyphony = settings.MaximumPolyphony
	result.EnableReverbAndChorus = settings.EnableReverbAndChorus

	result.minimumVoiceDuration = settings.SampleRate / 500

	result.presetLookup = make(map[int32]*Preset)

	minPresetId := int32(math.MaxInt32)
	for i := 0; i < len(soundFont.Presets); i++ {
		preset := soundFont.Presets[i]
		// The preset ID is Int32, where the upper 16 bits represent the bank number
		// and the lower 16 bits represent the patch number.
		// This ID is used to search for presets by the combination of bank number
		// and patch number.
		presetId := (preset.BankNumber << 16) | preset.PatchNumber
		result.presetLookup[presetId] = preset

		// The preset with the minimum ID number will be default.
		// If the SoundFont is GM compatible, the piano will be chosen.
		if presetId < minPresetId {
			result.defaultPreset = preset
			minPresetId = presetId
		}
	}

	result.channels = make([]*channel, synth_channelCount)
	for i := int32(0); int(i) < len(result.channels); i++ {
		result.channels[i] = newChannel(result, i == synth_percussionChannel)
	}

	result.voices = newVoiceCollection(result, result.MaximumPolyphony)

	result.blockLeft = make([]float32, result.BlockSize)
	result.blockRight = make([]float32, result.BlockSize)

	result.inverseBlockSize = 1 / float32(result.BlockSize)

	result.blockRead = result.BlockSize

	result.MasterVolume = 0.5

	return result, nil
}

func (synthesizer *Synthesizer) ProcessMidiMessage(channel int32, command int32, data1 int32, data2 int32) {

	if !(0 <= channel && int(channel) < len(synthesizer.channels)) {
		return
	}

	channelInfo := synthesizer.channels[channel]

	switch command {
	case 0x80: // Note Off
		synthesizer.NoteOff(channel, data1)

	case 0x90: // Note On
		synthesizer.NoteOn(channel, data1, data2)

	case 0xB0: // Controller
		switch data1 {
		case 0x00: // Bank Selection
			channelInfo.setBank(data2)

		case 0x01: // Modulation Coarse
			channelInfo.setModulationCoarse(data2)

		case 0x21: // Modulation Fine
			channelInfo.setModulationFine(data2)

		case 0x06: // Data Entry Coarse
			channelInfo.dataEntryCoarse(data2)

		case 0x26: // Data Entry Fine
			channelInfo.dataEntryFine(data2)

		case 0x07: // Channel Volume Coarse
			channelInfo.setVolumeCoarse(data2)

		case 0x27: // Channel Volume Fine
			channelInfo.setVolumeFine(data2)

		case 0x0A: // Pan Coarse
			channelInfo.setPanCoarse(data2)

		case 0x2A: // Pan Fine
			channelInfo.setPanFine(data2)

		case 0x0B: // Expression Coarse
			channelInfo.setExpressionCoarse(data2)

		case 0x2B: // Expression Fine
			channelInfo.setExpressionFine(data2)

		case 0x40: // Hold Pedal
			channelInfo.setHoldPedal(data2)

		case 0x5B: // Reverb Send
			channelInfo.setReverbSend(data2)

		case 0x5D: // Chorus Send
			channelInfo.setChorusSend(data2)

		case 0x65: // RPN Coarse
			channelInfo.setRpnCoarse(data2)

		case 0x64: // RPN Fine
			channelInfo.setRpnFine(data2)

		case 0x78: // All Sound Off
			synthesizer.NoteOffAllChannel(channel, true)

		case 0x79: // Reset All Controllers
			synthesizer.ResetAllControllersChannel(channel)

		case 0x7B: // All Note Off
			synthesizer.NoteOffAllChannel(channel, false)
		}

	case 0xC0: // Program Change
		channelInfo.setPatch(data1)

	case 0xE0: // Pitch Bend
		channelInfo.setPitchBend(data1, data2)
	}
}

func (synthesizer *Synthesizer) NoteOff(channel int32, key int32) {

	if !(0 <= channel && int(channel) < len(synthesizer.channels)) {
		return
	}

	for i := int32(0); i < synthesizer.voices.activeVoiceCount; i++ {
		voice := synthesizer.voices.voices[i]
		if voice.channel == channel && voice.key == key {
			voice.end()
		}
	}
}

func (synthesizer *Synthesizer) NoteOn(channel int32, key int32, velocity int32) {

	if velocity == 0 {
		synthesizer.NoteOff(channel, key)
		return
	}

	if !(0 <= channel && int(channel) < len(synthesizer.channels)) {
		return
	}

	channelInfo := synthesizer.channels[channel]

	presetId := (channelInfo.bankNumber << 16) | channelInfo.patchNumber

	preset, found := synthesizer.presetLookup[presetId]
	if !found {
		// Try fallback to the GM sound set.
		// Normally, the given patch number + the bank number 0 will work.
		// For drums (bank number >= 128), it seems to be better to select the standard set (128:0).
		var gmPresetId int32
		if channelInfo.bankNumber < 128 {
			gmPresetId = channelInfo.patchNumber
		} else {
			gmPresetId = 128 << 16
		}

		preset, found = synthesizer.presetLookup[gmPresetId]
		if !found {
			// No corresponding preset was found. Use the default one...
			preset = synthesizer.defaultPreset
		}
	}

	presetCount := len(preset.Regions)
	for i := 0; i < presetCount; i++ {
		presetRegion := preset.Regions[i]
		if presetRegion.contains(key, velocity) {
			instrumentCount := len(presetRegion.Instrument.Regions)
			for j := 0; j < instrumentCount; j++ {
				instrumentRegion := presetRegion.Instrument.Regions[j]
				if instrumentRegion.contains(key, velocity) {
					regionPair := newRegionPair(presetRegion, instrumentRegion)

					voice := synthesizer.voices.requestNew(instrumentRegion, channel)
					if voice != nil {
						voice.start(regionPair, channel, key, velocity)
					}
				}
			}
		}
	}
}

func (synthesizer *Synthesizer) NoteOffAll(immediate bool) {

	if immediate {
		synthesizer.voices.clear()
	} else {
		for i := 0; i < int(synthesizer.voices.activeVoiceCount); i++ {
			synthesizer.voices.voices[i].end()
		}
	}
}

func (synthesizer *Synthesizer) NoteOffAllChannel(channel int32, immediate bool) {

	if immediate {
		for i := 0; i < int(synthesizer.voices.activeVoiceCount); i++ {
			if synthesizer.voices.voices[i].channel == channel {
				synthesizer.voices.voices[i].kill()
			}
		}
	} else {
		for i := 0; i < int(synthesizer.voices.activeVoiceCount); i++ {
			if synthesizer.voices.voices[i].channel == channel {
				synthesizer.voices.voices[i].end()
			}
		}
	}
}

func (synthesizer *Synthesizer) ResetAllControllers() {

	channelCount := len(synthesizer.channels)
	for i := 0; i < channelCount; i++ {
		synthesizer.channels[i].resetAllControllers()
	}
}

func (synthesizer *Synthesizer) ResetAllControllersChannel(channel int32) {

	if !(0 <= channel && int(channel) < len(synthesizer.channels)) {
		return
	}

	synthesizer.channels[channel].resetAllControllers()
}

func (synthesizer *Synthesizer) Reset() {

	synthesizer.voices.clear()

	channelCount := len(synthesizer.channels)
	for i := 0; i < channelCount; i++ {
		synthesizer.channels[i].reset()
	}

	synthesizer.blockRead = synthesizer.BlockSize
}

func (synthesizer *Synthesizer) Render(left []float32, right []float32) {

	wrote := int32(0)
	length := int32(len(left))
	for wrote < length {
		if synthesizer.blockRead == synthesizer.BlockSize {
			synthesizer.renderBlock()
			synthesizer.blockRead = 0
		}

		srcRem := synthesizer.BlockSize - synthesizer.blockRead
		dstRem := int32(length - wrote)
		rem := int32(math.Min(float64(srcRem), float64(dstRem)))

		for i := int32(0); i < rem; i++ {
			left[wrote+i] = synthesizer.blockLeft[synthesizer.blockRead+i]
			right[wrote+i] = synthesizer.blockRight[synthesizer.blockRead+i]
		}

		synthesizer.blockRead += rem
		wrote += rem
	}
}

func (synthesizer *Synthesizer) renderBlock() {

	synthesizer.voices.process()

	for i := 0; i < int(synthesizer.BlockSize); i++ {
		synthesizer.blockLeft[i] = 0
		synthesizer.blockRight[i] = 0
	}

	for i := 0; i < int(synthesizer.voices.activeVoiceCount); i++ {
		voice := synthesizer.voices.voices[i]
		previousGainLeft := synthesizer.MasterVolume * voice.previousMixGainLeft
		currentGainLeft := synthesizer.MasterVolume * voice.currentMixGainLeft
		synthesizer.writeBlock(previousGainLeft, currentGainLeft, voice.block, synthesizer.blockLeft)
		var previousGainRight = synthesizer.MasterVolume * voice.previousMixGainRight
		var currentGainRight = synthesizer.MasterVolume * voice.currentMixGainRight
		synthesizer.writeBlock(previousGainRight, currentGainRight, voice.block, synthesizer.blockRight)
	}
}

func (synthesizer *Synthesizer) writeBlock(previousGain float32, currentGain float32, source []float32, destination []float32) {

	if math.Max(float64(previousGain), float64(currentGain)) < float64(nonAudible) {
		return
	}

	if math.Abs(float64(currentGain-previousGain)) < 1.0e-3 {
		arrayMultiplyAdd(currentGain, source, destination)
	} else {
		step := synthesizer.inverseBlockSize * (currentGain - previousGain)
		arrayMultiplyAddSlope(previousGain, step, source, destination)
	}
}
