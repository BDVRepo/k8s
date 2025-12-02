package topkelements

import (
	"reflect"
	"sort"
	"testing"
)

// TestFindKLargest проверяет поиск K наибольших элементов
func TestFindKLargest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{
			name:     "Обычный случай",
			nums:     []int{3, 1, 5, 12, 2, 11},
			k:        3,
			expected: []int{12, 11, 5},
		},
		{
			name:     "k равно длине массива",
			nums:     []int{1, 2, 3},
			k:        3,
			expected: []int{3, 2, 1},
		},
		{
			name:     "k=1",
			nums:     []int{3, 1, 5, 12, 2},
			k:        1,
			expected: []int{12},
		},
		{
			name:     "Пустой массив",
			nums:     []int{},
			k:        3,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindKLargest(tt.nums, tt.k)
			sort.Sort(sort.Reverse(sort.IntSlice(result)))
			sort.Sort(sort.Reverse(sort.IntSlice(tt.expected)))

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindKLargest(%v, %d) = %v; ожидалось %v",
					tt.nums, tt.k, result, tt.expected)
			}
		})
	}
}

// TestFindKthLargest проверяет поиск K-го наибольшего элемента
func TestFindKthLargest(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected int
	}{
		{
			name:     "Обычный случай",
			nums:     []int{3, 2, 1, 5, 6, 4},
			k:        2,
			expected: 5,
		},
		{
			name:     "k=1",
			nums:     []int{3, 2, 3, 1, 2, 4, 5, 5, 6},
			k:        4,
			expected: 4,
		},
		{
			name:     "k равно длине",
			nums:     []int{3, 2, 1},
			k:        3,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindKthLargest(tt.nums, tt.k)
			if result != tt.expected {
				t.Errorf("FindKthLargest(%v, %d) = %d; ожидалось %d",
					tt.nums, tt.k, result, tt.expected)
			}
		})
	}
}

// TestTopKFrequent проверяет поиск K наиболее частых элементов
func TestTopKFrequent(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []int
	}{
		{
			name:     "Обычный случай",
			nums:     []int{1, 1, 1, 2, 2, 3},
			k:        2,
			expected: []int{1, 2},
		},
		{
			name:     "k=1",
			nums:     []int{1},
			k:        1,
			expected: []int{1},
		},
		{
			name:     "Все элементы одинаковые",
			nums:     []int{1, 1, 1, 1},
			k:        1,
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TopKFrequent(tt.nums, tt.k)
			sort.Ints(result)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TopKFrequent(%v, %d) = %v; ожидалось %v",
					tt.nums, tt.k, result, tt.expected)
			}
		})
	}
}

// TestKClosestPoints проверяет поиск K ближайших точек
func TestKClosestPoints(t *testing.T) {
	tests := []struct {
		name     string
		points   [][]int
		k        int
		expected int // Количество точек в результате
	}{
		{
			name:     "Обычный случай",
			points:   [][]int{{1, 3}, {-2, 2}, {5, 8}},
			k:        1,
			expected: 1,
		},
		{
			name:     "k равно количеству точек",
			points:   [][]int{{1, 3}, {-2, 2}},
			k:        2,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KClosestPoints(tt.points, tt.k)
			if len(result) != tt.expected {
				t.Errorf("KClosestPoints(..., %d) вернул %d точек; ожидалось %d",
					tt.k, len(result), tt.expected)
			}
		})
	}
}

// TestTopKFrequentWords проверяет поиск K наиболее частых слов
func TestTopKFrequentWords(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		k        int
		expected int // Количество слов в результате
	}{
		{
			name:     "Обычный случай",
			words:    []string{"i", "love", "leetcode", "i", "love", "coding"},
			k:        2,
			expected: 2,
		},
		{
			name:     "k=1",
			words:    []string{"the", "day", "is", "sunny", "the", "the", "the", "sunny", "is", "is"},
			k:        1,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TopKFrequentWords(tt.words, tt.k)
			if len(result) != tt.expected {
				t.Errorf("TopKFrequentWords(..., %d) вернул %d слов; ожидалось %d",
					tt.k, len(result), tt.expected)
			}
		})
	}
}

