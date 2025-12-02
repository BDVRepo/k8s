package cyclicsort

import (
	"reflect"
	"sort"
	"testing"
)

// TestCyclicSort проверяет циклическую сортировку
func TestCyclicSort(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{3, 1, 5, 4, 2},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Уже отсортирован",
			arr:      []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Обратный порядок",
			arr:      []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Один элемент",
			arr:      []int{1},
			expected: []int{1},
		},
		{
			name:     "Два элемента",
			arr:      []int{2, 1},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)
			CyclicSort(arrCopy)
			if !reflect.DeepEqual(arrCopy, tt.expected) {
				t.Errorf("CyclicSort(%v) = %v; ожидалось %v",
					tt.arr, arrCopy, tt.expected)
			}
		})
	}
}

// TestFindMissingNumbers проверяет поиск пропущенных чисел
func TestFindMissingNumbers(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{2, 3, 1, 8, 2, 3, 5, 1},
			expected: []int{4, 6, 7},
		},
		{
			name:     "Нет пропущенных",
			arr:      []int{1, 2, 3, 4, 5},
			expected: []int{},
		},
		{
			name:     "Пропущено первое",
			arr:      []int{2, 3, 4, 5},
			expected: []int{1},
		},
		{
			name:     "Пропущено последнее",
			arr:      []int{1, 2, 3, 4, 6, 7}, // Длина 6, но нет 5
			expected: []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)
			result := FindMissingNumbers(arrCopy)

			// Сортируем для сравнения, так как порядок может отличаться
			sort.Ints(result)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindMissingNumbers(%v) = %v; ожидалось %v",
					tt.arr, result, tt.expected)
			}
		})
	}
}

// TestFindDuplicate проверяет поиск дубликата
func TestFindDuplicate(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{1, 4, 4, 3, 2},
			expected: 4,
		},
		{
			name:     "Дубликат в начале",
			arr:      []int{2, 1, 3, 3, 5, 4},
			expected: 3,
		},
		{
			name:     "Дубликат в конце",
			arr:      []int{2, 4, 1, 4, 4},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)
			result := FindDuplicate(arrCopy)
			if result != tt.expected {
				t.Errorf("FindDuplicate(%v) = %d; ожидалось %d",
					tt.arr, result, tt.expected)
			}
		})
	}
}

// TestFindAllDuplicates проверяет поиск всех дубликатов
func TestFindAllDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected []int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{4, 3, 2, 7, 8, 2, 3, 1},
			expected: []int{2, 3},
		},
		{
			name:     "Один дубликат",
			arr:      []int{1, 2, 3, 2, 4},
			expected: []int{2},
		},
		{
			name:     "Нет дубликатов",
			arr:      []int{1, 2, 3, 4, 5},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)
			result := FindAllDuplicates(arrCopy)

			// Сортируем для сравнения
			sort.Ints(result)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindAllDuplicates(%v) = %v; ожидалось %v",
					tt.arr, result, tt.expected)
			}
		})
	}
}

// TestFindFirstMissingPositive проверяет поиск первого пропущенного положительного
func TestFindFirstMissingPositive(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		expected int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{3, 4, -1, 1},
			expected: 2,
		},
		{
			name:     "Все положительные",
			arr:      []int{1, 2, 0},
			expected: 3,
		},
		{
			name:     "Большие числа",
			arr:      []int{7, 8, 9, 11, 12},
			expected: 1,
		},
		{
			name:     "Отрицательные и нули",
			arr:      []int{-1, -2, 0},
			expected: 1,
		},
		{
			name:     "Последовательность",
			arr:      []int{1, 2, 3, 4},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)
			result := FindFirstMissingPositive(arrCopy)
			if result != tt.expected {
				t.Errorf("FindFirstMissingPositive(%v) = %d; ожидалось %d",
					tt.arr, result, tt.expected)
			}
		})
	}
}

// TestFindCorruptPair проверяет поиск пары (дубликат, пропущенное)
func TestFindCorruptPair(t *testing.T) {
	tests := []struct {
		name           string
		arr            []int
		expectedDup    int
		expectedMissing int
	}{
		{
			name:           "Обычный случай",
			arr:            []int{3, 1, 2, 5, 2},
			expectedDup:    2,
			expectedMissing: 4,
		},
		{
			name:           "Другой случай",
			arr:            []int{3, 1, 2, 3, 6, 4},
			expectedDup:    3,
			expectedMissing: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arrCopy := make([]int, len(tt.arr))
			copy(arrCopy, tt.arr)
			dup, missing := FindCorruptPair(arrCopy)
			if dup != tt.expectedDup || missing != tt.expectedMissing {
				t.Errorf("FindCorruptPair(%v) = (%d, %d); ожидалось (%d, %d)",
					tt.arr, dup, missing, tt.expectedDup, tt.expectedMissing)
			}
		})
	}
}

