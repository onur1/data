package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/onur1/data"
	"github.com/onur1/data/nilable"
	"github.com/onur1/data/result"
	"github.com/stretchr/testify/assert"
)

var (
	errFailed  = errors.New("failed")
	errWrapped = fmt.Errorf("wrapped: %w", errFailed)
)

func TestResult(t *testing.T) {
	testCases := []struct {
		desc        string
		result      data.Result[int]
		expected    int
		expectedErr error
	}{
		{
			desc:     "Ok",
			result:   result.Ok(42),
			expected: 42,
		},
		{
			desc:        "Error",
			result:      result.Error[int](errFailed),
			expectedErr: errFailed,
		},
		{
			desc:     "Map (succeed)",
			result:   result.Map(result.Ok(1), double),
			expected: 2,
		},
		{
			desc:        "Map (fail)",
			result:      result.Map(result.Error[int](errFailed), double),
			expectedErr: errFailed,
		},
		{
			desc:     "MapError (succeed)",
			result:   result.MapError(result.Ok(42), wrappedError),
			expected: 42,
		},
		{
			desc:        "MapError (fail)",
			result:      result.MapError(result.Error[int](errFailed), wrappedError),
			expectedErr: errWrapped,
		},
		{
			desc: "Bimap (succeed)",
			result: result.Map(
				result.Bimap(result.Ok(-1), wrappedError, isPositive),
				func(n bool) int {
					if n {
						return 1
					} else {
						return 2
					}
				},
			),
			expected: 2,
		},
		{
			desc: "Bimap (fail)",
			result: result.Map(
				result.Bimap(result.Error[int](errFailed), wrappedError, isPositive),
				func(n bool) int {
					if n {
						return 1
					} else {
						return 2
					}
				},
			),
			expectedErr: errWrapped,
		},
		{
			desc:     "Ap (succeed)",
			result:   result.Ap(result.Ok(double), result.Ok(42)),
			expected: 84,
		},
		{
			desc:        "Ap (fail)",
			result:      result.Ap(result.Ok(double), result.Error[int](errFailed)),
			expectedErr: errFailed,
		},
		{
			desc:     "ApFirst (succeed)",
			result:   result.ApFirst(result.Ok(1), result.Ok(2)),
			expected: 1,
		},
		{
			desc:        "ApFirst (fail)",
			result:      result.ApFirst(result.Error[int](errFailed), result.Ok(2)),
			expectedErr: errFailed,
		},
		{
			desc:     "ApSecond (succeed)",
			result:   result.ApSecond(result.Ok(1), result.Ok(2)),
			expected: 2,
		},
		{
			desc:        "ApSecond (fail)",
			result:      result.ApSecond(result.Ok(1), result.Error[int](errFailed)),
			expectedErr: errFailed,
		},
		{
			desc: "Chain (succeed)",
			result: result.Chain(result.Ok(42), func(a int) data.Result[int] {
				return result.Ok(a + 1)
			}),
			expected: 43,
		},
		{
			desc: "Chain (fail)",
			result: result.Chain(result.Error[int](errFailed), func(a int) data.Result[int] {
				return result.Ok(a + 1)
			}),
			expectedErr: errFailed,
		},
		{
			desc: "FromNilable (some)",
			result: result.FromNilable(nilable.Some(42), func() error {
				return errFailed
			}),
			expected: 42,
		},
		{
			desc: "FromNilable (nil)",
			result: result.FromNilable(nilable.Nil[int](), func() error {
				return errFailed
			}),
			expectedErr: errFailed,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assertEq(t, tC.result, tC.expected, tC.expectedErr)
		})
	}
}

func assertEq(t *testing.T, res data.Result[int], expected int, expectedErr error) {
	x, err := res()
	if err != nil {
		assert.Equal(t, expectedErr, err)
	} else {
		assert.Equal(t, expected, x)
	}
}

func double(n int) int {
	return n * 2
}

func isPositive(n int) bool {
	return n > 0
}

func wrappedError(err error) error {
	return fmt.Errorf("wrapped: %w", err)
}
