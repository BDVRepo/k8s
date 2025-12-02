// Package topologicalsort демонстрирует паттерн "Топологическая сортировка"
package topologicalsort

// TopologicalSort выполняет топологическую сортировку графа
// Возвращает порядок вершин или nil, если есть цикл
// Временная сложность: O(V + E)
// Пространственная сложность: O(V + E)
//
// Пример:
//
//	// Граф: 0 -> 1, 0 -> 2, 1 -> 3, 2 -> 3
//	graph := [][]int{
//		0: {1, 2},
//		1: {3},
//		2: {3},
//		3: {},
//	}
//	result := TopologicalSort(4, graph) // Вернёт [0, 1, 2, 3] или [0, 2, 1, 3]
func TopologicalSort(numVertices int, graph [][]int) []int {
	// Вычисляем входящие степени (in-degree) для каждой вершины
	inDegree := make([]int, numVertices)
	for _, neighbors := range graph {
		for _, neighbor := range neighbors {
			inDegree[neighbor]++
		}
	}

	// Находим все вершины с in-degree = 0
	queue := []int{}
	for i := 0; i < numVertices; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	result := []int{}

	// Обрабатываем вершины
	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		// Уменьшаем in-degree всех соседей
		for _, neighbor := range graph[vertex] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Если не все вершины обработаны, значит есть цикл
	if len(result) != numVertices {
		return nil
	}

	return result
}

// CanFinish проверяет, можно ли завершить все курсы с заданными зависимостями
// prerequisites[i] = [a, b] означает, что курс a требует завершения курса b
// Временная сложность: O(V + E)
// Пространственная сложность: O(V + E)
//
// Пример:
//
//	numCourses := 2
//	prerequisites := [][]int{{1, 0}}
//	result := CanFinish(numCourses, prerequisites) // Вернёт true
func CanFinish(numCourses int, prerequisites [][]int) bool {
	// Строим граф
	graph := make([][]int, numCourses)
	inDegree := make([]int, numCourses)

	for _, prereq := range prerequisites {
		course := prereq[0]
		required := prereq[1]
		graph[required] = append(graph[required], course)
		inDegree[course]++
	}

	// Находим курсы без зависимостей
	queue := []int{}
	for i := 0; i < numCourses; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	count := 0

	// Обрабатываем курсы
	for len(queue) > 0 {
		course := queue[0]
		queue = queue[1:]
		count++

		// Уменьшаем in-degree зависимых курсов
		for _, nextCourse := range graph[course] {
			inDegree[nextCourse]--
			if inDegree[nextCourse] == 0 {
				queue = append(queue, nextCourse)
			}
		}
	}

	// Если обработали все курсы, значит можно завершить
	return count == numCourses
}

// FindOrder находит порядок прохождения курсов
// Возвращает порядок или пустой слайс, если невозможно
// Временная сложность: O(V + E)
// Пространственная сложность: O(V + E)
//
// Пример:
//
//	numCourses := 4
//	prerequisites := [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}}
//	result := FindOrder(numCourses, prerequisites)
//	// Вернёт [0, 1, 2, 3] или [0, 2, 1, 3]
func FindOrder(numCourses int, prerequisites [][]int) []int {
	// Строим граф
	graph := make([][]int, numCourses)
	inDegree := make([]int, numCourses)

	for _, prereq := range prerequisites {
		course := prereq[0]
		required := prereq[1]
		graph[required] = append(graph[required], course)
		inDegree[course]++
	}

	// Находим курсы без зависимостей
	queue := []int{}
	for i := 0; i < numCourses; i++ {
		if inDegree[i] == 0 {
			queue = append(queue, i)
		}
	}

	result := []int{}

	// Обрабатываем курсы
	for len(queue) > 0 {
		course := queue[0]
		queue = queue[1:]
		result = append(result, course)

		// Уменьшаем in-degree зависимых курсов
		for _, nextCourse := range graph[course] {
			inDegree[nextCourse]--
			if inDegree[nextCourse] == 0 {
				queue = append(queue, nextCourse)
			}
		}
	}

	// Если не все курсы обработаны, значит есть цикл
	if len(result) != numCourses {
		return []int{}
	}

	return result
}

// FindMinHeightTrees находит корни деревьев с минимальной высотой
// Временная сложность: O(V + E)
// Пространственная сложность: O(V + E)
//
// Пример:
//
//	n := 4
//	edges := [][]int{{1, 0}, {1, 2}, {1, 3}}
//	result := FindMinHeightTrees(n, edges) // Вернёт [1]
func FindMinHeightTrees(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}

	// Строим граф
	graph := make([][]int, n)
	degree := make([]int, n)

	for _, edge := range edges {
		u, v := edge[0], edge[1]
		graph[u] = append(graph[u], v)
		graph[v] = append(graph[v], u)
		degree[u]++
		degree[v]++
	}

	// Находим листья (вершины со степенью 1)
	queue := []int{}
	for i := 0; i < n; i++ {
		if degree[i] == 1 {
			queue = append(queue, i)
		}
	}

	// Удаляем листья слой за слоем
	remaining := n
	for remaining > 2 {
		levelSize := len(queue)
		remaining -= levelSize

		for i := 0; i < levelSize; i++ {
			leaf := queue[0]
			queue = queue[1:]

			// Уменьшаем степень соседей
			for _, neighbor := range graph[leaf] {
				degree[neighbor]--
				if degree[neighbor] == 1 {
					queue = append(queue, neighbor)
				}
			}
		}
	}

	return queue
}

// IsValidTopologicalOrder проверяет, является ли порядок валидной топологической сортировкой
// Временная сложность: O(V + E)
// Пространственная сложность: O(V)
func IsValidTopologicalOrder(numVertices int, graph [][]int, order []int) bool {
	if len(order) != numVertices {
		return false
	}

	// Создаём маппинг позиций
	position := make(map[int]int)
	for i, vertex := range order {
		position[vertex] = i
	}

	// Проверяем, что для каждого рёбра (u, v) u идёт перед v
	for u := 0; u < numVertices; u++ {
		for _, v := range graph[u] {
			if position[u] >= position[v] {
				return false
			}
		}
	}

	return true
}

