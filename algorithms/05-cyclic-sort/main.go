// Package cyclicsort демонстрирует паттерн "Циклическая сортировка"
package cyclicsort

// CyclicSort сортирует массив чисел в диапазоне [0, n-1] или [1, n]
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{3, 1, 5, 4, 2}
//	CyclicSort(arr) // Массив станет [1, 2, 3, 4, 5]
func CyclicSort(arr []int) {
	i := 0
	for i < len(arr) {
		// Правильная позиция для arr[i] — это arr[i] - 1 (для диапазона [1, n])
		// или arr[i] (для диапазона [0, n-1])
		correctIndex := arr[i] - 1

		// Если элемент не на своей позиции, меняем местами
		if arr[i] != arr[correctIndex] {
			arr[i], arr[correctIndex] = arr[correctIndex], arr[i]
		} else {
			// Элемент на правильной позиции — переходим к следующему
			i++
		}
	}
}

// FindMissingNumbers находит все пропущенные числа в массиве [1, n]
// Временная сложность: O(n)
// Пространственная сложность: O(1), не считая результата
//
// Пример:
//
//	arr := []int{2, 3, 1, 8, 2, 3, 5, 1}
//	result := FindMissingNumbers(arr) // Вернёт [4, 6, 7]
func FindMissingNumbers(arr []int) []int {
	// Сначала сортируем массив циклической сортировкой
	i := 0
	for i < len(arr) {
		correctIndex := arr[i] - 1
		// Проверяем, что индекс в пределах массива
		if arr[i] > 0 && arr[i] <= len(arr) && arr[i] != arr[correctIndex] {
			arr[i], arr[correctIndex] = arr[correctIndex], arr[i]
		} else {
			i++
		}
	}

	// Ищем пропущенные числа
	missing := []int{}
	for i := 0; i < len(arr); i++ {
		if arr[i] != i+1 {
			missing = append(missing, i+1)
		}
	}

	return missing
}

// FindDuplicate находит дубликат в массиве [1, n] с одним дубликатом
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{1, 4, 4, 3, 2}
//	result := FindDuplicate(arr) // Вернёт 4
func FindDuplicate(arr []int) int {
	i := 0
	for i < len(arr) {
		if arr[i] != i+1 {
			correctIndex := arr[i] - 1
			if arr[i] != arr[correctIndex] {
				arr[i], arr[correctIndex] = arr[correctIndex], arr[i]
			} else {
				// Нашли дубликат: элемент уже на правильной позиции
				return arr[i]
			}
		} else {
			i++
		}
	}

	return -1
}

// FindAllDuplicates находит все дубликаты в массиве [1, n]
// Временная сложность: O(n)
// Пространственная сложность: O(1), не считая результата
//
// Пример:
//
//	arr := []int{4, 3, 2, 7, 8, 2, 3, 1}
//	result := FindAllDuplicates(arr) // Вернёт [2, 3]
func FindAllDuplicates(arr []int) []int {
	// Сортируем циклической сортировкой
	i := 0
	for i < len(arr) {
		correctIndex := arr[i] - 1
		if arr[i] != arr[correctIndex] {
			arr[i], arr[correctIndex] = arr[correctIndex], arr[i]
		} else {
			i++
		}
	}

	// Ищем дубликаты (элементы не на своих позициях)
	duplicates := []int{}
	for i := 0; i < len(arr); i++ {
		if arr[i] != i+1 {
			duplicates = append(duplicates, arr[i])
		}
	}

	return duplicates
}

// FindFirstMissingPositive находит первое пропущенное положительное число
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{3, 4, -1, 1}
//	result := FindFirstMissingPositive(arr) // Вернёт 2
func FindFirstMissingPositive(arr []int) int {
	n := len(arr)

	// Шаг 1: Размещаем числа на правильных позициях
	// Число i должно быть на позиции i-1
	i := 0
	for i < n {
		// Правильная позиция для arr[i] — это arr[i] - 1
		// Но только если arr[i] в диапазоне [1, n]
		if arr[i] > 0 && arr[i] <= n && arr[i] != arr[arr[i]-1] {
			arr[i], arr[arr[i]-1] = arr[arr[i]-1], arr[i]
		} else {
			i++
		}
	}

	// Шаг 2: Ищем первое место, где число не на своей позиции
	for i := 0; i < n; i++ {
		if arr[i] != i+1 {
			return i + 1
		}
	}

	// Все числа от 1 до n присутствуют
	return n + 1
}

// FindCorruptPair находит пару (дубликат, пропущенное число) в массиве [1, n]
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	arr := []int{3, 1, 2, 5, 2}
//	duplicate, missing := FindCorruptPair(arr) // Вернёт (2, 4)
func FindCorruptPair(arr []int) (int, int) {
	// Сортируем циклической сортировкой
	i := 0
	for i < len(arr) {
		correctIndex := arr[i] - 1
		if arr[i] != arr[correctIndex] {
			arr[i], arr[correctIndex] = arr[correctIndex], arr[i]
		} else {
			i++
		}
	}

	// Ищем несоответствие
	for i := 0; i < len(arr); i++ {
		if arr[i] != i+1 {
			// arr[i] — дубликат, i+1 — пропущенное число
			return arr[i], i + 1
		}
	}

	return -1, -1
}

