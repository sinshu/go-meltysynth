package meltysynth

const (
	voice_Playing          int32 = 0
	voice_ReleaseRequested int32 = 1
	voice_Released         int32 = 2
)

type voice struct {
	synthesizer *Synthesizer

	volEnv *volumeEnvelope
	modEnv *modulationEnvelope

	vibLfo *lfo
	modLfo *lfo

	oscillator *oscillator
	filter     *biQuadFilter

	block []float32

	// A sudden change in the mix gain will cause pop noise.
	// To avoid this, we save the mix gain of the previous block,
	// and smooth out the gain if the gap between the current and previous gain is too large.
	// The actual smoothing process is done in the WriteBlock method of the Synthesizer class.

	previousMixGainLeft  float32
	previousMixGainRight float32
	currentMixGainLeft   float32
	currentMixGainRight  float32

	previousReverbSend float32
	previousChorusSend float32
	currentReverbSend  float32
	currentChorusSend  float32

	exclusiveClass int32
	channel        int32
	key            int32
	velocity       int32

	noteGain float32

	cutoff    float32
	resonance float32

	vibLfoToPitch float32
	modLfoToPitch float32
	modEnvToPitch float32

	modLfoToCutoff int32
	modEnvToCutoff int32
	dynamicCutoff  bool

	modLfoToVolume float32
	dynamicVolume  bool

	instrumentPan    float32
	instrumentReverb float32
	instrumentChorus float32

	// Some instruments require fast cutoff change, which can cause pop noise.
	// This is used to smooth out the cutoff frequency.
	smoothedCutoff float32

	voiceState  int32
	voiceLength int32
}

func newVoice(synthesizer *Synthesizer) *voice {

	result := new(voice)

	result.synthesizer = synthesizer

	result.volEnv = newVolumeEnvelope(synthesizer)
	result.modEnv = newModulationEnvelope(synthesizer)

	result.vibLfo = newLfo(synthesizer)
	result.modLfo = newLfo(synthesizer)

	result.oscillator = newOscillator(synthesizer)
	result.filter = newBiQuadFilter(synthesizer)

	result.block = make([]float32, synthesizer.BlockSize)

	return result
}

func (voice *voice) start(region regionPair, channel int32, key int32, velocity int32) {

	voice.exclusiveClass = region.GetExclusiveClass()
	voice.channel = channel
	voice.key = key
	voice.velocity = velocity

	if velocity > 0 {
		// According to the Polyphone's implementation, the initial attenuation should be reduced to 40%.
		// I'm not sure why, but this indeed improves the loudness variability.
		sampleAttenuation := 0.4 * region.GetInitialAttenuation()
		filterAttenuation := 0.5 * region.GetInitialFilterQ()
		decibels := 2*calcLinearToDecibels(float32(velocity)/float32(127)) - sampleAttenuation - filterAttenuation
		voice.noteGain = calcDecibelsToLinear(decibels)
	} else {
		voice.noteGain = 0
	}

	voice.cutoff = region.GetInitialFilterCutoffFrequency()
	voice.resonance = calcDecibelsToLinear(region.GetInitialFilterQ())

	voice.vibLfoToPitch = 0.01 * float32(region.GetVibratoLfoToPitch())
	voice.modLfoToPitch = 0.01 * float32(region.GetModulationLfoToPitch())
	voice.modEnvToPitch = 0.01 * float32(region.GetModulationEnvelopeToPitch())

	voice.modLfoToCutoff = region.GetModulationLfoToFilterCutoffFrequency()
	voice.modEnvToCutoff = region.GetModulationEnvelopeToFilterCutoffFrequency()
	voice.dynamicCutoff = voice.modLfoToCutoff != 0 || voice.modEnvToCutoff != 0

	voice.modLfoToVolume = region.GetModulationLfoToVolume()
	voice.dynamicVolume = voice.modLfoToVolume > 0.05

	voice.instrumentPan = calcClamp(region.GetPan(), -50, 50)
	voice.instrumentReverb = 0.01 * region.GetReverbEffectsSend()
	voice.instrumentChorus = 0.01 * region.GetChorusEffectsSend()

	voice.volEnv.startByRegion(region, key, velocity)
	voice.modEnv.startByRegion(region, key, velocity)
	voice.vibLfo.startVibrato(region, key, velocity)
	voice.modLfo.startModulation(region, key, velocity)
	voice.oscillator.startByRegion(voice.synthesizer.SoundFont.WaveData, region)
	voice.filter.clearBuffer()
	voice.filter.setLowPassFilter(voice.cutoff, voice.resonance)

	voice.smoothedCutoff = voice.cutoff

	voice.voiceState = voice_Playing
	voice.voiceLength = 0
}
