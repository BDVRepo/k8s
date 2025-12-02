package twoheaps

import (
	"math"
	"testing"
)

// TestMedianFinder проверяет поиск медианы
func TestMedianFinder(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected float64
	}{
		{
			name:     "Нечётное количество",
			numbers:  []int{1, 2, 3, 4, 5},
			expected: 3.0,
		},
		{
			name:     "Чётное количество",
			numbers:  []int{1, 2, 3, 4},
			expected: 2.5, // (2 + 3) / 2
		},
		{
			name:     "Один элемент",
			numbers:  []int{5},
			expected: 5.0,
		},
		{
			name:     "Два элемента",
			numbers:  []int{1, 2},
			expected: 1.5,
		},
		{
			name:     "Неупорядоченные числа",
			numbers:  []int{5, 1, 3, 2, 4},
			expected: 3.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := NewMedianFinder()
			for _, num := range tt.numbers {
				finder.AddNum(num)
			}
			result := finder.FindMedian()

			// Сравниваем с небольшой погрешностью для float64
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("FindMedian() = %f; ожидалось %f",
					result, tt.expected)
			}
		})
	}
}

// TestFindMedianOfStream проверяет поиск медианы потока
func TestFindMedianOfStream(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected float64
	}{
		{
			name:     "Обычный случай",
			numbers:  []int{1, 2, 3, 4, 5},
			expected: 3.0,
		},
		{
			name:     "Чётное количество",
			numbers:  []int{1, 2, 3, 4},
			expected: 2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMedianOfStream(tt.numbers)
			if math.Abs(result-tt.expected) > 0.0001 {
				t.Errorf("FindMedianOfStream(%v) = %f; ожидалось %f",
					tt.numbers, result, tt.expected)
			}
		})
	}
}

// TestPartitionArray проверяет разделение массива
func TestPartitionArray(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		checkSum bool // Проверять ли суммы частей
	}{
		{
			name:     "Обычный случай",
			arr:      []int{1, 2, 3, 4, 5},
			checkSum: true,
		},
		{
			name:     "Чётное количество",
			arr:      []int{1, 2, 3, 4},
			checkSum: true,
		},
		{
			name:     "Один элемент",
			arr:      []int{5},
			checkSum: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part1, part2 := PartitionArray(tt.arr)

			// Проверяем, что все элементы присутствуют
			allElements := append(part1, part2...)
			if len(allElements) != len(tt.arr) {
				t.Errorf("PartitionArray(%v) разделил на %d элементов; ожидалось %d",
					tt.arr, len(allElements), len(tt.arr))
			}

			// Проверяем суммы, если требуется
			if tt.checkSum {
				sum1 := 0
				for _, v := range part1 {
					sum1 += v
				}
				sum2 := 0
				for _, v := range part2 {
					sum2 += v
				}

				// Суммы должны быть близки
				diff := sum1 - sum2
				if diff < 0 {
					diff = -diff
				}
			// Разница не должна быть слишком большой (кроме случая с одним элементом)
			if len(tt.arr) > 1 {
				totalSum := sum1 + sum2
				if diff > totalSum/2 {
					t.Errorf("PartitionArray(%v) создал части с большой разницей сумм: %d и %d",
						tt.arr, sum1, sum2)
				}
			}
			}
		})
	}
}

// TestSlidingWindowMedian проверяет медиану скользящего окна
func TestSlidingWindowMedian(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		k        int
		expected []float64
	}{
		{
			name:     "Обычный случай",
			nums:     []int{1, 3, -1, -3, 5, 3, 6, 7},
			k:        3,
			expected: []float64{1.0, -1.0, -1.0, 3.0, 5.0, 6.0},
		},
		{
			name:     "k равно длине массива",
			nums:     []int{1, 2, 3, 4, 5},
			k:        5,
			expected: []float64{3.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SlidingWindowMedian(tt.nums, tt.k)

			if len(result) != len(tt.expected) {
				t.Errorf("SlidingWindowMedian(%v, %d) вернул %d элементов; ожидалось %d",
					tt.nums, tt.k, len(result), len(tt.expected))
				return
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 0.0001 {
					t.Errorf("SlidingWindowMedian(%v, %d)[%d] = %f; ожидалось %f",
						tt.nums, tt.k, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

