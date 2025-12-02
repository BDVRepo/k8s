package reverselinkedlist

import (
	"reflect"
	"testing"
)

// TestReverseLinkedList проверяет разворот списка
func TestReverseLinkedList(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected []int
	}{
		{
			name:     "Обычный случай",
			values:   []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
		{
			name:     "Два элемента",
			values:   []int{1, 2},
			expected: []int{2, 1},
		},
		{
			name:     "Один элемент",
			values:   []int{1},
			expected: []int{1},
		},
		{
			name:     "Пустой список",
			values:   []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := CreateList(tt.values)
			reversed := ReverseLinkedList(head)
			result := ListToSlice(reversed)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReverseLinkedList(%v) = %v; ожидалось %v",
					tt.values, result, tt.expected)
			}
		})
	}
}

// TestReverseLinkedListRecursive проверяет рекурсивный разворот
func TestReverseLinkedListRecursive(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected []int
	}{
		{
			name:     "Обычный случай",
			values:   []int{1, 2, 3, 4},
			expected: []int{4, 3, 2, 1},
		},
		{
			name:     "Один элемент",
			values:   []int{5},
			expected: []int{5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := CreateList(tt.values)
			reversed := ReverseLinkedListRecursive(head)
			result := ListToSlice(reversed)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ReverseLinkedListRecursive(%v) = %v; ожидалось %v",
					tt.values, result, tt.expected)
			}
		})
	}
}

// TestReverseBetween проверяет частичный разворот
func TestReverseBetween(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		left     int
		right    int
		expected []int
	}{
		{
			name:     "Разворот середины",
			values:   []int{1, 2, 3, 4, 5},
			left:     2,
			right:    4,
			expected: []int{1, 4, 3, 2, 5},
		},
		{
			name:     "Разворот с начала",
			values:   []int{1, 2, 3, 4, 5},
			left:     1,
			right:    3,
			expected: []int{3, 2, 1, 4, 5},
		},
		{
			name:     "Разворот до конца",
			values:   []int{1, 2, 3, 4, 5},
			left:     3,
			right:    5,
			expected: []int{1, 2, 5, 4, 3},
		},
		{
			name:     "Разворот одного элемента",
			values:   []int{1, 2, 3},
			left:     2,
			right:    2,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := CreateList(tt.values)
			result := ReverseBetween(head, tt.left, tt.right)
			values := ListToSlice(result)

			if !reflect.DeepEqual(values, tt.expected) {
				t.Errorf("ReverseBetween(%v, %d, %d) = %v; ожидалось %v",
					tt.values, tt.left, tt.right, values, tt.expected)
			}
		})
	}
}

// TestReverseKGroup проверяет разворот по группам
func TestReverseKGroup(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		k        int
		expected []int
	}{
		{
			name:     "Группы по 2",
			values:   []int{1, 2, 3, 4, 5},
			k:        2,
			expected: []int{2, 1, 4, 3, 5},
		},
		{
			name:     "Группы по 3",
			values:   []int{1, 2, 3, 4, 5},
			k:        3,
			expected: []int{3, 2, 1, 4, 5},
		},
		{
			name:     "k=1 (без изменений)",
			values:   []int{1, 2, 3},
			k:        1,
			expected: []int{1, 2, 3},
		},
		{
			name:     "k равно длине списка",
			values:   []int{1, 2, 3},
			k:        3,
			expected: []int{3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := CreateList(tt.values)
			result := ReverseKGroup(head, tt.k)
			values := ListToSlice(result)

			if !reflect.DeepEqual(values, tt.expected) {
				t.Errorf("ReverseKGroup(%v, %d) = %v; ожидалось %v",
					tt.values, tt.k, values, tt.expected)
			}
		})
	}
}

// TestIsPalindrome проверяет определение палиндрома
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected bool
	}{
		{
			name:     "Палиндром чётной длины",
			values:   []int{1, 2, 2, 1},
			expected: true,
		},
		{
			name:     "Палиндром нечётной длины",
			values:   []int{1, 2, 3, 2, 1},
			expected: true,
		},
		{
			name:     "Не палиндром",
			values:   []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "Один элемент",
			values:   []int{1},
			expected: true,
		},
		{
			name:     "Два одинаковых элемента",
			values:   []int{1, 1},
			expected: true,
		},
		{
			name:     "Два разных элемента",
			values:   []int{1, 2},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := CreateList(tt.values)
			result := IsPalindrome(head)

			if result != tt.expected {
				t.Errorf("IsPalindrome(%v) = %v; ожидалось %v",
					tt.values, result, tt.expected)
			}
		})
	}
}


