package kwaymerge

import (
	"reflect"
	"testing"
)

// TestMergeKLists проверяет объединение K списков
func TestMergeKLists(t *testing.T) {
	tests := []struct {
		name     string
		lists    []*ListNode
		expected []int
	}{
		{
			name: "Три списка",
			lists: []*ListNode{
				CreateList([]int{1, 4, 5}),
				CreateList([]int{1, 3, 4}),
				CreateList([]int{2, 6}),
			},
			expected: []int{1, 1, 2, 3, 4, 4, 5, 6},
		},
		{
			name: "Один список",
			lists: []*ListNode{
				CreateList([]int{1, 2, 3}),
			},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Пустой список",
			lists:    []*ListNode{},
			expected: []int{},
		},
		{
			name: "С пустыми списками",
			lists: []*ListNode{
				CreateList([]int{1, 2}),
				nil,
				CreateList([]int{3, 4}),
			},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeKLists(tt.lists)
			values := ListToSlice(result)

			if !reflect.DeepEqual(values, tt.expected) {
				t.Errorf("MergeKLists(...) = %v; ожидалось %v",
					values, tt.expected)
			}
		})
	}
}

// TestMergeKSortedArrays проверяет объединение K массивов
func TestMergeKSortedArrays(t *testing.T) {
	tests := []struct {
		name     string
		arrays   [][]int
		expected []int
	}{
		{
			name: "Три массива",
			arrays: [][]int{
				{1, 4, 5},
				{1, 3, 4},
				{2, 6},
			},
			expected: []int{1, 1, 2, 3, 4, 4, 5, 6},
		},
		{
			name: "Один массив",
			arrays: [][]int{
				{1, 2, 3},
			},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Пустой массив",
			arrays:   [][]int{},
			expected: []int{},
		},
		{
			name: "С пустыми массивами",
			arrays: [][]int{
				{1, 2},
				{},
				{3, 4},
			},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeKSortedArrays(tt.arrays)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeKSortedArrays(...) = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}

// TestFindKthSmallestInMatrix проверяет поиск K-го наименьшего элемента
func TestFindKthSmallestInMatrix(t *testing.T) {
	tests := []struct {
		name     string
		matrix   [][]int
		k        int
		expected int
	}{
		{
			name: "Обычный случай",
			matrix: [][]int{
				{1, 5, 9},
				{10, 11, 13},
				{12, 13, 15},
			},
			k:        8,
			expected: 13,
		},
		{
			name: "k=1",
			matrix: [][]int{
				{1, 5, 9},
				{10, 11, 13},
			},
			k:        1,
			expected: 1,
		},
		{
			name: "k равно количеству элементов",
			matrix: [][]int{
				{1, 2},
				{3, 4},
			},
			k:        4,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindKthSmallestInMatrix(tt.matrix, tt.k)
			if result != tt.expected {
				t.Errorf("FindKthSmallestInMatrix(..., %d) = %d; ожидалось %d",
					tt.k, result, tt.expected)
			}
		})
	}
}

