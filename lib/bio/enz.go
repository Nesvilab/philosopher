package bio

import (
	"errors"
	"strings"
)

// Enzyme struct
type Enzyme struct {
	Name    string
	Pattern string
	Join    string
}

// Synth is an enzyme builder
func (e *Enzyme) Synth(t string) error {

	if strings.EqualFold(strings.ToLower(t), "trypsin") {
		e.Name = "trypsin"
		e.Pattern = "KR[^P]"
		e.Join = "KR"
	} else if strings.EqualFold(strings.ToLower(t), "lys_c") {
		e.Name = "lys_c"
		e.Pattern = "K[^P]"
		e.Join = "K"
	} else if strings.EqualFold(strings.ToLower(t), "lys_n") {
		e.Name = "lys_n"
		e.Pattern = "K"
		e.Join = "K"
	} else if strings.EqualFold(strings.ToLower(t), "chymotrypsin") {
		e.Name = "chymotrypsin"
		e.Pattern = "FWYL[^P]"
		e.Join = "K"
	} else {
		return errors.New("Enzyme not supported")
	}

	return nil
}
