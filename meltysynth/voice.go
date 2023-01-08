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

func newVoice(s *Synthesizer) *voice {
	return &voice{
		synthesizer: s,
		volEnv:      newVolumeEnvelope(s),
		modEnv:      newModulationEnvelope(s),
		vibLfo:      newLfo(s),
		modLfo:      newLfo(s),
		oscillator:  newOscillator(s),
		filter:      newBiQuadFilter(s),
		block:       make([]float32, s.BlockSize),
	}
}

func (v *voice) start(region regionPair, channel int32, key int32, velocity int32) {
	v.exclusiveClass = region.GetExclusiveClass()
	v.channel = channel
	v.key = key
	v.velocity = velocity

	if velocity > 0 {
		// According to the Polyphone's implementation, the initial attenuation should be reduced to 40%.
		// I'm not sure why, but this indeed improves the loudness variability.
		sampleAttenuation := 0.4 * region.GetInitialAttenuation()
		filterAttenuation := 0.5 * region.GetInitialFilterQ()
		decibels := 2*calcLinearToDecibels(float32(velocity)/float32(127)) - sampleAttenuation - filterAttenuation
		v.noteGain = calcDecibelsToLinear(decibels)
	} else {
		v.noteGain = 0
	}

	v.cutoff = region.GetInitialFilterCutoffFrequency()
	v.resonance = calcDecibelsToLinear(region.GetInitialFilterQ())

	v.vibLfoToPitch = 0.01 * float32(region.GetVibratoLfoToPitch())
	v.modLfoToPitch = 0.01 * float32(region.GetModulationLfoToPitch())
	v.modEnvToPitch = 0.01 * float32(region.GetModulationEnvelopeToPitch())

	v.modLfoToCutoff = region.GetModulationLfoToFilterCutoffFrequency()
	v.modEnvToCutoff = region.GetModulationEnvelopeToFilterCutoffFrequency()
	v.dynamicCutoff = v.modLfoToCutoff != 0 || v.modEnvToCutoff != 0

	v.modLfoToVolume = region.GetModulationLfoToVolume()
	v.dynamicVolume = v.modLfoToVolume > 0.05

	v.instrumentPan = calcClamp(region.GetPan(), -50, 50)
	v.instrumentReverb = 0.01 * region.GetReverbEffectsSend()
	v.instrumentChorus = 0.01 * region.GetChorusEffectsSend()

	v.volEnv.startByRegion(region, key, velocity)
	v.modEnv.startByRegion(region, key, velocity)
	v.vibLfo.startVibrato(region, key, velocity)
	v.modLfo.startModulation(region, key, velocity)
	v.oscillator.startByRegion(v.synthesizer.SoundFont.WaveData, region)
	v.filter.clearBuffer()
	v.filter.setLowPassFilter(v.cutoff, v.resonance)

	v.smoothedCutoff = v.cutoff

	v.voiceState = voice_Playing
	v.voiceLength = 0
}

func (v *voice) end() {
	if v.voiceState == voice_Playing {
		v.voiceState = voice_ReleaseRequested
	}
}

func (v *voice) kill() {
	v.noteGain = 0
}

func (v *voice) process() bool {
	if v.noteGain < nonAudible {
		return false
	}

	channelInfo := v.synthesizer.channels[v.channel]

	v.releaseIfNecessary(channelInfo)

	if !v.volEnv.process(v.synthesizer.BlockSize) {
		return false
	}

	v.modEnv.process(v.synthesizer.BlockSize)
	v.vibLfo.process()
	v.modLfo.process()

	vibPitchChange := (0.01*channelInfo.getModulation() + v.vibLfoToPitch) * v.vibLfo.value
	modPitchChange := v.modLfoToPitch*v.modLfo.value + v.modEnvToPitch*v.modEnv.value
	channelPitchChange := channelInfo.getTune() + channelInfo.getPitchBend()
	pitch := float32(v.key) + vibPitchChange + modPitchChange + channelPitchChange
	if !v.oscillator.process(v.block, pitch) {
		return false
	}

	if v.dynamicCutoff {
		cents := float32(v.modLfoToCutoff)*v.modLfo.value + float32(v.modEnvToCutoff)*v.modEnv.value
		factor := calcCentsToMultiplyingFactor(cents)
		newCutoff := factor * v.cutoff

		// The cutoff change is limited within x0.5 and x2 to reduce pop noise.
		lowerLimit := 0.5 * v.smoothedCutoff
		upperLimit := 2 * v.smoothedCutoff
		if newCutoff < lowerLimit {
			v.smoothedCutoff = lowerLimit
		} else if newCutoff > upperLimit {
			v.smoothedCutoff = upperLimit
		} else {
			v.smoothedCutoff = newCutoff
		}

		v.filter.setLowPassFilter(v.smoothedCutoff, v.resonance)
	}
	v.filter.process(v.block)

	v.previousMixGainLeft = v.currentMixGainLeft
	v.previousMixGainRight = v.currentMixGainRight
	v.previousReverbSend = v.currentReverbSend
	v.previousChorusSend = v.currentChorusSend

	// According to the GM spec, the following value should be squared.
	ve := channelInfo.getVolume() * channelInfo.getExpression()
	channelGain := ve * ve

	mixGain := v.noteGain * channelGain * v.volEnv.value
	if v.dynamicVolume {
		decibels := v.modLfoToVolume * v.modLfo.value
		mixGain *= calcDecibelsToLinear(decibels)
	}

	angle := float32(math.Pi/200) * (channelInfo.getPan() + v.instrumentPan + 50)
	switch {
	case angle <= 0:
		v.currentMixGainLeft = mixGain
		v.currentMixGainRight = 0
	case angle >= halfPi:
		v.currentMixGainLeft = 0
		v.currentMixGainRight = mixGain
	default:
		v.currentMixGainLeft = mixGain * float32(math.Cos(float64(angle)))
		v.currentMixGainRight = mixGain * float32(math.Sin(float64(angle)))
	}

	v.currentReverbSend = calcClamp(channelInfo.getReverbSend()+v.instrumentReverb, 0, 1)
	v.currentChorusSend = calcClamp(channelInfo.getChorusSend()+v.instrumentChorus, 0, 1)

	if v.voiceLength == 0 {
		v.previousMixGainLeft = v.currentMixGainLeft
		v.previousMixGainRight = v.currentMixGainRight
		v.previousReverbSend = v.currentReverbSend
		v.previousChorusSend = v.currentChorusSend
	}

	v.voiceLength += v.synthesizer.BlockSize

	return true
}

func (v *voice) releaseIfNecessary(channelInfo *channel) {
	if v.voiceLength < v.synthesizer.minimumVoiceDuration {
		return
	}

	if v.voiceState == voice_ReleaseRequested && !channelInfo.holdPedal {
		v.volEnv.release()
		v.modEnv.release()
		v.oscillator.release()

		v.voiceState = voice_Released
	}
}

func (v *voice) getPriority() float32 {
	if v.noteGain < nonAudible {
		return 0
	}
	return v.volEnv.priority
}
