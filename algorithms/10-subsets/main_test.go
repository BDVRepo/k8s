package subsets

import (
	"sort"
	"testing"
)

// TestSubsets проверяет генерацию подмножеств
func TestSubsets(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int // Количество подмножеств должно быть 2^n
	}{
		{
			name:     "Три элемента",
			nums:     []int{1, 2, 3},
			expected: 8, // 2^3
		},
		{
			name:     "Один элемент",
			nums:     []int{1},
			expected: 2, // 2^1
		},
		{
			name:     "Пустой массив",
			nums:     []int{},
			expected: 1, // 2^0 (только пустое множество)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subsets(tt.nums)
			if len(result) != tt.expected {
				t.Errorf("Subsets(%v) вернул %d подмножеств; ожидалось %d",
					tt.nums, len(result), tt.expected)
			}

			// Проверяем, что пустое множество присутствует
			hasEmpty := false
			for _, subset := range result {
				if len(subset) == 0 {
					hasEmpty = true
					break
				}
			}
			if !hasEmpty {
				t.Errorf("Subsets(%v) не содержит пустого множества", tt.nums)
			}
		})
	}
}

// TestPermute проверяет генерацию перестановок
func TestPermute(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int // Количество перестановок должно быть n!
	}{
		{
			name:     "Три элемента",
			nums:     []int{1, 2, 3},
			expected: 6, // 3!
		},
		{
			name:     "Два элемента",
			nums:     []int{1, 2},
			expected: 2, // 2!
		},
		{
			name:     "Один элемент",
			nums:     []int{1},
			expected: 1, // 1!
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Permute(tt.nums)
			if len(result) != tt.expected {
				t.Errorf("Permute(%v) вернул %d перестановок; ожидалось %d",
					tt.nums, len(result), tt.expected)
			}
		})
	}
}

// TestCombine проверяет генерацию комбинаций
func TestCombine(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		k        int
		expected int // Количество комбинаций C(n,k)
	}{
		{
			name:     "C(4,2)",
			n:        4,
			k:        2,
			expected: 6, // C(4,2) = 6
		},
		{
			name:     "C(5,3)",
			n:        5,
			k:        3,
			expected: 10, // C(5,3) = 10
		},
		{
			name:     "k=1",
			n:        4,
			k:        1,
			expected: 4,
		},
		{
			name:     "k=n",
			n:        4,
			k:        4,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Combine(tt.n, tt.k)
			if len(result) != tt.expected {
				t.Errorf("Combine(%d, %d) вернул %d комбинаций; ожидалось %d",
					tt.n, tt.k, len(result), tt.expected)
			}

			// Проверяем, что все комбинации имеют размер k
			for _, comb := range result {
				if len(comb) != tt.k {
					t.Errorf("Combine(%d, %d) вернул комбинацию размера %d; ожидалось %d",
						tt.n, tt.k, len(comb), tt.k)
				}
			}
		})
	}
}

// TestGenerateParenthesis проверяет генерацию скобок
func TestGenerateParenthesis(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int // Количество валидных комбинаций
	}{
		{
			name:     "n=3",
			n:        3,
			expected: 5,
		},
		{
			name:     "n=2",
			n:        2,
			expected: 2,
		},
		{
			name:     "n=1",
			n:        1,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateParenthesis(tt.n)
			if len(result) != tt.expected {
				t.Errorf("GenerateParenthesis(%d) вернул %d комбинаций; ожидалось %d",
					tt.n, len(result), tt.expected)
			}

			// Проверяем, что все комбинации валидны
			for _, paren := range result {
				if !isValidParenthesis(paren) {
					t.Errorf("GenerateParenthesis(%d) вернул невалидную комбинацию: %s",
						tt.n, paren)
				}
			}
		})
	}
}

// isValidParenthesis проверяет валидность скобок
func isValidParenthesis(s string) bool {
	count := 0
	for _, char := range s {
		if char == '(' {
			count++
		} else {
			count--
		}
		if count < 0 {
			return false
		}
	}
	return count == 0
}

// TestLetterCombinations проверяет генерацию комбинаций букв
func TestLetterCombinations(t *testing.T) {
	tests := []struct {
		name     string
		digits   string
		expected int
	}{
		{
			name:     "Две цифры",
			digits:   "23",
			expected: 9, // 3 * 3
		},
		{
			name:     "Одна цифра",
			digits:   "2",
			expected: 3,
		},
		{
			name:     "Пустая строка",
			digits:   "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LetterCombinations(tt.digits)
			if len(result) != tt.expected {
				t.Errorf("LetterCombinations(%q) вернул %d комбинаций; ожидалось %d",
					tt.digits, len(result), tt.expected)
			}
		})
	}
}

// TestSubsetsWithDup проверяет подмножества с дубликатами
func TestSubsetsWithDup(t *testing.T) {
	tests := []struct {
		name     string
		nums     []int
		expected int
	}{
		{
			name:     "С дубликатами",
			nums:     []int{1, 2, 2},
			expected: 6, // Меньше чем 2^3 = 8 из-за дубликатов
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SubsetsWithDup(tt.nums)
			if len(result) != tt.expected {
				t.Errorf("SubsetsWithDup(%v) вернул %d подмножеств; ожидалось %d",
					tt.nums, len(result), tt.expected)
			}

			// Проверяем, что нет дубликатов в результате
			seen := make(map[string]bool)
			for _, subset := range result {
				// Сортируем для создания ключа
				sorted := make([]int, len(subset))
				copy(sorted, subset)
				sort.Ints(sorted)

				key := ""
				for _, v := range sorted {
					key += string(rune('0' + v))
				}

				if seen[key] {
					t.Errorf("SubsetsWithDup(%v) вернул дубликат: %v", tt.nums, subset)
				}
				seen[key] = true
			}
		})
	}
}

