// Package subsets демонстрирует паттерн "Подмножества"
package subsets

// Subsets генерирует все подмножества массива
// Временная сложность: O(2^n)
// Пространственная сложность: O(2^n)
//
// Пример:
//
//	nums := []int{1, 2, 3}
//	result := Subsets(nums)
//	// Вернёт [[], [1], [2], [1,2], [3], [1,3], [2,3], [1,2,3]]
func Subsets(nums []int) [][]int {
	result := [][]int{}
	current := []int{}

	var backtrack func(start int)
	backtrack = func(start int) {
		// Добавляем текущее подмножество в результат
		subset := make([]int, len(current))
		copy(subset, current)
		result = append(result, subset)

		// Генерируем подмножества, начиная с каждого следующего элемента
		for i := start; i < len(nums); i++ {
			// Добавляем элемент
			current = append(current, nums[i])
			// Рекурсивно генерируем подмножества с этим элементом
			backtrack(i + 1)
			// Убираем элемент (backtracking)
			current = current[:len(current)-1]
		}
	}

	backtrack(0)
	return result
}

// SubsetsWithDup генерирует все подмножества массива с дубликатами
// Временная сложность: O(2^n)
// Пространственная сложность: O(2^n)
//
// Пример:
//
//	nums := []int{1, 2, 2}
//	result := SubsetsWithDup(nums)
//	// Вернёт [[], [1], [2], [1,2], [2,2], [1,2,2]]
func SubsetsWithDup(nums []int) [][]int {
	// Сначала сортируем массив
	sorted := make([]int, len(nums))
	copy(sorted, nums)
	// В реальном коде используйте sort.Ints(sorted)

	result := [][]int{}
	current := []int{}

	var backtrack func(start int)
	backtrack = func(start int) {
		subset := make([]int, len(current))
		copy(subset, current)
		result = append(result, subset)

		for i := start; i < len(sorted); i++ {
			// Пропускаем дубликаты на том же уровне
			if i > start && sorted[i] == sorted[i-1] {
				continue
			}

			current = append(current, sorted[i])
			backtrack(i + 1)
			current = current[:len(current)-1]
		}
	}

	backtrack(0)
	return result
}

// Permute генерирует все перестановки массива
// Временная сложность: O(n!)
// Пространственная сложность: O(n!)
//
// Пример:
//
//	nums := []int{1, 2, 3}
//	result := Permute(nums)
//	// Вернёт [[1,2,3], [1,3,2], [2,1,3], [2,3,1], [3,1,2], [3,2,1]]
func Permute(nums []int) [][]int {
	result := [][]int{}
	current := []int{}
	used := make([]bool, len(nums))

	var backtrack func()
	backtrack = func() {
		// Если текущая перестановка завершена
		if len(current) == len(nums) {
			perm := make([]int, len(current))
			copy(perm, current)
			result = append(result, perm)
			return
		}

		// Пробуем каждый элемент
		for i := 0; i < len(nums); i++ {
			if !used[i] {
				used[i] = true
				current = append(current, nums[i])
				backtrack()
				current = current[:len(current)-1]
				used[i] = false
			}
		}
	}

	backtrack()
	return result
}

// Combine генерирует все комбинации размера k из n элементов
// Временная сложность: O(C(n,k))
// Пространственная сложность: O(C(n,k))
//
// Пример:
//
//	n := 4
//	k := 2
//	result := Combine(n, k)
//	// Вернёт [[1,2], [1,3], [1,4], [2,3], [2,4], [3,4]]
func Combine(n int, k int) [][]int {
	result := [][]int{}
	current := []int{}

	var backtrack func(start int)
	backtrack = func(start int) {
		// Если комбинация готова
		if len(current) == k {
			comb := make([]int, len(current))
			copy(comb, current)
			result = append(result, comb)
			return
		}

		// Пробуем элементы от start до n
		for i := start; i <= n; i++ {
			current = append(current, i)
			backtrack(i + 1)
			current = current[:len(current)-1]
		}
	}

	backtrack(1)
	return result
}

// GenerateParenthesis генерирует все валидные комбинации n пар скобок
// Временная сложность: O(4^n / sqrt(n))
// Пространственная сложность: O(4^n / sqrt(n))
//
// Пример:
//
//	n := 3
//	result := GenerateParenthesis(3)
//	// Вернёт ["((()))", "(()())", "(())()", "()(())", "()()()"]
func GenerateParenthesis(n int) []string {
	result := []string{}
	current := ""

	var backtrack func(open, close int)
	backtrack = func(open, close int) {
		// Если достигли нужной длины
		if len(current) == 2*n {
			result = append(result, current)
			return
		}

		// Можем добавить открывающую скобку, если их меньше n
		if open < n {
			current += "("
			backtrack(open+1, close)
			current = current[:len(current)-1]
		}

		// Можем добавить закрывающую скобку, если их меньше открывающих
		if close < open {
			current += ")"
			backtrack(open, close+1)
			current = current[:len(current)-1]
		}
	}

	backtrack(0, 0)
	return result
}

// LetterCombinations генерирует все комбинации букв для телефонного номера
// Временная сложность: O(4^n), где n — длина digits
// Пространственная сложность: O(4^n)
//
// Пример:
//
//	digits := "23"
//	result := LetterCombinations(digits)
//	// Вернёт ["ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"]
func LetterCombinations(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}

	// Маппинг цифр на буквы
	digitMap := map[byte]string{
		'2': "abc",
		'3': "def",
		'4': "ghi",
		'5': "jkl",
		'6': "mno",
		'7': "pqrs",
		'8': "tuv",
		'9': "wxyz",
	}

	result := []string{}
	current := ""

	var backtrack func(index int)
	backtrack = func(index int) {
		// Если обработали все цифры
		if index == len(digits) {
			result = append(result, current)
			return
		}

		// Получаем буквы для текущей цифры
		letters := digitMap[digits[index]]
		for i := 0; i < len(letters); i++ {
			current += string(letters[i])
			backtrack(index + 1)
			current = current[:len(current)-1]
		}
	}

	backtrack(0)
	return result
}

