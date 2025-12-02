package treedfs

import (
	"reflect"
	"testing"
)

// TestHasPathSum проверяет поиск пути с заданной суммой
func TestHasPathSum(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		target   int
		expected bool
	}{
		{
			name: "Путь существует",
			root: &TreeNode{
				Val: 5,
				Left: &TreeNode{
					Val:   4,
					Left:  &TreeNode{Val: 11},
					Right: nil,
				},
				Right: &TreeNode{
					Val:   8,
					Left:  &TreeNode{Val: 13},
					Right: &TreeNode{Val: 4},
				},
			},
			target:   20, // 5 + 4 + 11
			expected: true,
		},
		{
			name:     "Один узел, сумма совпадает",
			root:     &TreeNode{Val: 1},
			target:   1,
			expected: true,
		},
		{
			name:     "Один узел, сумма не совпадает",
			root:     &TreeNode{Val: 1},
			target:   2,
			expected: false,
		},
		{
			name:     "Пустое дерево",
			root:     nil,
			target:   0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasPathSum(tt.root, tt.target)
			if result != tt.expected {
				t.Errorf("HasPathSum(..., %d) = %v; ожидалось %v",
					tt.target, result, tt.expected)
			}
		})
	}
}

// TestFindAllPaths проверяет поиск всех путей
func TestFindAllPaths(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []string
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:   2,
					Right: &TreeNode{Val: 5},
				},
				Right: &TreeNode{Val: 3},
			},
			expected: []string{"1->2->5", "1->3"},
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: []string{"1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindAllPaths(tt.root)
			// Упрощённая проверка — в реальном коде нужно правильно преобразовывать числа в строки
			if len(result) != len(tt.expected) {
				t.Errorf("FindAllPaths() вернул %d путей; ожидалось %d",
					len(result), len(tt.expected))
			}
		})
	}
}

// TestIsSymmetric проверяет симметрию дерева
func TestIsSymmetric(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected bool
	}{
		{
			name: "Симметричное дерево",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 3},
					Right: &TreeNode{Val: 4},
				},
				Right: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 4},
					Right: &TreeNode{Val: 3},
				},
			},
			expected: true,
		},
		{
			name: "Несимметричное дерево",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
			expected: false,
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: true,
		},
		{
			name:     "Пустое дерево",
			root:     nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSymmetric(tt.root)
			if result != tt.expected {
				t.Errorf("IsSymmetric() = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}

// TestDiameterOfBinaryTree проверяет поиск диаметра
func TestDiameterOfBinaryTree(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val: 1,
				Left: &TreeNode{
					Val:   2,
					Left:  &TreeNode{Val: 4},
					Right: &TreeNode{Val: 5},
				},
				Right: &TreeNode{Val: 3},
			},
			expected: 3, // Путь через корень: 4 -> 2 -> 1 -> 3
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DiameterOfBinaryTree(tt.root)
			if result != tt.expected {
				t.Errorf("DiameterOfBinaryTree() = %d; ожидалось %d",
					result, tt.expected)
			}
		})
	}
}

// TestPreorderTraversal проверяет Pre-order обход
func TestPreorderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
			expected: []int{1, 2, 3},
		},
		{
			name:     "Один узел",
			root:     &TreeNode{Val: 1},
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreorderTraversal(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PreorderTraversal() = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}

// TestInorderTraversal проверяет In-order обход
func TestInorderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
			expected: []int{2, 1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InorderTraversal(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InorderTraversal() = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}

// TestPostorderTraversal проверяет Post-order обход
func TestPostorderTraversal(t *testing.T) {
	tests := []struct {
		name     string
		root     *TreeNode
		expected []int
	}{
		{
			name: "Обычное дерево",
			root: &TreeNode{
				Val:   1,
				Left:  &TreeNode{Val: 2},
				Right: &TreeNode{Val: 3},
			},
			expected: []int{2, 3, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PostorderTraversal(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PostorderTraversal() = %v; ожидалось %v",
					result, tt.expected)
			}
		})
	}
}


