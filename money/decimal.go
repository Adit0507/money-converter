package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Decimal struct {
	subunits  int64
	precision byte
}

const (
	// returned if the decimal is malformed
	ErrInvalidDecimal = Error("unable to convert the deciaml")

	// returned if the quantity is too arge - this would cause floating point precision errors.
	ErrTooLarge = Error("quantity over 10^12 is too large")
)

func ParseDecimal(value string) (Decimal, error) {
	intPart, fracPart, _ := strings.Cut(value, ".")

	const maxDecimal = 12
	if len(intPart) > maxDecimal {
		return Decimal{}, ErrTooLarge
	}

	subunits, err := strconv.ParseInt(intPart+fracPart, 10, 64)
	if err != nil {
		return Decimal{}, fmt.Errorf("%w: %s", ErrInvalidDecimal, err.Error())
	}

	precision := byte(len(fracPart))
	return Decimal{subunits: subunits, precision: precision}, nil
}

func (d *Decimal) simplify() {
	for d.subunits%10 == 0 && d.precision > 0 {
		d.precision--
		d.subunits /= 10
	}
}

func pow10(power int) int {
	switch power {
	case 0:
		return 1
	case 1:
		return 10
	case 2:
		return 100
	case 3:
		return 1000
	default:
		return int(math.Pow(10, float64(power)))
	}
}
