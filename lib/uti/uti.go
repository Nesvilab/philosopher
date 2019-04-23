// Package uti provides general, low priority methods and functions for different purposes
package uti

import (
	"math"
)

// Round serves the rol of the missing math.Round function
func Round(val float64, roundOn float64, places int) (newVal float64) {

	var round float64

	pow := math.Pow(10, float64(places))
	digit := pow * val

	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)

	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	newVal = round / pow

	return
}

// ToFixed truncates float numbers to a specific position
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(toFixedRound(num*output)) / output
}

func toFixedRound(num float64) int {
	return int(num + math.Copysign(0.05, num))
}
