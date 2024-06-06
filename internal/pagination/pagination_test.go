package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateOffset(t *testing.T) {
	type testCase struct {
		name     string
		page     int
		size     int
		expected int
	}

	testCases := []testCase{
		{
			name:     "valid input",
			page:     2,
			size:     10,
			expected: 10,
		},
		{
			name:     "page is zero",
			page:     0,
			size:     10,
			expected: 0,
		},
		{
			name:     "size is zero",
			page:     2,
			size:     0,
			expected: 0,
		},
		{
			name:     "page and size are zero",
			page:     0,
			size:     0,
			expected: 0,
		},
		{
			name:     "page is negative",
			page:     -1,
			size:     10,
			expected: 0,
		},
		{
			name:     "size is negative",
			page:     2,
			size:     -10,
			expected: 0,
		},
		{
			name:     "page and size are negative",
			page:     -1,
			size:     -10,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateOffset(tc.page, tc.size)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCalculateLastPage(t *testing.T) {
	type testCase struct {
		name     string
		count    int
		pageSize int
		expected int
	}

	testCases := []testCase{
		{
			name:     "exact division",
			count:    100,
			pageSize: 10,
			expected: 10,
		},
		{
			name:     "division with remainder",
			count:    105,
			pageSize: 10,
			expected: 11,
		},
		{
			name:     "count is less than pageSize",
			count:    5,
			pageSize: 10,
			expected: 1,
		},
		{
			name:     "pageSize is zero",
			count:    100,
			pageSize: 0,
			expected: 1,
		},
		{
			name:     "pageSize is negative",
			count:    100,
			pageSize: -10,
			expected: 1,
		},
		{
			name:     "count is zero",
			count:    0,
			pageSize: 10,
			expected: 1,
		},
		{
			name:     "count and pageSize are zero",
			count:    0,
			pageSize: 0,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CalculateLastPage(tc.count, tc.pageSize)
			assert.Equal(t, tc.expected, result)
		})
	}
}
