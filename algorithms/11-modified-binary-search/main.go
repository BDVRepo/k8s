// Package modifiedbinarysearch демонстрирует паттерн "Модифицированный бинарный поиск"
package modifiedbinarysearch

// SearchInRotatedArray находит элемент в ротированном отсортированном массиве
// Временная сложность: O(log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	nums := []int{4, 5, 6, 7, 0, 1, 2}
//	target := 0
//	result := SearchInRotatedArray(nums, target) // Вернёт 4
func SearchInRotatedArray(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		}

		// Определяем, какая половина отсортирована
		if nums[left] <= nums[mid] {
			// Левая половина отсортирована
			if target >= nums[left] && target < nums[mid] {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			// Правая половина отсортирована
			if target > nums[mid] && target <= nums[right] {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}
	}

	return -1
}

// FindRange находит первую и последнюю позицию элемента в отсортированном массиве
// Временная сложность: O(log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	nums := []int{5, 7, 7, 8, 8, 10}
//	target := 8
//	first, last := FindRange(nums, target) // Вернёт (3, 4)
func FindRange(nums []int, target int) (int, int) {
	first := findFirst(nums, target)
	if first == -1 {
		return -1, -1
	}
	last := findLast(nums, target)
	return first, last
}

// findFirst находит первую позицию элемента
func findFirst(nums []int, target int) int {
	left, right := 0, len(nums)-1
	first := -1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			first = mid
			right = mid - 1 // Продолжаем искать слева
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return first
}

// findLast находит последнюю позицию элемента
func findLast(nums []int, target int) int {
	left, right := 0, len(nums)-1
	last := -1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			last = mid
			left = mid + 1 // Продолжаем искать справа
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return last
}

// FindPeakElement находит индекс пикового элемента
// Пик — элемент, который больше своих соседей
// Временная сложность: O(log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	nums := []int{1, 2, 3, 1}
//	result := FindPeakElement(nums) // Вернёт 2 (элемент 3)
func FindPeakElement(nums []int) int {
	left, right := 0, len(nums)-1

	for left < right {
		mid := left + (right-left)/2

		// Если элемент справа больше, пик справа
		if nums[mid] < nums[mid+1] {
			left = mid + 1
		} else {
			// Иначе пик слева или в mid
			right = mid
		}
	}

	return left
}

// SearchInfiniteArray находит элемент в "бесконечном" отсортированном массиве
// Временная сложность: O(log n), где n — позиция элемента
// Пространственная сложность: O(1)
//
// Пример:
//
//	// Предполагаем, что есть функция get(index) для получения элемента
//	// Для демонстрации используем слайс
//	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, ...}
//	target := 9
//	result := SearchInfiniteArray(arr, target) // Вернёт 4
func SearchInfiniteArray(arr []int, target int) int {
	// Сначала находим границы
	left, right := 0, 1

	// Расширяем правую границу, пока target не окажется в диапазоне
	for arr[right] < target {
		left = right
		right *= 2
		// В реальном коде нужно проверять границы массива
		if right >= len(arr) {
			right = len(arr) - 1
			break
		}
	}

	// Теперь выполняем обычный бинарный поиск
	return binarySearch(arr, target, left, right)
}

// binarySearch выполняет бинарный поиск в заданном диапазоне
func binarySearch(arr []int, target, left, right int) int {
	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}

// SearchMatrix находит элемент в 2D матрице, отсортированной по строкам и столбцам
// Временная сложность: O(m + n), где m и n — размеры матрицы
// Пространственная сложность: O(1)
//
// Пример:
//
//	matrix := [][]int{
//		{1, 4, 7, 11},
//		{2, 5, 8, 12},
//		{3, 6, 9, 16},
//	}
//	target := 5
//	result := SearchMatrix(matrix, target) // Вернёт true
func SearchMatrix(matrix [][]int, target int) bool {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return false
	}

	// Начинаем с правого верхнего угла
	row, col := 0, len(matrix[0])-1

	for row < len(matrix) && col >= 0 {
		if matrix[row][col] == target {
			return true
		} else if matrix[row][col] > target {
			// Текущий элемент больше — двигаемся влево
			col--
		} else {
			// Текущий элемент меньше — двигаемся вниз
			row++
		}
	}

	return false
}

// FindMinInRotatedArray находит минимальный элемент в ротированном отсортированном массиве
// Временная сложность: O(log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	nums := []int{4, 5, 6, 7, 0, 1, 2}
//	result := FindMinInRotatedArray(nums) // Вернёт 0
func FindMinInRotatedArray(nums []int) int {
	left, right := 0, len(nums)-1

	for left < right {
		mid := left + (right-left)/2

		// Если правая половина отсортирована, минимум слева
		if nums[mid] < nums[right] {
			right = mid
		} else {
			// Иначе минимум справа
			left = mid + 1
		}
	}

	return nums[left]
}

// SearchInsertPosition находит позицию для вставки элемента в отсортированный массив
// Временная сложность: O(log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	nums := []int{1, 3, 5, 6}
//	target := 5
//	result := SearchInsertPosition(nums, target) // Вернёт 2
func SearchInsertPosition(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	// Если не нашли, возвращаем позицию для вставки
	return left
}

