package packs

import (
	"errors"
	"math"
	"sort"
	"sync"
)

type PacksManager interface {
	SetPackSizes(sizes []int) error
	CalculatePacks(itemOrder int) ([]Pack, error)
}

var defaultPackSizes = packSizes{
	sizes: []int{250, 500, 1000, 2000, 5000},
	l:     &sync.Mutex{},
}

type packSizes struct {
	sizes []int
	l     *sync.Mutex
}

type Manager struct {
	packSizes packSizes
}

func NewManager() PacksManager {

	m := Manager{
		packSizes: defaultPackSizes,
	}

	return &m
}

type Pack struct {
	Size     int
	Quantity int
}

// SetPackSizes sets and sorts the pack sizes in descending order for the Manager.
//
// It returns an error if the provided slice is empty.
func (m *Manager) SetPackSizes(sizes []int) error {
	if len(sizes) < 1 {
		return errors.New("invalid length of pack sizes")
	}

	m.packSizes.l.Lock()
	defer m.packSizes.l.Unlock()

	newSizes := make([]int, len(sizes))
	copy(newSizes, sizes)
	// sort in ascending order
	// since our algo expects a sorted array in asc order of pack sizes
	sort.Ints(newSizes)
	m.packSizes.sizes = newSizes

	return nil
}

// Returns an error if less then 0 items are ordered.
func (m *Manager) CalculatePacks(itemOrder int) ([]Pack, error) {
	packs := []Pack{}

	if itemOrder < 0 {
		return packs, errors.New("invalid item order")
	}

	// early exit
	if itemOrder == 0 {
		return packs, nil
	}

	results := m.dynamicCalculation(itemOrder)
	for s, q := range results {
		if q > 0 {
			packs = append(packs, Pack{
				Size:     s,
				Quantity: q,
			})
		}
	}

	// sort the output in ascending pack size order
	sort.Slice(packs, func(i, j int) bool {
		return packs[i].Size < packs[j].Size
	})

	return packs, nil
}

type packCombination struct {
	LeastExcess     int
	LeastPacks      int
	PackCombination map[int]int
}

func (m *Manager) dynamicCalculation(order int) map[int]int {
	if order < m.packSizes.sizes[0] {
		return map[int]int{
			m.packSizes.sizes[0]: 1,
		}
	}

	packSizes := m.packSizes.sizes

	// Initialize DP array
	dp := make([]packCombination, order+1)
	for i := range dp {
		dp[i].LeastExcess = math.MaxInt32
		dp[i].LeastPacks = math.MaxInt32
		dp[i].PackCombination = make(map[int]int)
	}

	// Base case: 0 order size
	dp[0] = packCombination{0, 0, make(map[int]int)}

	// Fill the DP table
	for i := 1; i <= order; i++ {
		for _, size := range packSizes {
			if size <= i {
				leftOver := i - size
				excess := dp[leftOver].LeastExcess
				packs := dp[leftOver].LeastPacks + 1
				combination := make(map[int]int)
				for k, v := range dp[leftOver].PackCombination {
					combination[k] = v
				}
				combination[size]++

				if excess < dp[i].LeastExcess || (excess == dp[i].LeastExcess && dp[leftOver].LeastPacks <= dp[i].LeastPacks) {
					dp[i].LeastPacks = packs
					dp[i].LeastExcess = excess
					dp[i].PackCombination = combination
				}
			} else {
				excess := size - i
				packs := 1
				if excess < dp[i].LeastExcess || (excess == dp[i].LeastExcess && packs < dp[i].LeastPacks) {
					dp[i].LeastPacks = 1
					dp[i].LeastExcess = size - i
					dp[i].PackCombination = map[int]int{size: 1}
				}
			}
		}
	}

	return dp[order].PackCombination
}
