package treebfs

import (
	"reflect"
	"testing"
)

// TestLevelOrderTraversal проверяет обход по уровням
func TestLevelOrderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected [][]int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: [][]int{{3}, {9, 20}, {15, 7}},
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: [][]int{{1}},
		},
		{
			name:     "Пустое дерево",
			root:     nil,
			expected: [][]int{},
		},
		{
			name: "Только левые узлы",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val: 2,
					Left: &TreeNode{Val: 3},
				},
			},
			expected: [][]int{{1}, {2}, {3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LevelOrderTraversal(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("LevelOrderTraversal() = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}

// TestMinDepth проверяет поиск минимальной глубины
func TestMinDepth(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: 2,
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: 1,
		},
		{
			name:     "Пустое дерево",
			root:     nil,
			expected: 0,
		},
		{
			name: "Лист на уровне 2",
			root: &TreeNode{
				Val:  1,
				Left: &TreeNode{Val: 2},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinDepth(tt.root)
			if result != tt.expected {
				t.Errorf("MinDepth() = %d; ожидалось %d",
					result, tt.expected)
			}
		})
	}
}

// TestMaxDepth проверяет поиск максимальной глубины
func TestMaxDepth(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: 3,
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: 1,
		},
		{
			name:     "Пустое дерево",
			root:     nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxDepth(tt.root)
			if result != tt.expected {
				t.Errorf("MaxDepth() = %d; ожидалось %d",
					result, tt.expected)
			}
		})
	}
}

// TestZigzagLevelOrder проверяет зигзагообразный обход
func TestZigzagLevelOrder(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected [][]int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: [][]int{{3}, {20, 9}, {15, 7}},
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: [][]int{{1}},
		},
		{
			name: "Три уровня",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 4},
					Right: &TreeNode{Val: 5},
				},
				Right: &TreeNode{
					Val:   3,
					Left:  &TreeNode{Val: 6},
					Right: &TreeNode{Val: 7},
				},
			},
			expected: [][]int{{1}, {3, 2}, {4, 5, 6, 7}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZigzagLevelOrder(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ZigzagLevelOrder() = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}

// TestLevelOrderSuccessor проверяет поиск следующего узла на уровне
func TestLevelOrderSuccessor(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		key      int
		expected *int // nil если не найден
	}{
		{
			name: "Найти следующий после 9",
			root: &TreeNode{
				Val: 3,
				Left: &TreeNode{Val: 9},
				Right: &TreeNode{
					Val:   20,
					Left:  &TreeNode{Val: 15},
					Right: &TreeNode{Val: 7},
				},
			},
			key:      9,
			expected: intPtr(20),
		},
		{
			name: "Ключ не найден",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
			key:      99,
			expected: nil,
		},
		{
			name: "Последний узел на уровне",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
			key:      3,
			expected: nil, // Нет следующего узла
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LevelOrderSuccessor(tt.root, tt.key)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("LevelOrderSuccessor() = %v; ожидалось nil", result.Val)
				}
			} else {
				if result == nil || result.Val != *tt.expected {
					val := -1
					if result != nil {
						val = result.Val
					}
					t.Errorf("LevelOrderSuccessor() = %d; ожидалось %d", val, *tt.expected)
				}
			}
		})
	}
}

// Вспомогательная функция для создания указателя на int
func intPtr(v int) *int {
	return &v
}

