package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateAverage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ratings := []int{10, 20, 30}
		result := CalculateAverage(ratings)
		assert.Equal(t, 20, result)
	})

	t.Run("single element", func(t *testing.T) {
		assert.Equal(t, 5, CalculateAverage([]int{5}))
	})

	t.Run("empty slice ", func(t *testing.T) {
		assert.Equal(t, 0, CalculateAverage([]int{}))
	})
}
