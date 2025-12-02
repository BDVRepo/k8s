package topologicalsort

import (
	"reflect"
	"testing"
)

// TestTopologicalSort проверяет топологическую сортировку
func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name        string
		numVertices int
		graph       [][]int
		hasCycle    bool // Ожидается ли цикл
		validLength int  // Ожидаемая длина результата
	}{
		{
			name:        "Обычный граф",
			numVertices: 4,
			graph: [][]int{
				0: {1, 2},
				1: {3},
				2: {3},
				3: {},
			},
			hasCycle:    false,
			validLength: 4,
		},
		{
			name:        "Граф с циклом",
			numVertices: 3,
			graph: [][]int{
				0: {1},
				1: {2},
				2: {0}, // Цикл
			},
			hasCycle:    true,
			validLength: 0,
		},
		{
			name:        "Линейный граф",
			numVertices: 3,
			graph: [][]int{
				0: {1},
				1: {2},
				2: {},
			},
			hasCycle:    false,
			validLength: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TopologicalSort(tt.numVertices, tt.graph)

			if tt.hasCycle {
				if result != nil {
					t.Errorf("TopologicalSort() вернул результат для графа с циклом: %v", result)
				}
			} else {
				if result == nil {
					t.Error("TopologicalSort() вернул nil для графа без цикла")
					return
				}

				if len(result) != tt.validLength {
					t.Errorf("TopologicalSort() вернул результат длины %d; ожидалось %d",
						len(result), tt.validLength)
				}

				// Проверяем валидность порядка
				if !IsValidTopologicalOrder(tt.numVertices, tt.graph, result) {
					t.Errorf("TopologicalSort() вернул невалидный порядок: %v", result)
				}
			}
		})
	}
}

// TestCanFinish проверяет возможность завершения курсов
func TestCanFinish(t *testing.T) {
	tests := []struct {
		name         string
		numCourses   int
		prerequisites [][]int
		expected     bool
	}{
		{
			name:       "Можно завершить",
			numCourses: 2,
			prerequisites: [][]int{
				{1, 0},
			},
			expected: true,
		},
		{
			name:       "Циклическая зависимость",
			numCourses: 2,
			prerequisites: [][]int{
				{1, 0},
				{0, 1}, // Цикл
			},
			expected: false,
		},
		{
			name:         "Нет зависимостей",
			numCourses:   3,
			prerequisites: [][]int{},
			expected:     true,
		},
		{
			name:       "Сложные зависимости",
			numCourses: 4,
			prerequisites: [][]int{
				{1, 0},
				{2, 0},
				{3, 1},
				{3, 2},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanFinish(tt.numCourses, tt.prerequisites)
			if result != tt.expected {
				t.Errorf("CanFinish(%d, %v) = %v; ожидалось %v",
					tt.numCourses, tt.prerequisites, result, tt.expected)
			}
		})
	}
}

// TestFindOrder проверяет поиск порядка курсов
func TestFindOrder(t *testing.T) {
	tests := []struct {
		name         string
		numCourses   int
		prerequisites [][]int
		expectedLen  int
		hasCycle     bool
	}{
		{
			name:       "Обычный случай",
			numCourses: 4,
			prerequisites: [][]int{
				{1, 0},
				{2, 0},
				{3, 1},
				{3, 2},
			},
			expectedLen: 4,
			hasCycle:    false,
		},
		{
			name:       "Циклическая зависимость",
			numCourses: 2,
			prerequisites: [][]int{
				{1, 0},
				{0, 1},
			},
			expectedLen: 0,
			hasCycle:    true,
		},
		{
			name:         "Нет зависимостей",
			numCourses:   3,
			prerequisites: [][]int{},
			expectedLen:  3,
			hasCycle:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindOrder(tt.numCourses, tt.prerequisites)

			if tt.hasCycle {
				if len(result) != 0 {
					t.Errorf("FindOrder() вернул результат для графа с циклом: %v", result)
				}
			} else {
				if len(result) != tt.expectedLen {
					t.Errorf("FindOrder() вернул результат длины %d; ожидалось %d",
						len(result), tt.expectedLen)
				}
			}
		})
	}
}

// TestFindMinHeightTrees проверяет поиск корней с минимальной высотой
func TestFindMinHeightTrees(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		edges    [][]int
		expected []int
	}{
		{
			name: "Обычный случай",
			n:    4,
			edges: [][]int{
				{1, 0},
				{1, 2},
				{1, 3},
			},
			expected: []int{1},
		},
		{
			name: "Два корня",
			n:    6,
			edges: [][]int{
				{3, 0}, {3, 1}, {3, 2}, {3, 4}, {5, 4},
			},
			expected: []int{3, 4},
		},
		{
			name:     "Одна вершина",
			n:        1,
			edges:    [][]int{},
			expected: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindMinHeightTrees(tt.n, tt.edges)
			// Сортируем для сравнения
			sortInts(result)
			sortInts(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindMinHeightTrees(%d, %v) = %v; ожидалось %v",
					tt.n, tt.edges, result, tt.expected)
			}
		})
	}
}

// sortInts сортирует слайс int
func sortInts(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// TestIsValidTopologicalOrder проверяет валидацию порядка
func TestIsValidTopologicalOrder(t *testing.T) {
	tests := []struct {
		name        string
		numVertices int
		graph       [][]int
		order       []int
		expected    bool
	}{
		{
			name:        "Валидный порядок",
			numVertices: 4,
			graph: [][]int{
				0: {1, 2},
				1: {3},
				2: {3},
				3: {},
			},
			order:    []int{0, 1, 2, 3},
			expected: true,
		},
		{
			name:        "Невалидный порядок",
			numVertices: 3,
			graph: [][]int{
				0: {1},
				1: {2},
				2: {},
			},
			order:    []int{2, 1, 0}, // Неправильный порядок
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTopologicalOrder(tt.numVertices, tt.graph, tt.order)
			if result != tt.expected {
				t.Errorf("IsValidTopologicalOrder(..., %v) = %v; ожидалось %v",
					tt.order, result, tt.expected)
			}
		})
	}
}

