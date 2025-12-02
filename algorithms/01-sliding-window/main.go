// Package slidingwindow демонстрирует паттерн "Скользящее окно"
package slidingwindow

// MaxSumSubarray находит максимальную сумму подмассива размера k
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{2, 1, 5, 1, 3, 2}
//	k := 3
//	result := MaxSumSubarray(arr, k) // Вернёт 9 (подмассив [5, 1, 3])
func MaxSumSubarray(arr []int, k int) int {
	// Проверка граничных условий
	if len(arr) < k || k <= 0 {
		return 0
	}

	// Вычисляем сумму первого окна
	windowSum := 0
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}

	maxSum := windowSum

	// Скользим окном по массиву
	// На каждом шаге: добавляем новый элемент справа, убираем элемент слева
	for i := k; i < len(arr); i++ {
		windowSum += arr[i] - arr[i-k] // Добавляем новый, убираем старый
		if windowSum > maxSum {
			maxSum = windowSum
		}
	}

	return maxSum
}

// MinSubarrayLen находит минимальную длину подмассива с суммой >= target
// Использует динамическое скользящее окно
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{2, 1, 5, 2, 3, 2}
//	target := 7
//	result := MinSubarrayLen(arr, target) // Вернёт 2 (подмассив [5, 2])
func MinSubarrayLen(arr []int, target int) int {
	// Проверка граничных условий
	if len(arr) == 0 {
		return 0
	}

	minLength := len(arr) + 1 // Инициализируем значением больше максимально возможного
	windowSum := 0
	windowStart := 0

	for windowEnd := 0; windowEnd < len(arr); windowEnd++ {
		// Расширяем окно, добавляя элемент справа
		windowSum += arr[windowEnd]

		// Сжимаем окно слева, пока сумма >= target
		for windowSum >= target {
			currentLength := windowEnd - windowStart + 1
			if currentLength < minLength {
				minLength = currentLength
			}
			// Убираем элемент слева и сдвигаем начало окна
			windowSum -= arr[windowStart]
			windowStart++
		}
	}

	// Если не нашли подходящий подмассив
	if minLength > len(arr) {
		return 0
	}

	return minLength
}

// LongestSubstringKDistinct находит длину самой длинной подстроки с не более чем k различными символами
// Временная сложность: O(n)
// Пространственная сложность: O(k) для хранения символов в окне
//
// Пример:
//
//	s := "araaci"
//	k := 2
//	result := LongestSubstringKDistinct(s, k) // Вернёт 4 ("araa")
func LongestSubstringKDistinct(s string, k int) int {
	if len(s) == 0 || k == 0 {
		return 0
	}

	// Хеш-таблица для подсчёта символов в текущем окне
	charCount := make(map[byte]int)
	maxLength := 0
	windowStart := 0

	for windowEnd := 0; windowEnd < len(s); windowEnd++ {
		// Добавляем символ справа в окно
		rightChar := s[windowEnd]
		charCount[rightChar]++

		// Сжимаем окно, пока различных символов больше k
		for len(charCount) > k {
			leftChar := s[windowStart]
			charCount[leftChar]--
			if charCount[leftChar] == 0 {
				delete(charCount, leftChar) // Удаляем символ из мапы, если его счётчик = 0
			}
			windowStart++
		}

		// Обновляем максимальную длину
		currentLength := windowEnd - windowStart + 1
		if currentLength > maxLength {
			maxLength = currentLength
		}
	}

	return maxLength
}

// FindAverages вычисляет средние значения для всех подмассивов размера k
// Временная сложность: O(n)
// Пространственная сложность: O(n) для результата
//
// Пример:
//
//	arr := []int{1, 3, 2, 6, -1, 4, 1, 8, 2}
//	k := 5
//	result := FindAverages(arr, k) // Вернёт [2.2, 2.8, 2.4, 3.6, 2.8]
func FindAverages(arr []int, k int) []float64 {
	if len(arr) < k || k <= 0 {
		return []float64{}
	}

	result := make([]float64, len(arr)-k+1)
	windowSum := 0

	// Считаем сумму первого окна
	for i := 0; i < k; i++ {
		windowSum += arr[i]
	}
	result[0] = float64(windowSum) / float64(k)

	// Скользим по массиву
	for i := k; i < len(arr); i++ {
		windowSum += arr[i] - arr[i-k]
		result[i-k+1] = float64(windowSum) / float64(k)
	}

	return result
}

