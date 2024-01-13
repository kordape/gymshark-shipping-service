package packs

import (
	"fmt"
	"testing"
)

func TestCalculatePacks(t *testing.T) {
	t.Run("Test 1", func(t *testing.T) {
		m := NewManager()
		m.SetPackSizes([]int{31, 23, 53})
		packs, _ := m.CalculatePacks(263)
		fmt.Println(packs)
	})
}
