package main

import (
	"math"
	"testing"
)

func TestPercentile(t *testing.T) {
	acceptableError := 0.0000001
	values := []float64{3, 5, 6, 6, 7, 9, 10, 12, 13, 13, 15}

	tests := []struct {
		name     string
		p        int
		expected float64
	}{
		{"0th percentile", 0, 3},
		{"25th percentile", 25, 6},
		{"50th percentile", 50, 9},
		{"70th percentile", 70, 12.4},
		{"75th percentile", 75, 13},
		{"85th percentile", 85, 13.4},
		{"95th percentile", 95, 15},
		{"100th percentile", 100, 15},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			result := percentile(values, testCase.p)
			variance := math.Abs(result - testCase.expected)
			if variance > acceptableError {
				t.Errorf("Expected %s to be %f but got %f", testCase.name, testCase.expected, result)
			}
		})
	}
}
