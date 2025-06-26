package filesize

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"

	"github.com/zekrotja/parsables/internal/utils"
)

// FileSize implements encoding.TextUnmarshaler for file sizes.
//
// Parsable formats follow the following expression:
//
//	^\s*[0-9,\.]+\s*[a-zA-Z]+$
//
// Examples:
//
//	200
//	42 Gigabyte
//	2 Tebibyte
//	2.34 GiB
//	2,456,789kb
//
// Allowed units are as following (case insensitive):
//   - b, byte > factor 1000^0
//   - kb, kbyte, kilobyte > factor 1000^1
//   - kib, bibyte, kibibyte > factor 1024^1
//   - mb, mbyte, megabyte > factor 1000^2
//   - mib, mibyte, mibibyte > factor 1024^2
//   - gb, gbyte, gigabyte > factor 1000^3
//   - gib, gibyte, gibibyte > factor 1024^3
//   - tb, tbyte, terrabyte > factor 1000^4
//   - tib, tibyte, tebibyte > factor 1024^4
//   - eb, ebyte, exabyte > factor 1000^5
//   - eib, eibyte, exibyte > factor 1024^5
type FileSize int64

// FromString parses FileSize from the given string s.
// For format details, see documentation of FileSize.
func FromString(s string) (FileSize, error) {
	runes := []rune(s)

	if len(runes) == 0 {
		return 0, errors.New("invalid format: empty value")
	}

	rdr := utils.NewReader(runes)

	if cont := rdr.SkipWhile(unicode.IsSpace); !cont {
		return 0, errors.New("invalid format: empty value")
	}

	r, _ := rdr.Peek()
	if !unicode.IsDigit(r) {
		return 0, errors.New("invalid format: must start with digit")
	}

	var bufNumber, bufUnit bytes.Buffer
	bufNumber.Grow(len(s))

	for {
		r, ok := rdr.Take()
		if !ok {
			break
		}
		if r == ',' || unicode.IsSpace(r) {
			continue
		}
		if r != '.' && !unicode.IsDigit(r) {
			rdr.Back()
			break
		}
		bufNumber.WriteRune(r)
	}

	for {
		r, ok := rdr.Take()
		if !ok {
			break
		}
		if unicode.IsSpace(r) {
			continue
		}
		bufUnit.WriteRune(unicode.ToLower(r))
	}

	nb, err := strconv.ParseFloat(bufNumber.String(), 64)
	if err != nil {
		return 0, err
	}

	if bufUnit.Len() == 0 {
		return FileSize(math.Round(nb)), nil
	}

	switch bufUnit.String() {
	case "b", "byte":
		nb *= math.Pow(1000, 0)

	case "kb", "kbyte", "kilobyte":
		nb *= math.Pow(1000, 1)
	case "mb", "mbyte", "megabyte":
		nb *= math.Pow(1000, 2)
	case "gb", "gbyte", "gigabyte":
		nb *= math.Pow(1000, 3)
	case "tb", "tbyte", "terrabyte":
		nb *= math.Pow(1000, 4)
	case "eb", "ebyte", "exabyte":
		nb *= math.Pow(1000, 5)

	case "kib", "kibyte", "kibibyte":
		nb *= math.Pow(1024, 1)
	case "mib", "mibyte", "mebibyte":
		nb *= math.Pow(1024, 2)
	case "gib", "gibyte", "gibibyte":
		nb *= math.Pow(1024, 3)
	case "tib", "tibyte", "tebibyte":
		nb *= math.Pow(1024, 4)
	case "eib", "eibyte", "exibyte":
		nb *= math.Pow(1024, 5)

	default:
		return 0, fmt.Errorf("invalid format: invalid unit '%s'", bufUnit.String())
	}

	return FileSize(math.Round(nb)), nil
}

func (t *FileSize) UnmarshalText(p []byte) (err error) {
	*t, err = FromString(string(p))
	return err
}
