package twopointers

import (
	"reflect"
	"testing"
)

// TestPairWithTargetSum проверяет поиск пары с целевой суммой
func TestPairWithTargetSum(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected []int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{1, 2, 3, 4, 6},
			target:   6,
			expected: []int{1, 3}, // 2 + 4 = 6
		},
		{
			name:     "Первый и последний элементы",
			arr:      []int{2, 5, 9, 11},
			target:   11,
			expected: []int{0, 2}, // 2 + 9 = 11
		},
		{
			name:     "Соседние элементы",
			arr:      []int{1, 2, 3, 4},
			target:   5,
			expected: []int{0, 3}, // 1 + 4 = 5 (алгоритм находит первую пару)
		},
		{
			name:     "Пара не найдена",
			arr:      []int{1, 2, 3, 4},
			target:   100,
			expected: []int{-1, -1},
		},
		{
			name:     "Пустой массив",
			arr:      []int{},
			target:   5,
			expected: []int{-1, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PairWithTargetSum(tt.arr, tt.target)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PairWithTargetSum(%v, %d) = %v; ожидалось %v",
					tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

// TestRemoveDuplicates проверяет удаление дубликатов
func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name           string
		arr            []int
		expectedLength int
		expectedArr    []int
	}{
		{
			name:           "Обычный случай",
			arr:            []int{2, 3, 3, 3, 6, 9, 9},
			expectedLength: 4,
			expectedArr:    []int{2, 3, 6, 9},
		},
		{
			name:           "Без дубликатов",
			arr:            []int{1, 2, 3, 4},
			expectedLength: 4,
			expectedArr:    []int{1, 2, 3, 4},
		},
		{
			name:           "Все одинаковые",
			arr:            []int{5, 5, 5, 5},
			expectedLength: 1,
			expectedArr:    []int{5},
		},
		{
			name:           "Пустой массив",
			arr:            []int{},
			expectedLength: 0,
			expectedArr:    []int{},
		},
		{
			name:           "Один элемент",
			arr:            []int{7},
			expectedLength: 1,
			expectedArr:    []int{7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)

			result := RemoveDuplicates(arrCopy)
			if result != tt.expectedLength {
				t.Errorf("RemoveDuplicates(%v) вернул длину %d; ожидалось %d",
					tt.arr, result, tt.expectedLength)
			}

			// Проверяем первые result элементов
			for i := 0; i < result && i < len(tt.expectedArr); i++ {
				if arrCopy[i] != tt.expectedArr[i] {
					t.Errorf("После RemoveDuplicates массив[%d] = %d; ожидалось %d",
						i, arrCopy[i], tt.expectedArr[i])
				}
			}
		})
	}
}

// TestSquareSortedArray проверяет возведение в квадрат
func TestSquareSortedArray(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{
			name:     "С отрицательными числами",
			arr:      []int{-2, -1, 0, 2, 3},
			expected: []int{0, 1, 4, 4, 9},
		},
		{
			name:     "Только положительные",
			arr:      []int{1, 2, 3, 4},
			expected: []int{1, 4, 9, 16},
		},
		{
			name:     "Только отрицательные",
			arr:      []int{-4, -3, -2, -1},
			expected: []int{1, 4, 9, 16},
		},
		{
			name:     "Пустой массив",
			arr:      []int{},
			expected: []int{},
		},
		{
			name:     "Один элемент",
			arr:      []int{-5},
			expected: []int{25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SquareSortedArray(tt.arr)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SquareSortedArray(%v) = %v; ожидалось %v",
					tt.arr, result, tt.expected)
			}
		})
	}
}

// TestThreeSum проверяет поиск троек с нулевой суммой
func TestThreeSum(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected [][]int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{-3, 0, 1, 2, -1, 1, -2},
			expected: [][]int{{-3, 1, 2}, {-2, 0, 2}, {-2, 1, 1}, {-1, 0, 1}},
		},
		{
			name:     "Нет троек",
			arr:      []int{1, 2, 3},
			expected: [][]int{},
		},
		{
			name:     "Три нуля",
			arr:      []int{0, 0, 0},
			expected: [][]int{{0, 0, 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ThreeSum(tt.arr)
			if len(result) != len(tt.expected) {
				t.Errorf("ThreeSum(%v) вернул %d троек; ожидалось %d",
					tt.arr, len(result), len(tt.expected))
			}
		})
	}
}

// TestIsPalindrome проверяет определение палиндрома
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "Палиндром нечётной длины",
			s:        "racecar",
			expected: true,
		},
		{
			name:     "Палиндром чётной длины",
			s:        "abba",
			expected: true,
		},
		{
			name:     "Не палиндром",
			s:        "hello",
			expected: false,
		},
		{
			name:     "Пустая строка",
			s:        "",
			expected: true,
		},
		{
			name:     "Один символ",
			s:        "a",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.s)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; ожидалось %v",
					tt.s, result, tt.expected)
			}
		})
	}
}

