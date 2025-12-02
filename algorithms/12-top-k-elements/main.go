// Package topkelements демонстрирует паттерн "Топ K-элементов"
package topkelements

import (
	"container/heap"
)

// IntHeap — минимальная куча для хранения K наибольших элементов
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// FindKLargest находит K наибольших элементов
// Временная сложность: O(n log k)
// Пространственная сложность: O(k)
//
// Пример:
//
//	nums := []int{3, 1, 5, 12, 2, 11}
//	k := 3
//	result := FindKLargest(nums, k) // Вернёт [12, 11, 5]
func FindKLargest(nums []int, k int) []int {
	if k <= 0 || len(nums) == 0 {
		return []int{}
	}

	// Используем min-heap размера k
	h := &IntHeap{}
	heap.Init(h)

	for _, num := range nums {
		if h.Len() < k {
			heap.Push(h, num)
		} else if num > (*h)[0] {
			// Если текущий элемент больше минимума в куче, заменяем
			heap.Pop(h)
			heap.Push(h, num)
		}
	}

	// Извлекаем элементы из кучи
	result := make([]int, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(int)
	}

	return result
}

// MaxIntHeap — максимальная куча для хранения K наименьших элементов
type MaxIntHeap []int

func (h MaxIntHeap) Len() int           { return len(h) }
func (h MaxIntHeap) Less(i, j int) bool { return h[i] > h[j] } // Инвертируем для max-heap
func (h MaxIntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxIntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxIntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// FindKSmallest находит K наименьших элементов
// Временная сложность: O(n log k)
// Пространственная сложность: O(k)
func FindKSmallest(nums []int, k int) []int {
	if k <= 0 || len(nums) == 0 {
		return []int{}
	}

	// Используем max-heap размера k
	maxHeap := &MaxIntHeap{}
	heap.Init(maxHeap)

	for _, num := range nums {
		if maxHeap.Len() < k {
			heap.Push(maxHeap, num)
		} else if num < (*maxHeap)[0] {
			// Если текущий элемент меньше максимума в куче, заменяем
			heap.Pop(maxHeap)
			heap.Push(maxHeap, num)
		}
	}

	// Извлекаем элементы из кучи
	result := make([]int, maxHeap.Len())
	for i := maxHeap.Len() - 1; i >= 0; i-- {
		result[i] = heap.Pop(maxHeap).(int)
	}

	return result
}

// FindKthLargest находит K-й наибольший элемент
// Временная сложность: O(n log k)
// Пространственная сложность: O(k)
//
// Пример:
//
//	nums := []int{3, 2, 1, 5, 6, 4}
//	k := 2
//	result := FindKthLargest(nums, k) // Вернёт 5
func FindKthLargest(nums []int, k int) int {
	if k <= 0 || len(nums) == 0 || k > len(nums) {
		return -1
	}

	h := &IntHeap{}
	heap.Init(h)

	for _, num := range nums {
		if h.Len() < k {
			heap.Push(h, num)
		} else if num > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, num)
		}
	}

	return (*h)[0] // Минимум в min-heap размера k — это k-й наибольший
}

// Pair представляет пару (частота, элемент)
type Pair struct {
	freq int
	num  int
}

// PairHeap — минимальная куча для пар
type PairHeap []Pair

func (h PairHeap) Len() int           { return len(h) }
func (h PairHeap) Less(i, j int) bool { return h[i].freq < h[j].freq }
func (h PairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *PairHeap) Push(x interface{}) {
	*h = append(*h, x.(Pair))
}

func (h *PairHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TopKFrequent находит K наиболее частых элементов
// Временная сложность: O(n log k)
// Пространственная сложность: O(n)
//
// Пример:
//
//	nums := []int{1, 1, 1, 2, 2, 3}
//	k := 2
//	result := TopKFrequent(nums, k) // Вернёт [1, 2]
func TopKFrequent(nums []int, k int) []int {
	if k <= 0 || len(nums) == 0 {
		return []int{}
	}

	// Подсчитываем частоту
	freq := make(map[int]int)
	for _, num := range nums {
		freq[num]++
	}

	// Создаём кучу пар (частота, элемент)
	pairHeap := &PairHeap{}
	heap.Init(pairHeap)

	for num, count := range freq {
		if pairHeap.Len() < k {
			heap.Push(pairHeap, Pair{freq: count, num: num})
		} else if count > (*pairHeap)[0].freq {
			heap.Pop(pairHeap)
			heap.Push(pairHeap, Pair{freq: count, num: num})
		}
	}

	result := make([]int, pairHeap.Len())
	for i := 0; i < pairHeap.Len(); i++ {
		result[i] = (*pairHeap)[i].num
	}

	return result
}

// PointHeap — максимальная куча для точек (по расстоянию)
type PointHeap [][]int

func (h PointHeap) Len() int { return len(h) }
func (h PointHeap) Less(i, j int) bool {
	distI := h[i][0]*h[i][0] + h[i][1]*h[i][1]
	distJ := h[j][0]*h[j][0] + h[j][1]*h[j][1]
	return distI > distJ // Max-heap по расстоянию
}
func (h PointHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PointHeap) Push(x interface{}) {
	*h = append(*h, x.([]int))
}

func (h *PointHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// KClosestPoints находит K ближайших точек к началу координат
// Временная сложность: O(n log k)
// Пространственная сложность: O(k)
//
// Пример:
//
//	points := [][]int{{1, 3}, {-2, 2}, {5, 8}}
//	k := 1
//	result := KClosestPoints(points, k) // Вернёт [[-2, 2]]
func KClosestPoints(points [][]int, k int) [][]int {
	if k <= 0 || len(points) == 0 {
		return [][]int{}
	}

	// Вычисляем расстояние до начала координат
	distance := func(point []int) int {
		return point[0]*point[0] + point[1]*point[1]
	}

	// Используем max-heap для хранения K ближайших
	pointHeap := &PointHeap{}
	heap.Init(pointHeap)

	for _, point := range points {
		if pointHeap.Len() < k {
			heap.Push(pointHeap, point)
		} else {
			// Сравниваем расстояния
			currDist := distance(point)
			farthestDist := distance((*pointHeap)[0])

			if currDist < farthestDist {
				heap.Pop(pointHeap)
				heap.Push(pointHeap, point)
			}
		}
	}

	result := make([][]int, pointHeap.Len())
	for i := 0; i < pointHeap.Len(); i++ {
		result[i] = (*pointHeap)[i]
	}

	return result
}

// WordPair представляет пару (слово, частота)
type WordPair struct {
	word string
	freq int
}

// WordHeap — минимальная куча для пар слов
type WordHeap []WordPair

func (h WordHeap) Len() int           { return len(h) }
func (h WordHeap) Less(i, j int) bool { return h[i].freq < h[j].freq }
func (h WordHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *WordHeap) Push(x interface{}) {
	*h = append(*h, x.(WordPair))
}

func (h *WordHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// TopKFrequentWords находит K наиболее частых слов
// Временная сложность: O(n log k)
// Пространственная сложность: O(n)
func TopKFrequentWords(words []string, k int) []string {
	if k <= 0 || len(words) == 0 {
		return []string{}
	}

	// Подсчитываем частоту
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}

	// Создаём кучу пар (частота, слово)
	wordHeap := &WordHeap{}
	heap.Init(wordHeap)

	for word, count := range freq {
		if wordHeap.Len() < k {
			heap.Push(wordHeap, WordPair{word: word, freq: count})
		} else if count > (*wordHeap)[0].freq {
			heap.Pop(wordHeap)
			heap.Push(wordHeap, WordPair{word: word, freq: count})
		}
	}

	result := make([]string, wordHeap.Len())
	for i := 0; i < wordHeap.Len(); i++ {
		result[i] = (*wordHeap)[i].word
	}

	return result
}

