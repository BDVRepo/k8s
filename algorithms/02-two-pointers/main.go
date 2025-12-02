// Package twopointers демонстрирует паттерн "Два указателя"
package twopointers

// PairWithTargetSum находит индексы двух чисел в отсортированном массиве,
// которые дают целевую сумму
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{1, 2, 3, 4, 6}
//	target := 6
//	result := PairWithTargetSum(arr, target) // Вернёт [1, 3] (числа 2 и 4)
func PairWithTargetSum(arr []int, target int) []int {
	// Два указателя: один в начале, другой в конце
	left := 0
	right := len(arr) - 1

	for left < right {
		currentSum := arr[left] + arr[right]

		if currentSum == target {
			// Нашли пару
			return []int{left, right}
		}

		if currentSum < target {
			// Сумма меньше целевой — сдвигаем левый указатель вправо
			// (чтобы увеличить сумму, так как массив отсортирован)
			left++
		} else {
			// Сумма больше целевой — сдвигаем правый указатель влево
			// (чтобы уменьшить сумму)
			right--
		}
	}

	// Пара не найдена
	return []int{-1, -1}
}

// RemoveDuplicates удаляет дубликаты из отсортированного массива на месте
// и возвращает длину уникальной части
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{2, 3, 3, 3, 6, 9, 9}
//	length := RemoveDuplicates(arr) // Вернёт 4, массив станет [2, 3, 6, 9, ...]
func RemoveDuplicates(arr []int) int {
	if len(arr) == 0 {
		return 0
	}

	// nextNonDuplicate указывает на позицию для следующего уникального элемента
	nextNonDuplicate := 1

	for i := 1; i < len(arr); i++ {
		// Если текущий элемент отличается от предыдущего уникального
		if arr[nextNonDuplicate-1] != arr[i] {
			arr[nextNonDuplicate] = arr[i]
			nextNonDuplicate++
		}
	}

	return nextNonDuplicate
}

// SquareSortedArray возвращает квадраты элементов отсортированного массива
// в отсортированном порядке
// Временная сложность: O(n)
// Пространственная сложность: O(n) для результата
//
// Пример:
//
//	arr := []int{-2, -1, 0, 2, 3}
//	result := SquareSortedArray(arr) // Вернёт [0, 1, 4, 4, 9]
func SquareSortedArray(arr []int) []int {
	n := len(arr)
	if n == 0 {
		return []int{}
	}

	// Создаём массив результата
	squares := make([]int, n)

	// Два указателя с разных концов
	left := 0
	right := n - 1

	// Заполняем массив с конца (максимальные квадраты)
	highestSquareIdx := n - 1

	for left <= right {
		leftSquare := arr[left] * arr[left]
		rightSquare := arr[right] * arr[right]

		// Выбираем больший квадрат и помещаем в конец результата
		if leftSquare > rightSquare {
			squares[highestSquareIdx] = leftSquare
			left++
		} else {
			squares[highestSquareIdx] = rightSquare
			right--
		}
		highestSquareIdx--
	}

	return squares
}

// ThreeSum находит все уникальные тройки чисел с нулевой суммой
// Временная сложность: O(n²)
// Пространственная сложность: O(n) для сортировки
//
// Пример:
//
//	arr := []int{-3, 0, 1, 2, -1, 1, -2}
//	result := ThreeSum(arr) // Вернёт [[-3, 1, 2], [-2, 0, 2], [-2, 1, 1], [-1, 0, 1]]
func ThreeSum(arr []int) [][]int {
	// Сначала сортируем массив
	sortArray(arr)

	triplets := [][]int{}

	for i := 0; i < len(arr)-2; i++ {
		// Пропускаем дубликаты для первого числа
		if i > 0 && arr[i] == arr[i-1] {
			continue
		}

		// Ищем пары, которые дадут -arr[i]
		searchPair(arr, -arr[i], i+1, &triplets)
	}

	return triplets
}

// searchPair - вспомогательная функция для поиска пар с целевой суммой
func searchPair(arr []int, targetSum int, left int, triplets *[][]int) {
	right := len(arr) - 1

	for left < right {
		currentSum := arr[left] + arr[right]

		if currentSum == targetSum {
			// Нашли тройку
			*triplets = append(*triplets, []int{-targetSum, arr[left], arr[right]})
			left++
			right--

			// Пропускаем дубликаты
			for left < right && arr[left] == arr[left-1] {
				left++
			}
			for left < right && arr[right] == arr[right+1] {
				right--
			}
		} else if currentSum < targetSum {
			left++
		} else {
			right--
		}
	}
}

// sortArray - простая сортировка пузырьком для демонстрации
// В реальном коде используйте sort.Ints
func sortArray(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// IsPalindrome проверяет, является ли строка палиндромом
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	s := "racecar"
//	result := IsPalindrome(s) // Вернёт true
func IsPalindrome(s string) bool {
	left := 0
	right := len(s) - 1

	for left < right {
		if s[left] != s[right] {
			return false
		}
		left++
		right--
	}

	return true
}

