package meltysynth

import "math"

type chorus struct {
	bufferL []float32
	bufferR []float32

	delayTable []float32

	bufferIndexL int
	bufferIndexR int

	delayTableIndexL int
	delayTableIndexR int
}

func newChorus(sampleRate int32, delay float64, depth float64, frequency float64) *chorus {
	c := &chorus{}

	c.bufferL = make([]float32, int(float64(sampleRate)*(delay+depth))+2)
	c.bufferR = make([]float32, int(float64(sampleRate)*(delay+depth))+2)

	c.delayTable = make([]float32, int(math.Round(float64(sampleRate)/frequency)))
	delayTableLength := len(c.delayTable)
	for t := 0; t < delayTableLength; t++ {
		phase := 2 * math.Pi * float64(t) / float64(delayTableLength)
		c.delayTable[t] = float32(float64(sampleRate) * (delay + depth*math.Sin(phase)))
	}

	c.bufferIndexL = 0
	c.bufferIndexR = 0

	c.delayTableIndexL = 0
	c.delayTableIndexR = delayTableLength / 4

	return c
}

func (c *chorus) process(inputLeft []float32, intputRight []float32, outputLeft []float32, outputRight []float32) {
	bufferLength := len(c.bufferL)
	delayTableLength := len(c.delayTable)
	inputLength := len(inputLeft)

	for t := 0; t < inputLength; t++ {
		position := float64(c.bufferIndexL) - float64(c.delayTable[c.delayTableIndexL])
		if position < 0.0 {
			position += float64(bufferLength)
		}

		index1 := int(position)
		index2 := index1 + 1

		if index2 == bufferLength {
			index2 = 0
		}

		x1 := float64(c.bufferL[index1])
		x2 := float64(c.bufferL[index2])
		a := position - float64(index1)
		outputLeft[t] = float32(x1 + a*(x2-x1))

		c.bufferL[c.bufferIndexL] = inputLeft[t]
		c.bufferIndexL++
		if c.bufferIndexL == bufferLength {
			c.bufferIndexL = 0
		}

		c.delayTableIndexL++
		if c.delayTableIndexL == delayTableLength {
			c.delayTableIndexL = 0
		}
	}

	for t := 0; t < inputLength; t++ {
		position := float64(c.bufferIndexR) - float64(c.delayTable[c.delayTableIndexR])
		if position < 0.0 {
			position += float64(bufferLength)
		}

		index1 := int(position)
		index2 := index1 + 1

		if index2 == bufferLength {
			index2 = 0
		}

		x1 := float64(c.bufferR[index1])
		x2 := float64(c.bufferR[index2])
		a := position - float64(index1)
		outputRight[t] = float32(x1 + a*(x2-x1))

		c.bufferR[c.bufferIndexR] = intputRight[t]
		c.bufferIndexR++
		if c.bufferIndexR == bufferLength {
			c.bufferIndexR = 0
		}

		c.delayTableIndexR++
		if c.delayTableIndexR == delayTableLength {
			c.delayTableIndexR = 0
		}
	}
}

func (c *chorus) mute() {
	bufferLength := len(c.bufferL)
	for t := 0; t < bufferLength; t++ {
		c.bufferL[t] = 0
	}
	for t := 0; t < bufferLength; t++ {
		c.bufferR[t] = 0
	}
}
