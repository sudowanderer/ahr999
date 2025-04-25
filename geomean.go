package main

import (
	"errors"
	"math"
)

// GeometricMean calculates the geometric mean of a list of positive float64 numbers.
// Returns an error if the list is empty or contains non-positive numbers.
func GeometricMean(data []float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("input slice is empty")
	}

	sumLog := 0.0
	c := 0.0 // Kahan compensation

	for _, v := range data {
		if v <= 0 {
			return 0, errors.New("input contains non-positive number")
		}
		y := math.Log(v) - c
		t := sumLog + y
		c = (t - sumLog) - y
		sumLog = t
	}

	meanLog := sumLog / float64(len(data))
	return math.Exp(meanLog), nil
}
