package meltysynth

import (
	"errors"
	"io"
)

// Since modulators will not be supported, we discard the data.
func discardModulatorData(reader io.Reader, size int32) error {

	if size%10 != 0 {
		return errors.New("The modulator list is invalid.")
	}

	n, err := reader.Read(make([]byte, size, size))
	if err != nil {
		return err
	}
	if n != int(size) {
		return errors.New("Failed to read the modulator list.")
	}

	return nil
}
