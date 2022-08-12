package meltysynth

import "math"

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

func (voice *voice) end() {
	if voice.voiceState == voice_Playing {
		voice.voiceState = voice_ReleaseRequested
	}
}

func (voice *voice) kill() {
	voice.noteGain = 0
}

func (voice *voice) process() bool {

	if voice.noteGain < nonAudible {
		return false
	}

	channelInfo := voice.synthesizer.channels[voice.channel]

	voice.releaseIfNecessary(channelInfo)

	if !voice.volEnv.process(voice.synthesizer.BlockSize) {
		return false
	}

	voice.modEnv.process(voice.synthesizer.BlockSize)
	voice.vibLfo.process()
	voice.modLfo.process()

	vibPitchChange := (0.01*channelInfo.getModulation() + voice.vibLfoToPitch) * voice.vibLfo.value
	modPitchChange := voice.modLfoToPitch*voice.modLfo.value + voice.modEnvToPitch*voice.modEnv.value
	channelPitchChange := channelInfo.getTune() + channelInfo.getPitchBend()
	pitch := float32(voice.key) + vibPitchChange + modPitchChange + channelPitchChange
	if !voice.oscillator.process(voice.block, pitch) {
		return false
	}

	if voice.dynamicCutoff {
		cents := float32(voice.modLfoToCutoff)*voice.modLfo.value + float32(voice.modEnvToCutoff)*voice.modEnv.value
		factor := calcCentsToMultiplyingFactor(cents)
		newCutoff := factor * voice.cutoff

		// The cutoff change is limited within x0.5 and x2 to reduce pop noise.
		lowerLimit := 0.5 * voice.smoothedCutoff
		upperLimit := 2 * voice.smoothedCutoff
		if newCutoff < lowerLimit {
			voice.smoothedCutoff = lowerLimit
		} else if newCutoff > upperLimit {
			voice.smoothedCutoff = upperLimit
		} else {
			voice.smoothedCutoff = newCutoff
		}

		voice.filter.setLowPassFilter(voice.smoothedCutoff, voice.resonance)
	}
	voice.filter.process(voice.block)

	voice.previousMixGainLeft = voice.currentMixGainLeft
	voice.previousMixGainRight = voice.currentMixGainRight
	voice.previousReverbSend = voice.currentReverbSend
	voice.previousChorusSend = voice.currentChorusSend

	// According to the GM spec, the following value should be squared.
	ve := channelInfo.getVolume() * channelInfo.getExpression()
	channelGain := ve * ve

	mixGain := voice.noteGain * channelGain * voice.volEnv.value
	if voice.dynamicVolume {
		decibels := voice.modLfoToVolume * voice.modLfo.value
		mixGain *= calcDecibelsToLinear(decibels)
	}

	angle := float32(math.Pi/200) * (channelInfo.getPan() + voice.instrumentPan + 50)
	if angle <= 0 {
		voice.currentMixGainLeft = mixGain
		voice.currentMixGainRight = 0
	} else if angle >= halfPi {
		voice.currentMixGainLeft = 0
		voice.currentMixGainRight = mixGain
	} else {
		voice.currentMixGainLeft = mixGain * float32(math.Cos(float64(angle)))
		voice.currentMixGainRight = mixGain * float32(math.Sin(float64(angle)))
	}

	voice.currentReverbSend = calcClamp(channelInfo.getReverbSend()+voice.instrumentReverb, 0, 1)
	voice.currentChorusSend = calcClamp(channelInfo.getChorusSend()+voice.instrumentChorus, 0, 1)

	if voice.voiceLength == 0 {
		voice.previousMixGainLeft = voice.currentMixGainLeft
		voice.previousMixGainRight = voice.currentMixGainRight
		voice.previousReverbSend = voice.currentReverbSend
		voice.previousChorusSend = voice.currentChorusSend
	}

	voice.voiceLength += voice.synthesizer.BlockSize

	return true
}

func (voice *voice) releaseIfNecessary(channelInfo *channel) {

	if voice.voiceLength < voice.synthesizer.minimumVoiceDuration {
		return
	}

	if voice.voiceState == voice_ReleaseRequested && !channelInfo.holdPedal {
		voice.volEnv.release()
		voice.modEnv.release()
		voice.oscillator.release()

		voice.voiceState = voice_Released
	}
}

func (voice *voice) getPriority() float32 {

	if voice.noteGain < nonAudible {
		return 0
	} else {
		return voice.volEnv.priority
	}
}
