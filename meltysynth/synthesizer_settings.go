package meltysynth

import "errors"

const (
	synth_DefaultBlockSize             int32 = 64
	synth_DefaultMaximumPolyphony      int32 = 64
	synth_DefaultEnableReverbAndChorus bool  = true
)

type SynthesizerSettings struct {
	SampleRate            int32
	BlockSize             int32
	MaximumPolyphony      int32
	EnableReverbAndChorus bool
}

func NewSynthesizerSettings(sampleRate int32) *SynthesizerSettings {
	result := new(SynthesizerSettings)

	result.SampleRate = sampleRate
	result.BlockSize = synth_DefaultBlockSize
	result.MaximumPolyphony = synth_DefaultMaximumPolyphony
	result.EnableReverbAndChorus = synth_DefaultEnableReverbAndChorus

	return result
}

func (settings *SynthesizerSettings) validate() error {
	if !(16000 <= settings.SampleRate && settings.SampleRate <= 192000) {
		return errors.New("the sample rate must be between 16000 and 192000")
	}

	if !(8 <= settings.BlockSize && settings.BlockSize <= 1024) {
		return errors.New("the block size must be between 8 and 1024")
	}

	if !(8 <= settings.MaximumPolyphony && settings.MaximumPolyphony <= 256) {
		return errors.New("the maximum number of polyphony must be between 8 and 256")
	}

	return nil
}
