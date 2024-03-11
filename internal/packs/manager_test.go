package packs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetPackSizes(t *testing.T) {
	t.Run("Invalid input", func(t *testing.T) {
		m := NewManager()
		err := m.SetPackSizes([]int{})
		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid length of pack sizes")
	})

	t.Run("Valid sizes", func(t *testing.T) {
		m := NewManager()
		err := m.SetPackSizes([]int{500, 200, 1000, 300})
		assert.Nil(t, err)

		expected := []int{200, 300, 500, 1000}
		assert.Equal(t, expected, m.(*Manager).packSizes.sizes)
	})

}

func TestCalculatePacks(t *testing.T) {
	t.Run("Negative order", func(t *testing.T) {
		m := NewManager()
		packs, err := m.CalculatePacks(-1)
		assert.Empty(t, packs)
		assert.EqualError(t, err, "invalid item order")
	})

	t.Run("Zero order", func(t *testing.T) {
		m := NewManager()
		packs, err := m.CalculatePacks(0)
		assert.NotNil(t, packs)
		assert.Empty(t, packs)
		assert.NoError(t, err)
	})

	t.Run("Valid orders with defaul sizes", func(t *testing.T) {
		m := NewManager()

		testCases := []struct {
			input    int
			expected []Pack
		}{
			{
				input:    0,
				expected: []Pack{},
			},
			{
				input: 1,
				expected: []Pack{
					{Size: 250, Quantity: 1},
				},
			},
			{
				input: 250,
				expected: []Pack{
					{Size: 250, Quantity: 1},
				},
			},
			{
				input: 251,
				expected: []Pack{
					{Size: 500, Quantity: 1},
				},
			},
			{
				input: 501,
				expected: []Pack{
					{Size: 500, Quantity: 1},
					{Size: 250, Quantity: 1},
				},
			},
			{
				input: 12001,
				expected: []Pack{
					{Size: 5000, Quantity: 2},
					{Size: 2000, Quantity: 1},
					{Size: 250, Quantity: 1},
				},
			},
			{
				input: 500000,
				expected: []Pack{
					{Size: 5000, Quantity: 100},
				},
			},
		}

		for _, tc := range testCases {
			packs, err := m.CalculatePacks(tc.input)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, packs)
		}
	})

	t.Run("Valid orders with custom sizes", func(t *testing.T) {
		m := NewManager()
		err := m.SetPackSizes([]int{
			12, 45, 100, 234, 654,
		})
		assert.NoError(t, err)

		testCases := []struct {
			input    int
			expected []Pack
		}{
			{
				input: 1,
				expected: []Pack{
					{Size: 12, Quantity: 1},
				},
			},
			{
				input: 46,
				expected: []Pack{
					{Size: 12, Quantity: 4},
				},
			},
			{
				input: 53,
				expected: []Pack{
					{Size: 45, Quantity: 1},
					{Size: 12, Quantity: 1},
				},
			},
			{
				input: 234,
				expected: []Pack{
					{Size: 234, Quantity: 1},
				},
			},
			{
				input: 5002,
				expected: []Pack{
					{Size: 654, Quantity: 7},
					{Size: 234, Quantity: 1},
					{Size: 100, Quantity: 1},
					{Size: 45, Quantity: 2},
				},
			},
			{
				input: 200000,
				expected: []Pack{
					{Size: 654, Quantity: 303},
					{Size: 234, Quantity: 7},
					{Size: 100, Quantity: 2},
				},
			},
		}

		for _, tc := range testCases {
			packs, err := m.CalculatePacks(tc.input)
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.expected, packs)
		}
	})
}
