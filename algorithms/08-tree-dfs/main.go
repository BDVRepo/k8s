// Package treedfs демонстрирует паттерн "DFS (обход в глубину) для дерева"
package treedfs

import "strconv"

// TreeNode представляет узел бинарного дерева
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// HasPathSum проверяет, существует ли путь от корня до листа с заданной суммой
// Временная сложность: O(n)
// Пространственная сложность: O(h), где h — высота дерева
//
// Пример:
//
//	// Дерево:
//	//     5
//	//    / \
//	//   4   8
//	//  /   / \
//	// 11  13  4
//	root := &TreeNode{Val: 5, Left: &TreeNode{Val: 4, Left: &TreeNode{Val: 11}},
//		Right: &TreeNode{Val: 8, Left: &TreeNode{Val: 13}, Right: &TreeNode{Val: 4}}}
//	result := HasPathSum(root, 20) // Вернёт true (5 -> 4 -> 11)
func HasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}

	// Если это лист и сумма совпадает
	if root.Left == nil && root.Right == nil {
		return root.Val == targetSum
	}

	// Рекурсивно проверяем левое и правое поддеревья
	remainingSum := targetSum - root.Val
	return HasPathSum(root.Left, remainingSum) || HasPathSum(root.Right, remainingSum)
}

// FindAllPaths находит все пути от корня до листьев
// Временная сложность: O(n)
// Пространственная сложность: O(h * n), где n — количество листьев
//
// Пример:
//
//	// Дерево:
//	//     1
//	//    / \
//	//   2   3
//	//    \
//	//     5
//	root := &TreeNode{Val: 1, Left: &TreeNode{Val: 2, Right: &TreeNode{Val: 5}},
//		Right: &TreeNode{Val: 3}}
//	result := FindAllPaths(root)
//	// Вернёт ["1->2->5", "1->3"]
func FindAllPaths(root *TreeNode) []string {
	if root == nil {
		return []string{}
	}

	result := []string{}
	currentPath := []int{}

	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}

		// Добавляем текущий узел в путь
		currentPath = append(currentPath, node.Val)

		// Если это лист, сохраняем путь
		if node.Left == nil && node.Right == nil {
			pathStr := pathToString(currentPath)
			result = append(result, pathStr)
		} else {
			// Рекурсивно обходим дочерние узлы
			dfs(node.Left)
			dfs(node.Right)
		}

		// Убираем текущий узел из пути (backtracking)
		currentPath = currentPath[:len(currentPath)-1]
	}

	dfs(root)
	return result
}

// pathToString преобразует слайс чисел в строку пути
func pathToString(path []int) string {
	if len(path) == 0 {
		return ""
	}

	result := ""
	for i, val := range path {
		if i > 0 {
			result += "->"
		}
		result += strconv.Itoa(val)
	}

	return result
}

// SumOfPathNumbers вычисляет сумму всех чисел, образованных путями от корня до листьев
// Каждый путь представляет число (например, 1->2->3 = 123)
// Временная сложность: O(n)
// Пространственная сложность: O(h)
func SumOfPathNumbers(root *TreeNode) int {
	return sumOfPathNumbersHelper(root, 0)
}

func sumOfPathNumbersHelper(node *TreeNode, pathSum int) int {
	if node == nil {
		return 0
	}

	// Обновляем сумму пути
	pathSum = pathSum*10 + node.Val

	// Если это лист, возвращаем сумму пути
	if node.Left == nil && node.Right == nil {
		return pathSum
	}

	// Рекурсивно суммируем пути из левого и правого поддеревьев
	return sumOfPathNumbersHelper(node.Left, pathSum) +
		sumOfPathNumbersHelper(node.Right, pathSum)
}

// IsSymmetric проверяет, является ли дерево симметричным
// Временная сложность: O(n)
// Пространственная сложность: O(h)
//
// Пример:
//
//	// Дерево:
//	//     1
//	//    / \
//	//   2   2
//	//  / \ / \
//	// 3  4 4  3
//	root := &TreeNode{Val: 1,
//		Left:  &TreeNode{Val: 2, Left: &TreeNode{Val: 3}, Right: &TreeNode{Val: 4}},
//		Right: &TreeNode{Val: 2, Left: &TreeNode{Val: 4}, Right: &TreeNode{Val: 3}}}
//	result := IsSymmetric(root) // Вернёт true
func IsSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return isSymmetricHelper(root.Left, root.Right)
}

func isSymmetricHelper(left *TreeNode, right *TreeNode) bool {
	// Оба nil — симметрично
	if left == nil && right == nil {
		return true
	}

	// Один nil — не симметрично
	if left == nil || right == nil {
		return false
	}

	// Значения должны совпадать
	if left.Val != right.Val {
		return false
	}

	// Рекурсивно проверяем зеркальные поддеревья
	return isSymmetricHelper(left.Left, right.Right) &&
		isSymmetricHelper(left.Right, right.Left)
}

// DiameterOfBinaryTree находит диаметр бинарного дерева
// Диаметр — это количество рёбер в самом длинном пути между любыми двумя узлами
// Временная сложность: O(n)
// Пространственная сложность: O(h)
func DiameterOfBinaryTree(root *TreeNode) int {
	maxDiameter := 0

	var maxDepth func(node *TreeNode) int
	maxDepth = func(node *TreeNode) int {
		if node == nil {
			return 0
		}

		leftDepth := maxDepth(node.Left)
		rightDepth := maxDepth(node.Right)

		// Диаметр через текущий узел
		currentDiameter := leftDepth + rightDepth
		if currentDiameter > maxDiameter {
			maxDiameter = currentDiameter
		}

		// Возвращаем максимальную глубину через текущий узел
		if leftDepth > rightDepth {
			return leftDepth + 1
		}
		return rightDepth + 1
	}

	maxDepth(root)
	return maxDiameter
}

// MaxPathSum находит максимальную сумму пути в бинарном дереве
// Путь может начинаться и заканчиваться в любом узле
// Временная сложность: O(n)
// Пространственная сложность: O(h)
func MaxPathSum(root *TreeNode) int {
	maxSum := root.Val // Инициализируем значением корня

	var maxGain func(node *TreeNode) int
	maxGain = func(node *TreeNode) int {
		if node == nil {
			return 0
		}

		// Максимальная выгода от левого и правого поддеревьев
		leftGain := maxGain(node.Left)
		rightGain := maxGain(node.Right)

		// Если выгода отрицательная, не берём её
		if leftGain < 0 {
			leftGain = 0
		}
		if rightGain < 0 {
			rightGain = 0
		}

		// Сумма пути через текущий узел
		currentPathSum := node.Val + leftGain + rightGain
		if currentPathSum > maxSum {
			maxSum = currentPathSum
		}

		// Возвращаем максимальную выгоду через текущий узел
		return node.Val + max(leftGain, rightGain)
	}

	maxGain(root)
	return maxSum
}

// Вспомогательная функция max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// PreorderTraversal выполняет обход дерева в порядке Pre-order
func PreorderTraversal(root *TreeNode) []int {
	result := []int{}

	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}

		result = append(result, node.Val) // Корень
		dfs(node.Left)                    // Левое
		dfs(node.Right)                   // Правое
	}

	dfs(root)
	return result
}

// InorderTraversal выполняет обход дерева в порядке In-order
func InorderTraversal(root *TreeNode) []int {
	result := []int{}

	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}

		dfs(node.Left)                    // Левое
		result = append(result, node.Val) // Корень
		dfs(node.Right)                   // Правое
	}

	dfs(root)
	return result
}

// PostorderTraversal выполняет обход дерева в порядке Post-order
func PostorderTraversal(root *TreeNode) []int {
	result := []int{}

	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}

		dfs(node.Left)                    // Левое
		dfs(node.Right)                   // Правое
		result = append(result, node.Val) // Корень
	}

	dfs(root)
	return result
}

