package utils

import (
	"testing"
)

func TestSplitSize(t *testing.T) {
	cases := []struct {
		size     int64
		expected Sizes
	}{
		{size: 0, expected: Sizes{0, 0, 0, 0, 0, 0}},
		{size: 1, expected: Sizes{0, 0, 0, 0, 0, 1}},
		{size: 2, expected: Sizes{0, 0, 0, 0, 0, 2}},
		// ...
		{size: 7, expected: Sizes{0, 0, 0, 0, 0, 7}},
		{size: 8, expected: Sizes{8, 0, 0, 0, 0, 0}},
		{size: 9, expected: Sizes{8, 0, 0, 0, 0, 1}},
		{size: 10, expected: Sizes{8, 0, 0, 0, 0, 2}},
		{size: 11, expected: Sizes{8, 0, 0, 0, 0, 3}},
		// ...
		{size: 15, expected: Sizes{8, 0, 0, 0, 0, 7}},
		{size: 16, expected: Sizes{8, 8, 0, 0, 0, 0}},
		{size: 17, expected: Sizes{8, 8, 0, 0, 0, 1}},
		// ...
		{size: 23, expected: Sizes{8, 8, 0, 0, 0, 7}},
		{size: 24, expected: Sizes{8, 8, 8, 0, 0, 0}},
		{size: 25, expected: Sizes{8, 8, 8, 0, 0, 1}},
		// ...
		{size: 47, expected: Sizes{8, 8, 8, 8, 8, 7}},
		{size: 48, expected: Sizes{8, 8, 8, 8, 8, 8}},
		{size: 49, expected: Sizes{8, 8, 8, 8, 8, 9}},
		// ...
		{size: 95, expected: Sizes{16, 16, 16, 16, 16, 15}},
		{size: 96, expected: Sizes{16, 16, 16, 16, 16, 16}},
		{size: 97, expected: Sizes{16, 16, 16, 16, 16, 17}},
		// ...
		{size: 103, expected: Sizes{16, 16, 16, 16, 16, 23}},
		{size: 104, expected: Sizes{24, 16, 16, 16, 16, 16}},
		{size: 105, expected: Sizes{24, 16, 16, 16, 16, 17}},
		// ...
		{size: 143, expected: Sizes{24, 24, 24, 24, 24, 23}},
		{size: 144, expected: Sizes{24, 24, 24, 24, 24, 24}},
		{size: 145, expected: Sizes{24, 24, 24, 24, 24, 25}},
	}

	for _, testCase := range cases {
		actual := SplitSize(testCase.size)
		if actual != testCase.expected {
			t.Errorf("With size %d. Expected: %v, but Actual: %v",
				testCase.size, testCase.expected, actual)
		}
	}
}
