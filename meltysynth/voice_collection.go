package meltysynth

import "math"

type voiceCollection struct {
	synthesizer      *Synthesizer
	voices           []*voice
	activeVoiceCount int32
}

func newVoiceCollection(synthesizer *Synthesizer, maxActiveVoiceCount int32) *voiceCollection {

	result := new(voiceCollection)

	result.synthesizer = synthesizer

	result.voices = make([]*voice, maxActiveVoiceCount)
	for i := 0; i < len(result.voices); i++ {
		result.voices[i] = newVoice(synthesizer)
	}

	result.activeVoiceCount = 0

	return result
}

func (collection *voiceCollection) requestNew(region *InstrumentRegion, channel int32) *voice {

	// If an exclusive class is assigned to the region, find a voice with the same class.
	// If found, reuse it to avoid playing multiple voices with the same class at a time.
	exclusiveClass := region.GetExclusiveClass()
	if exclusiveClass != 0 {
		for i := int32(0); i < collection.activeVoiceCount; i++ {
			voice := collection.voices[i]
			if voice.exclusiveClass == exclusiveClass && voice.channel == channel {
				return voice
			}
		}
	}

	// If the number of active voices is less than the limit, use a free one.
	if int(collection.activeVoiceCount) < len(collection.voices) {
		free := collection.voices[collection.activeVoiceCount]
		collection.activeVoiceCount++
		return free
	}

	// Too many active voices...
	// Find one which has the lowest priority.
	var candidate *voice = nil
	var lowestPriority float32 = math.MaxFloat32
	for i := int32(0); i < collection.activeVoiceCount; i++ {
		voice := collection.voices[i]
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

func (collection *voiceCollection) process() {

	var i int32 = 0

	for {
		if i == collection.activeVoiceCount {
			return
		}

		if collection.voices[i].process() {
			i++
		} else {
			collection.activeVoiceCount--

			tmp := collection.voices[i]
			collection.voices[i] = collection.voices[collection.activeVoiceCount]
			collection.voices[collection.activeVoiceCount] = tmp
		}
	}
}

func (collection *voiceCollection) clear() {
	collection.activeVoiceCount = 0
}
