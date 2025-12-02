package fastslowpointers

import (
	"testing"
)

// TestHasCycle проверяет обнаружение цикла
func TestHasCycle(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		cyclePos  int // Позиция, к которой присоединяется последний узел (-1 = нет цикла)
		expected  bool
	}{
		{
			name:     "Список с циклом",
			values:   []int{1, 2, 3, 4},
			cyclePos: 1, // Последний узел указывает на второй
			expected: true,
		},
		{
			name:     "Список без цикла",
			values:   []int{1, 2, 3, 4},
			cyclePos: -1,
			expected: false,
		},
		{
			name:     "Один элемент без цикла",
			values:   []int{1},
			cyclePos: -1,
			expected: false,
		},
		{
			name:     "Один элемент с циклом на себя",
			values:   []int{1},
			cyclePos: 0,
			expected: true,
		},
		{
			name:     "Пустой список",
			values:   []int{},
			cyclePos: -1,
			expected: false,
		},
		{
			name:     "Два элемента с циклом",
			values:   []int{1, 2},
			cyclePos: 0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createListWithCycle(tt.values, tt.cyclePos)
			result := HasCycle(head)
			if result != tt.expected {
				t.Errorf("HasCycle() = %v; ожидалось %v", result, tt.expected)
			}
		})
	}
}

// TestFindCycleStart проверяет поиск начала цикла
func TestFindCycleStart(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		cyclePos int
	}{
		{
			name:     "Цикл со второго элемента",
			values:   []int{1, 2, 3, 4, 5},
			cyclePos: 1,
		},
		{
			name:     "Цикл с первого элемента",
			values:   []int{1, 2, 3},
			cyclePos: 0,
		},
		{
			name:     "Цикл с последнего элемента на себя",
			values:   []int{1, 2, 3, 4},
			cyclePos: 3,
		},
		{
			name:     "Нет цикла",
			values:   []int{1, 2, 3},
			cyclePos: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, cycleStart := createListWithCycleAndGetStart(tt.values, tt.cyclePos)
			result := FindCycleStart(head)

			if tt.cyclePos == -1 {
				if result != nil {
					t.Errorf("FindCycleStart() = %v; ожидалось nil", result.Val)
				}
			} else {
				if result != cycleStart {
					t.Errorf("FindCycleStart() нашёл неверное начало цикла")
				}
			}
		})
	}
}

// TestFindMiddle проверяет поиск середины списка
func TestFindMiddle(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected int
	}{
		{
			name:     "Нечётное количество элементов",
			values:   []int{1, 2, 3, 4, 5},
			expected: 3,
		},
		{
			name:     "Чётное количество элементов",
			values:   []int{1, 2, 3, 4, 5, 6},
			expected: 4, // Второй из двух средних
		},
		{
			name:     "Один элемент",
			values:   []int{1},
			expected: 1,
		},
		{
			name:     "Два элемента",
			values:   []int{1, 2},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := CreateLinkedList(tt.values)
			result := FindMiddle(head)

			if result == nil {
				t.Error("FindMiddle() вернул nil")
				return
			}

			if result.Val != tt.expected {
				t.Errorf("FindMiddle() = %d; ожидалось %d", result.Val, tt.expected)
			}
		})
	}
}

// TestFindMiddleEmpty проверяет поиск середины для пустого списка
func TestFindMiddleEmpty(t *testing.T) {
	result := FindMiddle(nil)
	if result != nil {
		t.Errorf("FindMiddle(nil) = %v; ожидалось nil", result)
	}
}

// TestIsHappyNumber проверяет определение счастливого числа
func TestIsHappyNumber(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected bool
	}{
		{
			name:     "Счастливое число 23",
			n:        23,
			expected: true, // 2²+3²=13 → 1²+3²=10 → 1²+0²=1
		},
		{
			name:     "Счастливое число 1",
			n:        1,
			expected: true,
		},
		{
			name:     "Счастливое число 19",
			n:        19,
			expected: true,
		},
		{
			name:     "Несчастливое число 2",
			n:        2,
			expected: false,
		},
		{
			name:     "Несчастливое число 12",
			n:        12,
			expected: false,
		},
		{
			name:     "Счастливое число 7",
			n:        7,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHappyNumber(tt.n)
			if result != tt.expected {
				t.Errorf("IsHappyNumber(%d) = %v; ожидалось %v",
					tt.n, result, tt.expected)
			}
		})
	}
}

// TestGetCycleLength проверяет вычисление длины цикла
func TestGetCycleLength(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		cyclePos int
		expected int
	}{
		{
			name:     "Цикл длиной 3",
			values:   []int{1, 2, 3, 4, 5},
			cyclePos: 2, // Последний → третий, цикл: 3→4→5→3
			expected: 3,
		},
		{
			name:     "Цикл на весь список",
			values:   []int{1, 2, 3},
			cyclePos: 0,
			expected: 3,
		},
		{
			name:     "Нет цикла",
			values:   []int{1, 2, 3},
			cyclePos: -1,
			expected: 0,
		},
		{
			name:     "Цикл длиной 1 (на себя)",
			values:   []int{1, 2, 3},
			cyclePos: 2,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head := createListWithCycle(tt.values, tt.cyclePos)
			result := GetCycleLength(head)
			if result != tt.expected {
				t.Errorf("GetCycleLength() = %d; ожидалось %d", result, tt.expected)
			}
		})
	}
}

// Вспомогательные функции для создания тестовых списков

// createListWithCycle создаёт связанный список с циклом
func createListWithCycle(values []int, cyclePos int) *ListNode {
	if len(values) == 0 {
		return nil
	}

	head := &ListNode{Val: values[0]}
	current := head
	var cycleNode *ListNode

	if cyclePos == 0 {
		cycleNode = head
	}

	for i := 1; i < len(values); i++ {
		current.Next = &ListNode{Val: values[i]}
		current = current.Next
		if i == cyclePos {
			cycleNode = current
		}
	}

	// Создаём цикл, если cyclePos >= 0
	if cyclePos >= 0 && cycleNode != nil {
		current.Next = cycleNode
	}

	return head
}

// createListWithCycleAndGetStart создаёт список и возвращает начало цикла
func createListWithCycleAndGetStart(values []int, cyclePos int) (*ListNode, *ListNode) {
	if len(values) == 0 {
		return nil, nil
	}

	head := &ListNode{Val: values[0]}
	current := head
	var cycleNode *ListNode

	if cyclePos == 0 {
		cycleNode = head
	}

	for i := 1; i < len(values); i++ {
		current.Next = &ListNode{Val: values[i]}
		current = current.Next
		if i == cyclePos {
			cycleNode = current
		}
	}

	if cyclePos >= 0 && cycleNode != nil {
		current.Next = cycleNode
	}

	return head, cycleNode
}

