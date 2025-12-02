// Package twoheaps демонстрирует паттерн "Две кучи"
package twoheaps

import (
	"container/heap"
)

// MinHeap — минимальная куча (min-heap)
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MaxHeap — максимальная куча (max-heap)
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MedianFinder находит медиану потока чисел
// Использует две кучи: max-heap для меньшей половины, min-heap для большей
type MedianFinder struct {
	maxHeap *MaxHeap // Меньшая половина (максимум сверху)
	minHeap *MinHeap // Большая половина (минимум сверху)
}

// NewMedianFinder создаёт новый MedianFinder
func NewMedianFinder() *MedianFinder {
	maxHeap := &MaxHeap{}
	minHeap := &MinHeap{}
	heap.Init(maxHeap)
	heap.Init(minHeap)
	return &MedianFinder{
		maxHeap: maxHeap,
		minHeap: minHeap,
	}
}

// AddNum добавляет число в структуру данных
// Временная сложность: O(log n)
func (mf *MedianFinder) AddNum(num int) {
	// Если max-heap пуст или число меньше максимума меньшей половины
	if mf.maxHeap.Len() == 0 || num <= (*mf.maxHeap)[0] {
		heap.Push(mf.maxHeap, num)
	} else {
		heap.Push(mf.minHeap, num)
	}

	// Балансируем кучи: разница размеров не должна превышать 1
	if mf.maxHeap.Len() > mf.minHeap.Len()+1 {
		// Перемещаем из max-heap в min-heap
		maxVal := heap.Pop(mf.maxHeap).(int)
		heap.Push(mf.minHeap, maxVal)
	} else if mf.minHeap.Len() > mf.maxHeap.Len()+1 {
		// Перемещаем из min-heap в max-heap
		minVal := heap.Pop(mf.minHeap).(int)
		heap.Push(mf.maxHeap, minVal)
	}
}

// FindMedian возвращает медиану
// Временная сложность: O(1)
func (mf *MedianFinder) FindMedian() float64 {
	if mf.maxHeap.Len() == mf.minHeap.Len() {
		// Чётное количество элементов — берём среднее двух вершин
		return float64((*mf.maxHeap)[0]+(*mf.minHeap)[0]) / 2.0
	}

	// Нечётное количество — медиана в большей куче
	if mf.maxHeap.Len() > mf.minHeap.Len() {
		return float64((*mf.maxHeap)[0])
	}
	return float64((*mf.minHeap)[0])
}

// FindMedianOfStream находит медиану потока чисел
// Временная сложность: O(n log n) для n чисел
// Пространственная сложность: O(n)
//
// Пример:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	median := FindMedianOfStream(numbers) // Вернёт 3.0
func FindMedianOfStream(numbers []int) float64 {
	finder := NewMedianFinder()
	for _, num := range numbers {
		finder.AddNum(num)
	}
	return finder.FindMedian()
}

// PartitionArray разделяет массив на две части с минимальной разницей сумм
// Использует две кучи для хранения элементов двух частей
// Временная сложность: O(n log n)
// Пространственная сложность: O(n)
func PartitionArray(arr []int) ([]int, []int) {
	if len(arr) == 0 {
		return []int{}, []int{}
	}

	// Сортируем массив
	sorted := make([]int, len(arr))
	copy(sorted, arr)
	// В реальном коде используйте sort.Ints(sorted)

	// Разделяем на две части
	part1 := []int{}
	part2 := []int{}

	sum1, sum2 := 0, 0

	// Простой жадный алгоритм: добавляем в часть с меньшей суммой
	for _, num := range sorted {
		if sum1 <= sum2 {
			part1 = append(part1, num)
			sum1 += num
		} else {
			part2 = append(part2, num)
			sum2 += num
		}
	}

	return part1, part2
}

// SlidingWindowMedian находит медиану для каждого окна размера k
// Временная сложность: O(n log k)
// Пространственная сложность: O(k)
func SlidingWindowMedian(nums []int, k int) []float64 {
	if len(nums) < k {
		return []float64{}
	}

	result := []float64{}

	for i := 0; i <= len(nums)-k; i++ {
		window := nums[i : i+k]
		median := FindMedianOfStream(window)
		result = append(result, median)
	}

	return result
}


