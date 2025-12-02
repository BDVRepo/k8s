// Package treebfs демонстрирует паттерн "BFS (обход в ширину) для дерева"
package treebfs

// TreeNode представляет узел бинарного дерева
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// LevelOrderTraversal выполняет обход дерева по уровням
// Возвращает слайс слайсов, где каждый внутренний слайс — это один уровень
// Временная сложность: O(n)
// Пространственная сложность: O(w), где w — максимальная ширина дерева
//
// Пример:
//
//	// Дерево:
//	//     3
//	//    / \
//	//   9   20
//	//      /  \
//	//     15   7
//	root := &TreeNode{Val: 3, Left: &TreeNode{Val: 9},
//		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}}}
//	result := LevelOrderTraversal(root)
//	// Вернёт [[3], [9, 20], [15, 7]]
func LevelOrderTraversal(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	result := [][]int{}
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		levelValues := []int{}

		// Обрабатываем все узлы текущего уровня
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			levelValues = append(levelValues, node.Val)

			// Добавляем дочерние узлы в очередь
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		result = append(result, levelValues)
	}

	return result
}

// MinDepth находит минимальную глубину бинарного дерева
// Временная сложность: O(n)
// Пространственная сложность: O(w)
//
// Пример:
//
//	// Дерево:
//	//     3
//	//    / \
//	//   9   20
//	//      /  \
//	//     15   7
//	root := &TreeNode{Val: 3, Left: &TreeNode{Val: 9},
//		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}}}
//	depth := MinDepth(root) // Вернёт 2
func MinDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	queue := []*TreeNode{root}
	depth := 1

	for len(queue) > 0 {
		levelSize := len(queue)

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			// Если это лист — нашли минимальную глубину
			if node.Left == nil && node.Right == nil {
				return depth
			}

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		depth++
	}

	return depth
}

// MaxDepth находит максимальную глубину бинарного дерева
// Временная сложность: O(n)
// Пространственная сложность: O(w)
func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	queue := []*TreeNode{root}
	depth := 0

	for len(queue) > 0 {
		levelSize := len(queue)
		depth++

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}

	return depth
}

// ZigzagLevelOrder выполняет зигзагообразный обход дерева
// Чётные уровни слева направо, нечётные — справа налево
// Временная сложность: O(n)
// Пространственная сложность: O(w)
//
// Пример:
//
//	// Дерево:
//	//     3
//	//    / \
//	//   9   20
//	//      /  \
//	//     15   7
//	root := &TreeNode{Val: 3, Left: &TreeNode{Val: 9},
//		Right: &TreeNode{Val: 20, Left: &TreeNode{Val: 15}, Right: &TreeNode{Val: 7}}}
//	result := ZigzagLevelOrder(root)
//	// Вернёт [[3], [20, 9], [15, 7]]
func ZigzagLevelOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	result := [][]int{}
	queue := []*TreeNode{root}
	leftToRight := true

	for len(queue) > 0 {
		levelSize := len(queue)
		levelValues := make([]int, levelSize)

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			// Заполняем массив в нужном направлении
			index := i
			if !leftToRight {
				index = levelSize - 1 - i
			}
			levelValues[index] = node.Val

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}

		result = append(result, levelValues)
		leftToRight = !leftToRight
	}

	return result
}

// LevelOrderSuccessor находит следующий узел на том же уровне после заданного
// Временная сложность: O(n)
// Пространственная сложность: O(w)
func LevelOrderSuccessor(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}

	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			// Если нашли ключ и есть следующий узел на этом уровне
			if node.Val == key && i < levelSize-1 {
				return queue[0]
			}

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}

	return nil
}

// ConnectLevelOrderSiblings связывает узлы на одном уровне (для N-арного дерева)
// Используем структуру с полем Next
type NodeWithNext struct {
	Val   int
	Left  *NodeWithNext
	Right *NodeWithNext
	Next  *NodeWithNext
}

// ConnectLevelOrderSiblings связывает узлы на каждом уровне
func ConnectLevelOrderSiblings(root *NodeWithNext) *NodeWithNext {
	if root == nil {
		return nil
	}

	queue := []*NodeWithNext{root}

	for len(queue) > 0 {
		levelSize := len(queue)

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			// Связываем с следующим узлом на том же уровне
			if i < levelSize-1 {
				node.Next = queue[0]
			}

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}

	return root
}

