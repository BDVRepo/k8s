package slidingwindow

import (
	"testing"
)

// TestMaxSumSubarray проверяет функцию поиска максимальной суммы подмассива
func TestMaxSumSubarray(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		k        int
		expected int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{2, 1, 5, 1, 3, 2},
			k:        3,
			expected: 9, // [5, 1, 3]
		},
		{
			name:     "Максимум в начале",
			arr:      []int{9, 1, 2, 3, 4, 5},
			k:        2,
			expected: 10, // [9, 1]
		},
		{
			name:     "Максимум в конце",
			arr:      []int{1, 2, 3, 4, 5, 9},
			k:        2,
			expected: 14, // [5, 9]
		},
		{
			name:     "K равно длине массива",
			arr:      []int{1, 2, 3},
			k:        3,
			expected: 6,
		},
		{
			name:     "K больше длины массива",
			arr:      []int{1, 2},
			k:        3,
			expected: 0,
		},
		{
			name:     "Пустой массив",
			arr:      []int{},
			k:        1,
			expected: 0,
		},
		{
			name:     "Отрицательные числа",
			arr:      []int{-1, -2, 5, 3, -1},
			k:        2,
			expected: 8, // [5, 3]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxSumSubarray(tt.arr, tt.k)
			if result != tt.expected {
				t.Errorf("MaxSumSubarray(%v, %d) = %d; ожидалось %d",
					tt.arr, tt.k, result, tt.expected)
			}
		})
	}
}

// TestMinSubarrayLen проверяет функцию поиска минимального подмассива
func TestMinSubarrayLen(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{
			name:     "Обычный случай",
			arr:      []int{2, 1, 5, 2, 3, 2},
			target:   7,
			expected: 2, // [5, 2]
		},
		{
			name:     "Один элемент достаточен",
			arr:      []int{2, 1, 5, 2, 8},
			target:   7,
			expected: 1, // [8]
		},
		{
			name:     "Весь массив нужен",
			arr:      []int{1, 1, 1, 1, 1},
			target:   5,
			expected: 5,
		},
		{
			name:     "Невозможно достичь суммы",
			arr:      []int{1, 2, 3},
			target:   100,
			expected: 0,
		},
		{
			name:     "Пустой массив",
			arr:      []int{},
			target:   5,
			expected: 0,
		},
		{
			name:     "Большой target в начале",
			arr:      []int{3, 4, 1, 1, 6},
			target:   8,
			expected: 3, // [3, 4, 1]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinSubarrayLen(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("MinSubarrayLen(%v, %d) = %d; ожидалось %d",
					tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

// TestLongestSubstringKDistinct проверяет функцию поиска самой длинной подстроки
func TestLongestSubstringKDistinct(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		k        int
		expected int
	}{
		{
			name:     "Обычный случай",
			s:        "araaci",
			k:        2,
			expected: 4, // "araa"
		},
		{
			name:     "K больше количества уникальных символов",
			s:        "araaci",
			k:        10,
			expected: 6, // вся строка
		},
		{
			name:     "K равно 1",
			s:        "araaci",
			k:        1,
			expected: 2, // "aa"
		},
		{
			name:     "Пустая строка",
			s:        "",
			k:        2,
			expected: 0,
		},
		{
			name:     "K равно 0",
			s:        "abc",
			k:        0,
			expected: 0,
		},
		{
			name:     "Все символы одинаковые",
			s:        "aaaa",
			k:        1,
			expected: 4,
		},
		{
			name:     "Длинная строка cbbebi",
			s:        "cbbebi",
			k:        3,
			expected: 5, // "cbbeb"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LongestSubstringKDistinct(tt.s, tt.k)
			if result != tt.expected {
				t.Errorf("LongestSubstringKDistinct(%q, %d) = %d; ожидалось %d",
					tt.s, tt.k, result, tt.expected)
			}
		})
	}
}

// TestFindAverages проверяет функцию вычисления средних значений
func TestFindAverages(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		k        int
		expected []float64
	}{
		{
			name:     "Обычный случай",
			arr:      []int{1, 3, 2, 6, -1, 4, 1, 8, 2},
			k:        5,
			expected: []float64{2.2, 2.8, 2.4, 3.6, 2.8},
		},
		{
			name:     "K равно длине массива",
			arr:      []int{1, 2, 3},
			k:        3,
			expected: []float64{2.0},
		},
		{
			name:     "K больше длины массива",
			arr:      []int{1, 2},
			k:        5,
			expected: []float64{},
		},
		{
			name:     "Пустой массив",
			arr:      []int{},
			k:        2,
			expected: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindAverages(tt.arr, tt.k)
			if len(result) != len(tt.expected) {
				t.Errorf("FindAverages(%v, %d) вернул %d элементов; ожидалось %d",
					tt.arr, tt.k, len(result), len(tt.expected))
				return
			}
			for i := range result {
				// Сравниваем с небольшой погрешностью для float64
				diff := result[i] - tt.expected[i]
				if diff < -0.0001 || diff > 0.0001 {
					t.Errorf("FindAverages(%v, %d)[%d] = %f; ожидалось %f",
						tt.arr, tt.k, i, result[i], tt.expected[i])
				}
			}
		})
	}
}


