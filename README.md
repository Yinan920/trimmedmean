# trimmedmean

A small, dependency-free Go package for computing **symmetric and asymmetric
trimmed means** of integer or floating-point samples.

A *trimmed mean* sorts the data, discards a proportion of the smallest and
largest values, and averages what remains. Removing the extremes makes the
estimate resistant to outliers, which is why the trimmed mean is a common
*robust* alternative to the ordinary mean.

> Replace `Yinan920` below with your own GitHub username throughout.

## Installation

```bash
go get github.com/Yinan920/trimmedmean
```

Requires Go 1.18 or newer (the package uses generics).

## Usage

```go
package main

import (
	"fmt"

	"github.com/Yinan920/trimmedmean"
)

func main() {
	data := []float64{2, 4, 6, 8, 10, 1000} // one large outlier

	plain, _ := trimmedmean.Compute(data)            // ordinary mean
	sym, _ := trimmedmean.Compute(data, 0.10)        // 10% off each end
	asym, _ := trimmedmean.Compute(data, 0.10, 0.20) // 10% low, 20% high

	fmt.Println(plain, sym, asym)
}
```

## API

```go
func Compute[T Number](data []T, proportions ...float64) (float64, error)
```

`T` may be any integer or floating-point type. The variadic `proportions`
argument controls trimming:

| Arguments              | Behaviour                                              |
| ---------------------- | ------------------------------------------------------ |
| `Compute(data)`        | ordinary, untrimmed mean                               |
| `Compute(data, p)`     | **symmetric** trim: `p` from the low *and* the high end |
| `Compute(data, lo, hi)`| **asymmetric** trim: `lo` from the low end, `hi` from the high end |

The number of observations removed at an end is `floor(n * proportion)`. This
matches R's `mean(x, trim = p)` for the symmetric case, so results are directly
comparable. The input slice is never modified.

Each proportion must lie in `[0, 0.5]`, at most two may be given, and the
trimming must leave at least one observation. Violations return one of the
exported sentinel errors (`ErrEmptyData`, `ErrTooManyArgs`,
`ErrProportionRange`, `ErrTrimsEverything`), testable with `errors.Is`.

## Testing

```bash
go test ./...            # run unit tests
go test -v ./...         # verbose
go test -bench=. -benchmem ./...   # benchmark + allocation report
```

The test suite covers the untrimmed case, symmetric and asymmetric trimming,
the R floor rule, input immutability, and every error path.

## Repositories

- This package: `github.com/Yinan920/trimmedmean`
- Example program that uses it: `github.com/Yinan920/trimmedmean-demo`

## GenAI Tools

See the demo repository's README for the project-wide note on generative-AI use.

## License

MIT — see [LICENSE](LICENSE).
