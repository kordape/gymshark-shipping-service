package packs

import (
	"errors"
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

// SetPackSizes sets and sorts the pack sizes in descending order for the Manager.
//
// It returns an error if the provided slice is empty.
func (m *Manager) SetPackSizes(sizes []int) error {
	if len(sizes) < 1 {
		return errors.New("invalid length of pack sizes")
	}

	m.packSizes.l.Lock()
	defer m.packSizes.l.Unlock()

	m.packSizes.sizes = sizes
	// sort in descending order
	// since our algo expects a sorted array in desc order of pack sizes
	sort.Sort(sort.Reverse(sort.IntSlice(m.packSizes.sizes)))

	return nil
}

// CalculatePacks calculates the optimal number of packs needed to fullfil the item order
//
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

	results := m.calculatePacksCombination(itemOrder)
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

// calculatePacksCombination computes the optimal combination of pack sizes to fulfill a given order.
// The function aims to minimize the number of packs used and the excess number of items sent, while ensuring the order is fully met.
// If the order is smaller than the smallest pack size, it simply returns one pack of the smallest size.
//
// It uses a recursive approach to explore all possible combinations of the available pack sizes.
// During the recursion, it keeps track of the combination that results in the least number of excess items and the least number of packs.
// Once all combinations have been considered, it returns the most efficient combination found.
//
// NOTE: brute force approach since we are not expecting large number of pack sizes
// TODO: there is a more optimized solution using Dynamic Programming
func (m *Manager) calculatePacksCombination(order int) map[int]int {
	bestCombination := make(map[int]int)

	// early exit
	if order < m.packSizes.sizes[len(m.packSizes.sizes)-1] {
		bestCombination[m.packSizes.sizes[len(m.packSizes.sizes)-1]] = 1

		return bestCombination
	}

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
