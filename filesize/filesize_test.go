package filesize

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testVariants(t *testing.T, exp int64, v int, units ...string) {
	t.Helper()

	t.Run(fmt.Sprintf("%d %s", v, units), func(t *testing.T) {
		var (
			s   FileSize
			err error
		)

		for _, unit := range units {
			err = s.UnmarshalText(fmt.Appendf(nil, "%d%s", v, unit))
			assert.Nil(t, err)
			assert.EqualValues(t, exp, s)
		}

		for _, unit := range units {
			err = s.UnmarshalText(fmt.Appendf(nil, "%d %s", v, unit))
			assert.Nil(t, err)
			assert.EqualValues(t, exp, s)
		}
	})
}

func TestFileSize_UnmarshalText_positive(t *testing.T) {
	var (
		s   FileSize
		err error
	)

	err = s.UnmarshalText([]byte("0"))
	assert.Nil(t, err)
	assert.EqualValues(t, 0, s)

	err = s.UnmarshalText([]byte("5"))
	assert.Nil(t, err)
	assert.EqualValues(t, 5, s)

	err = s.UnmarshalText([]byte("5,156"))
	assert.Nil(t, err)
	assert.EqualValues(t, 5156, s)

	err = s.UnmarshalText([]byte("  \t5156 \n"))
	assert.Nil(t, err)
	assert.EqualValues(t, 5156, s)

	err = s.UnmarshalText([]byte("5.12 mb"))
	assert.Nil(t, err)
	assert.EqualValues(t, int64(5.12*1000*1000), s)

	err = s.UnmarshalText([]byte("12.3456GiB"))
	assert.Nil(t, err)
	assert.EqualValues(t, int64(math.Round(12.3456*1024*1024*1024)), s)

	testVariants(t, 69, 69, "b", "B", "byte", "Byte", "BYTE")

	testVariants(t, 12*1000, 12, "kb", "kB", "kByte", "kilobyte", "KILOBYTE")
	testVariants(t, 7*1000*1000, 7, "mb", "MB", "MByte", "megabyte", "MEGABYTE")
	testVariants(t, 123*1000*1000*1000, 123, "gb", "GB", "GByte", "gigabyte", "GIGABYTE")
	testVariants(t, 34*1000*1000*1000*1000, 34, "tb", "TB", "TByte", "terrabyte", "TERRABYTE")
	testVariants(t, 34*1000*1000*1000*1000*1000, 34, "eb", "EB", "EByte", "exabyte", "EXABYTE")

	testVariants(t, 12*1024, 12, "kib", "kiB", "kiByte", "kibibyte", "KIBIBYTE")
	testVariants(t, 7*1024*1024, 7, "mib", "MiB", "MiByte", "mebibyte", "MEBIBYTE")
	testVariants(t, 123*1024*1024*1024, 123, "gib", "GiB", "GiByte", "gibibyte", "GIBIBYTE")
	testVariants(t, 34*1024*1024*1024*1024, 34, "tib", "TiB", "TiByte", "tebibyte", "TEBIBYTE")
	testVariants(t, 34*1024*1024*1024*1024*1024, 34, "eib", "EiB", "EiByte", "exibyte", "EXIBYTE")
}

func TestFileSize_UnmarshalText_negative(t *testing.T) {
	var (
		s   FileSize
		err error
	)

	t.Run("empty", func(t *testing.T) {
		vals := []string{"", "  \t\n"}

		for _, v := range vals {
			err = s.UnmarshalText([]byte(v))
			assert.Equal(t, FileSize(0), s)
			assert.ErrorContains(t, err, "empty value")
		}
	})

	t.Run("not starting with digit", func(t *testing.T) {
		vals := []string{
			"foo",
			"  foo",
			"kb",
			"megabyte",
		}

		for _, v := range vals {
			err = s.UnmarshalText([]byte(v))
			assert.Equal(t, FileSize(0), s)
			assert.ErrorContains(t, err, "must start with digit")
		}
	})

	t.Run("invalid unit", func(t *testing.T) {
		vals := []string{
			"3 foo",
			"3foo",
			"234 mbb",
			"123mbb",
			"123megab",
			"123 megabytes",
		}

		for _, v := range vals {
			err = s.UnmarshalText([]byte(v))
			assert.Equal(t, FileSize(0), s)
			assert.ErrorContains(t, err, "invalid unit")
		}
	})
}

func BenchmarkFileSize_UnmarshalText(b *testing.B) {
	var s FileSize

	for b.Loop() {
		_ = s.UnmarshalText([]byte("512 EiB"))
	}
}
