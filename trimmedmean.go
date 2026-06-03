// Package trimmedmean computes symmetric and asymmetric trimmed means.
//
// A trimmed mean discards a proportion of the smallest and largest values
// before averaging the rest, which makes it robust to outliers. The amount
// removed from each end is expressed as a proportion of the sample size and
// is rounded down (floored) to a whole number of observations, matching the
// behaviour of R's mean(x, trim = ...).
package trimmedmean

import (
	"errors"
	"math"
	"sort"
)

// Number constrains the element types the package can average over: any
// built-in integer or floating-point type, including named types based on
// them. This lets a single function serve both integer and float samples.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sentinel errors let callers test for a specific cause with errors.Is.
var (
	ErrEmptyData       = errors.New("trimmedmean: data slice is empty")
	ErrTooManyArgs     = errors.New("trimmedmean: expected at most two trimming proportions")
	ErrProportionRange = errors.New("trimmedmean: each proportion must be in [0, 0.5]")
	ErrTrimsEverything = errors.New("trimmedmean: trimming would remove every observation")
)

// Compute returns the trimmed mean of data.
//
// The trimming proportions give the share of observations to drop from each
// end after sorting from lowest to highest:
//
//   - no proportion         -> ordinary (untrimmed) mean
//   - one proportion p       -> symmetric trim: p from the low end and p from the high end
//   - two proportions lo, hi -> asymmetric trim: lo from the low end, hi from the high end
//
// The count removed at an end is floor(n * proportion), so the symmetric case
// matches R's mean(x, trim = p) exactly. The caller's slice is never modified.
func Compute[T Number](data []T, proportions ...float64) (float64, error) {
	if len(data) == 0 {
		return 0, ErrEmptyData
	}
	if len(proportions) > 2 {
		return 0, ErrTooManyArgs
	}

	lowProportion, highProportion := resolveProportions(proportions)
	if !inRange(lowProportion) || !inRange(highProportion) {
		return 0, ErrProportionRange
	}

	sorted := sortedFloatCopy(data)
	n := len(sorted)
	lowCount := int(math.Floor(float64(n) * lowProportion))
	highCount := int(math.Floor(float64(n) * highProportion))
	if lowCount+highCount >= n {
		return 0, ErrTrimsEverything
	}

	return average(sorted[lowCount : n-highCount]), nil
}

// resolveProportions turns the variadic argument into explicit low and high
// proportions, applying the symmetric default when only one value is given.
func resolveProportions(proportions []float64) (low, high float64) {
	switch len(proportions) {
	case 0:
		return 0, 0
	case 1:
		return proportions[0], proportions[0]
	default:
		return proportions[0], proportions[1]
	}
}

// inRange reports whether a single trimming proportion is usable. R rejects
// anything outside [0, 0.5], and so do we.
func inRange(p float64) bool {
	return p >= 0 && p <= 0.5
}

// sortedFloatCopy converts every element to float64 and returns the values
// sorted ascending, leaving the caller's slice untouched.
func sortedFloatCopy[T Number](data []T) []float64 {
	sorted := make([]float64, len(data))
	for i, value := range data {
		sorted[i] = float64(value)
	}
	sort.Float64s(sorted)
	return sorted
}

// average returns the arithmetic mean of a non-empty slice.
func average(values []float64) float64 {
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}
