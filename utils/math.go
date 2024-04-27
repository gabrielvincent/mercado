package utils

import (
	"math"
	"slices"
)

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func CalcAverage[T Numeric](numbers []T) float64 {
	var total T = 0.0
	for _, n := range numbers {
		total += n
	}
	return float64(total) / float64(len(numbers))
}

func CalcStandardDeviation[T Numeric](numbers []T) float64 {
	avg := CalcAverage(numbers)
	sumOfSquares := 0.0
	for _, number := range numbers {
		sumOfSquares += math.Pow(float64(number)-avg, 2)
	}
	variance := sumOfSquares / float64(len(numbers))
	return math.Sqrt(variance)
}

func CalcMedian[T Numeric](numbers []T) float64 {

	slices.SortFunc(numbers, func(a T, b T) int {
		switch true {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	})

	n := len(numbers)
	if n == 0 {
		return 0
	}
	middle := n / 2
	if n&1 == 1 {
		return float64(
			numbers[middle],
		)
	}
	return float64(
		numbers[middle-1]+numbers[middle],
	) / 2.0
}
