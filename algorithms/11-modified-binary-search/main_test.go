package modifiedbinarysearch

import (
	"testing"
)

// TestSearchInRotatedArray проверяет поиск в ротированном массиве
func TestSearchInRotatedArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{
			name:     "Обычный случай",
			nums:     []int{4, 5, 6, 7, 0, 1, 2},
			target:   0,
			expected: 4,
		},
		{
			name:     "Элемент не найден",
			nums:     []int{4, 5, 6, 7, 0, 1, 2},
			target:   3,
			expected: -1,
		},
		{
			name:     "Не ротированный массив",
			nums:     []int{1, 2, 3, 4, 5},
			target:   3,
			expected: 2,
		},
		{
			name:     "Один элемент, найден",
			nums:     []int{1},
			target:   1,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchInRotatedArray(tt.nums, tt.target)
			if result != tt.expected {
				t.Errorf("SearchInRotatedArray(%v, %d) = %d; ожидалось %d",
					tt.nums, tt.target, result, tt.expected)
			}
		})
	}
}

// TestFindRange проверяет поиск диапазона
func TestFindRange(t *testing.T) {
	tests := []struct {
		name           string
		nums           []int
		target         int
		expectedFirst  int
		expectedLast   int
	}{
		{
			name:          "Обычный случай",
			nums:          []int{5, 7, 7, 8, 8, 10},
			target:        8,
			expectedFirst: 3,
			expectedLast:  4,
		},
		{
			name:          "Элемент не найден",
			nums:          []int{5, 7, 7, 8, 8, 10},
			target:        6,
			expectedFirst: -1,
			expectedLast:  -1,
		},
		{
			name:          "Один элемент",
			nums:          []int{1},
			target:        1,
			expectedFirst: 0,
			expectedLast:  0,
		},
		{
			name:          "Все элементы одинаковые",
			nums:          []int{2, 2, 2, 2},
			target:        2,
			expectedFirst: 0,
			expectedLast:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			first, last := FindRange(tt.nums, tt.target)
			if first != tt.expectedFirst || last != tt.expectedLast {
				t.Errorf("FindRange(%v, %d) = (%d, %d); ожидалось (%d, %d)",
					tt.nums, tt.target, first, last, tt.expectedFirst, tt.expectedLast)
			}
		})
	}
}

// TestFindPeakElement проверяет поиск пика
func TestFindPeakElement(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int // Ожидаемый индекс пика
	}{
		{
			name:     "Обычный случай",
			nums:     []int{1, 2, 3, 1},
			expected: 2, // Элемент 3
		},
		{
			name:     "Пик в начале",
			nums:     []int{3, 2, 1},
			expected: 0,
		},
		{
			name:     "Пик в конце",
			nums:     []int{1, 2, 3},
			expected: 2,
		},
		{
			name:     "Один элемент",
			nums:     []int{1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindPeakElement(tt.nums)
			if result != tt.expected {
				t.Errorf("FindPeakElement(%v) = %d; ожидалось %d",
					tt.nums, result, tt.expected)
			}

			// Проверяем, что найденный элемент действительно пик
			if result > 0 && result < len(tt.nums)-1 {
				if tt.nums[result] <= tt.nums[result-1] || tt.nums[result] <= tt.nums[result+1] {
					t.Errorf("FindPeakElement(%v) вернул индекс %d, который не является пиком",
						tt.nums, result)
				}
			}
		})
	}
}

// TestSearchMatrix проверяет поиск в матрице
func TestSearchMatrix(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		target   int
		expected bool
	}{
		{
			name: "Элемент найден",
			matrix: [][]int{
				{1, 4, 7, 11},
				{2, 5, 8, 12},
				{3, 6, 9, 16},
			},
			target:   5,
			expected: true,
		},
		{
			name: "Элемент не найден",
			matrix: [][]int{
				{1, 4, 7, 11},
				{2, 5, 8, 12},
				{3, 6, 9, 16},
			},
			target:   10,
			expected: false,
		},
		{
			name: "Пустая матрица",
			matrix: [][]int{},
			target:   5,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchMatrix(tt.matrix, tt.target)
			if result != tt.expected {
				t.Errorf("SearchMatrix(..., %d) = %v; ожидалось %v",
					tt.target, result, tt.expected)
			}
		})
	}
}

// TestFindMinInRotatedArray проверяет поиск минимума в ротированном массиве
func TestFindMinInRotatedArray(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "Обычный случай",
			nums:     []int{4, 5, 6, 7, 0, 1, 2},
			expected: 0,
		},
		{
			name:     "Не ротированный",
			nums:     []int{1, 2, 3, 4, 5},
			expected: 1,
		},
		{
			name:     "Один элемент",
			nums:     []int{1},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMinInRotatedArray(tt.nums)
			if result != tt.expected {
				t.Errorf("FindMinInRotatedArray(%v) = %d; ожидалось %d",
					tt.nums, result, tt.expected)
			}
		})
	}
}

// TestSearchInsertPosition проверяет поиск позиции для вставки
func TestSearchInsertPosition(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		target   int
		expected int
	}{
		{
			name:     "Элемент найден",
			nums:     []int{1, 3, 5, 6},
			target:   5,
			expected: 2,
		},
		{
			name:     "Элемент не найден, вставка в середину",
			nums:     []int{1, 3, 5, 6},
			target:   2,
			expected: 1,
		},
		{
			name:     "Вставка в конец",
			nums:     []int{1, 3, 5, 6},
			target:   7,
			expected: 4,
		},
		{
			name:     "Вставка в начало",
			nums:     []int{1, 3, 5, 6},
			target:   0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SearchInsertPosition(tt.nums, tt.target)
			if result != tt.expected {
				t.Errorf("SearchInsertPosition(%v, %d) = %d; ожидалось %d",
					tt.nums, tt.target, result, tt.expected)
			}
		})
	}
}

