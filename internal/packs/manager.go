package packs

import (
	"errors"
	"math"
	"sort"
	"sync"
)

var defaultPackSizes = packSizes{
	sizes: []int{5000, 2000, 1000, 500, 250},
	l:     &sync.Mutex{},
}

type packSizes struct {
	sizes []int
	l     *sync.Mutex
}

type Manager struct {
	packSizes packSizes
}

func NewManager() *Manager {

	m := Manager{
		packSizes: defaultPackSizes,
	}

	return &m
}

type Pack struct {
	Size     int
	Quantity int
}

func (m *Manager) SetPackSizes(sizes []int) error {
	if len(sizes) < 1 {
		return errors.New("invalid lenght of pack sizes")
	}

	m.packSizes.l.Lock()
	defer m.packSizes.l.Unlock()

	m.packSizes.sizes = sizes
	sort.Sort(sort.Reverse(sort.IntSlice(m.packSizes.sizes)))

	return nil
}

func (m *Manager) CalculatePacks(itemOrder int) ([]Pack, error) {
	packs := []Pack{}

	if itemOrder < 0 {
		return packs, errors.New("invalid item order")
	}

	results := m.calculatePacksCombination(itemOrder)
	for s, q := range results {
		if q > 0 {
			packs = append(packs, Pack{
				Size:     s,
				Quantity: q,
			})
		}
	}

	return packs, nil
}

// packInfo stores the number of packs and the total items for a given order size
type packInfo struct {
	NumPacks   int
	TotalItems int
}

func (m *Manager) calculatePacksDynamic(order int) map[int]int {
	sort.Ints(m.packSizes.sizes)

	// Initialize the dynamic programming table
	dp := make([]packInfo, order+1)
	for i := range dp {
		dp[i] = packInfo{NumPacks: math.MaxInt32, TotalItems: 0}
	}
	dp[0] = packInfo{NumPacks: 0, TotalItems: 0}

	for _, size := range m.packSizes.sizes {
		for i := size; i <= order; i++ {
			if dp[i-size].NumPacks+1 < dp[i].NumPacks || (dp[i-size].NumPacks+1 == dp[i].NumPacks && dp[i-size].TotalItems+size < dp[i].TotalItems) {
				dp[i].NumPacks = dp[i-size].NumPacks + 1
				dp[i].TotalItems = dp[i-size].TotalItems + size
			}
		}
	}

	// Backtrack to find the combination of packs
	result := make(map[int]int)
	for i := order; i > 0; {
		for _, size := range m.packSizes.sizes {
			if i-size >= 0 && dp[i].NumPacks == dp[i-size].NumPacks+1 && dp[i].TotalItems == dp[i-size].TotalItems+size {
				result[size]++
				i -= size
				break
			}
		}
	}

	return result
}

func (m *Manager) calculatePacksCombination(order int) map[int]int {
	sort.Sort(sort.Reverse(sort.IntSlice(m.packSizes.sizes)))
	bestCombination := make(map[int]int)
	leastExcess := order
	leastPacks := order

	// Recursive function to find the best combination of packs
	var findCombination func(int, map[int]int, int, int)
	findCombination = func(index int, currentCombination map[int]int, currentTotal int, currentPacks int) {
		if currentTotal >= order {
			excess := currentTotal - order
			if excess < leastExcess || (excess == leastExcess && currentPacks < leastPacks) {
				leastExcess = excess
				leastPacks = currentPacks
				// Copy currentCombination to bestCombination
				bestCombination = make(map[int]int)
				for k, v := range currentCombination {
					bestCombination[k] = v
				}
			}
			return
		}

		for i := index; i < len(m.packSizes.sizes); i++ {
			size := m.packSizes.sizes[i]
			currentCombination[size]++
			findCombination(i, currentCombination, currentTotal+size, currentPacks+1)
			currentCombination[size]--
		}
	}

	findCombination(0, make(map[int]int), 0, 0)

	return bestCombination
}
