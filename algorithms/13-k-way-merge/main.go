// Package kwaymerge демонстрирует паттерн "K-way слияние"
package kwaymerge

import (
	"container/heap"
)

// ListNode представляет узел связанного списка
type ListNode struct {
	Val  int
	Next *ListNode
}

// NodeHeap — минимальная куча для узлов связанных списков
type NodeHeap []*ListNode

func (h NodeHeap) Len() int           { return len(h) }
func (h NodeHeap) Less(i, j int) bool { return h[i].Val < h[j].Val }
func (h NodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *NodeHeap) Push(x interface{}) {
	*h = append(*h, x.(*ListNode))
}

func (h *NodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MergeKLists объединяет K отсортированных связанных списков
// Временная сложность: O(n log k), где n — общее количество узлов
// Пространственная сложность: O(k)
//
// Пример:
//
//	list1 := &ListNode{Val: 1, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}}
//	list2 := &ListNode{Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}}
//	list3 := &ListNode{Val: 2, Next: &ListNode{Val: 6}}
//	lists := []*ListNode{list1, list2, list3}
//	result := MergeKLists(lists)
//	// Вернёт: 1 -> 1 -> 2 -> 3 -> 4 -> 4 -> 5 -> 6
func MergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}

	// Создаём min-heap и добавляем первые узлы каждого списка
	h := &NodeHeap{}
	heap.Init(h)

	for _, list := range lists {
		if list != nil {
			heap.Push(h, list)
		}
	}

	// Создаём фиктивный узел для упрощения
	dummy := &ListNode{}
	current := dummy

	// Извлекаем минимум и добавляем следующий узел из того же списка
	for h.Len() > 0 {
		node := heap.Pop(h).(*ListNode)
		current.Next = node
		current = current.Next

		// Добавляем следующий узел из того же списка
		if node.Next != nil {
			heap.Push(h, node.Next)
		}
	}

	return dummy.Next
}

// Element представляет элемент массива с индексами
type Element struct {
	value       int
	arrayIndex  int
	elementIndex int
}

// ElementHeap — минимальная куча для элементов
type ElementHeap []Element

func (h ElementHeap) Len() int           { return len(h) }
func (h ElementHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h ElementHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ElementHeap) Push(x interface{}) {
	*h = append(*h, x.(Element))
}

func (h *ElementHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MergeKSortedArrays объединяет K отсортированных массивов
// Временная сложность: O(n log k)
// Пространственная сложность: O(k)
//
// Пример:
//
//	arrays := [][]int{{1, 4, 5}, {1, 3, 4}, {2, 6}}
//	result := MergeKSortedArrays(arrays)
//	// Вернёт [1, 1, 2, 3, 4, 4, 5, 6]
func MergeKSortedArrays(arrays [][]int) []int {
	if len(arrays) == 0 {
		return []int{}
	}

	h := &ElementHeap{}
	heap.Init(h)

	// Добавляем первые элементы каждого массива
	for i, arr := range arrays {
		if len(arr) > 0 {
			heap.Push(h, Element{
				value:       arr[0],
				arrayIndex:  i,
				elementIndex: 0,
			})
		}
	}

	result := []int{}

	// Извлекаем минимум и добавляем следующий элемент из того же массива
	for h.Len() > 0 {
		elem := heap.Pop(h).(Element)
		result = append(result, elem.value)

		// Добавляем следующий элемент из того же массива
		if elem.elementIndex+1 < len(arrays[elem.arrayIndex]) {
			heap.Push(h, Element{
				value:       arrays[elem.arrayIndex][elem.elementIndex+1],
				arrayIndex:  elem.arrayIndex,
				elementIndex: elem.elementIndex + 1,
			})
		}
	}

	return result
}

// Cell представляет ячейку матрицы с координатами
type Cell struct {
	value int
	row   int
	col   int
}

// CellHeap — минимальная куча для ячеек
type CellHeap []Cell

func (h CellHeap) Len() int           { return len(h) }
func (h CellHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h CellHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *CellHeap) Push(x interface{}) {
	*h = append(*h, x.(Cell))
}

func (h *CellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// FindKthSmallestInMatrix находит K-й наименьший элемент в отсортированной матрице
// Временная сложность: O(k log k)
// Пространственная сложность: O(k)
//
// Пример:
//
//	matrix := [][]int{
//		{1, 5, 9},
//		{10, 11, 13},
//		{12, 13, 15},
//	}
//	k := 8
//	result := FindKthSmallestInMatrix(matrix, k) // Вернёт 13
func FindKthSmallestInMatrix(matrix [][]int, k int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return -1
	}

	h := &CellHeap{}
	heap.Init(h)

	// Добавляем первые элементы каждой строки
	for i := 0; i < len(matrix) && i < k; i++ {
		heap.Push(h, Cell{
			value: matrix[i][0],
			row:   i,
			col:   0,
		})
	}

	count := 0
	result := -1

	// Извлекаем k элементов
	for h.Len() > 0 && count < k {
		cell := heap.Pop(h).(Cell)
		result = cell.value
		count++

		// Добавляем следующий элемент из той же строки
		if cell.col+1 < len(matrix[cell.row]) {
			heap.Push(h, Cell{
				value: matrix[cell.row][cell.col+1],
				row:   cell.row,
				col:   cell.col + 1,
			})
		}
	}

	return result
}

// CreateList создаёт связанный список из слайса
func CreateList(values []int) *ListNode {
	if len(values) == 0 {
		return nil
	}

	head := &ListNode{Val: values[0]}
	current := head

	for i := 1; i < len(values); i++ {
		current.Next = &ListNode{Val: values[i]}
		current = current.Next
	}

	return head
}

// ListToSlice преобразует связанный список в слайс
func ListToSlice(head *ListNode) []int {
	result := []int{}
	current := head

	for current != nil {
		result = append(result, current.Val)
		current = current.Next
	}

	return result
}

