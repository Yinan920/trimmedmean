package trimmedmean

import (
	"errors"
	"math"
	"testing"
)

const tolerance = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestUntrimmedMean(t *testing.T) {
	got, err := Compute([]int{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(got, 3.0) {
		t.Errorf("got %v, want 3.0", got)
	}
}

func TestSymmetricTrimRemovesOutlier(t *testing.T) {
	// A single huge outlier should not move a 10% symmetric trimmed mean.
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 1000}
	got, err := Compute(data, 0.1) // floor(10*0.1)=1 from each end -> mean of 2..9
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(got, 5.5) {
		t.Errorf("got %v, want 5.5", got)
	}
}

func TestAsymmetricTrim(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got, err := Compute(data, 0.2, 0.1) // drop 2 low, 1 high -> mean of 3..9
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := (3.0 + 4 + 5 + 6 + 7 + 8 + 9) / 7
	if !almostEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMatchesRFloorRule(t *testing.T) {
	// With n=100 and trim 0.05, R floors floor(100*0.05)=5 from each end,
	// leaving the values 6..95, whose mean is 50.5.
	data := make([]float64, 100)
	for i := range data {
		data[i] = float64(i + 1) // 1..100
	}
	got, err := Compute(data, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(got, 50.5) {
		t.Errorf("got %v, want 50.5", got)
	}
}

func TestDoesNotMutateInput(t *testing.T) {
	data := []int{5, 3, 1, 4, 2}
	if _, err := Compute(data, 0.2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, want := range []int{5, 3, 1, 4, 2} {
		if data[i] != want {
			t.Fatalf("input slice was mutated: %v", data)
		}
	}
}

func TestErrors(t *testing.T) {
	cases := []struct {
		name    string
		run     func() (float64, error)
		wantErr error
	}{
		{"empty", func() (float64, error) { return Compute([]int{}) }, ErrEmptyData},
		{"too many", func() (float64, error) { return Compute([]int{1, 2, 3}, 0.1, 0.2, 0.3) }, ErrTooManyArgs},
		{"out of range", func() (float64, error) { return Compute([]int{1, 2, 3}, 0.9) }, ErrProportionRange},
		{"trims everything", func() (float64, error) { return Compute([]int{1, 2}, 0.5, 0.5) }, ErrTrimsEverything},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.run(); !errors.Is(err, tc.wantErr) {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func BenchmarkComputeSymmetric(b *testing.B) {
	data := make([]float64, 10000)
	for i := range data {
		data[i] = float64((i * 7919) % 10000)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Compute(data, 0.05)
	}
}
