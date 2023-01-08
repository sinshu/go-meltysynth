package meltysynth

import "math"

type voiceCollection struct {
	synthesizer      *Synthesizer
	voices           []*voice
	activeVoiceCount int32
}

func newVoiceCollection(s *Synthesizer, maxActiveVoiceCount int32) *voiceCollection {
	result := &voiceCollection{
		synthesizer: s,
		voices:      make([]*voice, maxActiveVoiceCount),
	}
	for i := 0; i < len(result.voices); i++ {
		result.voices[i] = newVoice(s)
	}
	result.activeVoiceCount = 0

	return result
}

func (vc *voiceCollection) requestNew(region *InstrumentRegion, channel int32) *voice {
	// If an exclusive class is assigned to the region, find a voice with the same class.
	// If found, reuse it to avoid playing multiple voices with the same class at a time.
	exclusiveClass := region.GetExclusiveClass()
	if exclusiveClass != 0 {
		for i := int32(0); i < vc.activeVoiceCount; i++ {
			voice := vc.voices[i]
			if voice.exclusiveClass == exclusiveClass && voice.channel == channel {
				return voice
			}
		}
	}

	// If the number of active voices is less than the limit, use a free one.
	if int(vc.activeVoiceCount) < len(vc.voices) {
		free := vc.voices[vc.activeVoiceCount]
		vc.activeVoiceCount++
		return free
	}

	// Too many active voices...
	// Find one which has the lowest priority.
	var candidate *voice = nil
	var lowestPriority float32 = math.MaxFloat32
	for i := int32(0); i < vc.activeVoiceCount; i++ {
		voice := vc.voices[i]
		priority := voice.getPriority()
		if priority < lowestPriority {
			lowestPriority = priority
			candidate = voice
		} else if priority == lowestPriority {
			// Same priority...
			// The older one should be more suitable for reuse.
			if voice.voiceLength > candidate.voiceLength {
				candidate = voice
			}
		}
	}
	return candidate
}

func (vc *voiceCollection) process() {
	var i int32

	for {
		if i == vc.activeVoiceCount {
			return
		}

		if vc.voices[i].process() {
			i++
		} else {
			vc.activeVoiceCount--

			tmp := vc.voices[i]
			vc.voices[i] = vc.voices[vc.activeVoiceCount]
			vc.voices[vc.activeVoiceCount] = tmp
		}
	}
}

func (vc *voiceCollection) clear() {
	vc.activeVoiceCount = 0
}
