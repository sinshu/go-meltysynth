package meltysynth

type Synthesizer struct {
	SampleRate           int32
	BlockSize            int32
	SoundFont            *SoundFont
	Channels             []*channel
	minimumVoiceDuration int32
}
