package mergeintervals

import (
	"reflect"
	"testing"
)

// TestMergeIntervals проверяет объединение интервалов
func TestMergeIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals []Interval
		expected  []Interval
	}{
		{
			name:      "Перекрывающиеся интервалы",
			intervals: []Interval{{1, 4}, {2, 5}, {7, 9}},
			expected:  []Interval{{1, 5}, {7, 9}},
		},
		{
			name:      "Все интервалы перекрываются",
			intervals: []Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			expected:  []Interval{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			name:      "Нет перекрытий",
			intervals: []Interval{{1, 2}, {3, 4}, {5, 6}},
			expected:  []Interval{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:      "Один интервал",
			intervals: []Interval{{1, 4}},
			expected:  []Interval{{1, 4}},
		},
		{
			name:      "Пустой список",
			intervals: []Interval{},
			expected:  []Interval{},
		},
		{
			name:      "Все объединяются в один",
			intervals: []Interval{{1, 3}, {2, 4}, {3, 5}},
			expected:  []Interval{{1, 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeIntervals(tt.intervals)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeIntervals(%v) = %v; ожидалось %v",
					tt.intervals, result, tt.expected)
			}
		})
	}
}

// TestInsertInterval проверяет вставку интервала
func TestInsertInterval(t *testing.T) {
	tests := []struct {
		name        string
		intervals   []Interval
		newInterval Interval
		expected    []Interval
	}{
		{
			name:        "Вставка с объединением",
			intervals:   []Interval{{1, 3}, {5, 7}, {8, 12}},
			newInterval: Interval{4, 6},
			expected:    []Interval{{1, 3}, {4, 7}, {8, 12}},
		},
		{
			name:        "Вставка в начало",
			intervals:   []Interval{{5, 7}, {8, 12}},
			newInterval: Interval{1, 3},
			expected:    []Interval{{1, 3}, {5, 7}, {8, 12}},
		},
		{
			name:        "Вставка в конец",
			intervals:   []Interval{{1, 3}, {5, 7}},
			newInterval: Interval{8, 12},
			expected:    []Interval{{1, 3}, {5, 7}, {8, 12}},
		},
		{
			name:        "Вставка объединяет все",
			intervals:   []Interval{{1, 3}, {5, 7}},
			newInterval: Interval{2, 6},
			expected:    []Interval{{1, 7}},
		},
		{
			name:        "Вставка в пустой список",
			intervals:   []Interval{},
			newInterval: Interval{1, 3},
			expected:    []Interval{{1, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertInterval(tt.intervals, tt.newInterval)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertInterval(%v, %v) = %v; ожидалось %v",
					tt.intervals, tt.newInterval, result, tt.expected)
			}
		})
	}
}

// TestIntervalIntersection проверяет пересечение интервалов
func TestIntervalIntersection(t *testing.T) {
	tests := []struct {
		name      string
		intervalsA []Interval
		intervalsB []Interval
		expected   []Interval
	}{
		{
			name:       "Обычное пересечение",
			intervalsA: []Interval{{1, 3}, {5, 6}, {7, 9}},
			intervalsB: []Interval{{2, 3}, {5, 7}},
			expected:   []Interval{{2, 3}, {5, 6}, {7, 7}}, // {7, 7} - валидное пересечение в точке
		},
		{
			name:       "Нет пересечений",
			intervalsA: []Interval{{1, 2}, {5, 6}},
			intervalsB: []Interval{{3, 4}, {7, 8}},
			expected:   []Interval{},
		},
		{
			name:       "Полное перекрытие",
			intervalsA: []Interval{{1, 5}},
			intervalsB: []Interval{{2, 3}},
			expected:   []Interval{{2, 3}},
		},
		{
			name:       "Один список пуст",
			intervalsA: []Interval{{1, 2}},
			intervalsB: []Interval{},
			expected:   []Interval{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntervalIntersection(tt.intervalsA, tt.intervalsB)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IntervalIntersection(%v, %v) = %v; ожидалось %v",
					tt.intervalsA, tt.intervalsB, result, tt.expected)
			}
		})
	}
}

// TestCanAttendAllMeetings проверяет возможность посетить все встречи
func TestCanAttendAllMeetings(t *testing.T) {
	tests := []struct {
		name      string
		intervals []Interval
		expected  bool
	}{
		{
			name:      "Есть конфликт",
			intervals: []Interval{{1, 4}, {2, 5}, {7, 9}},
			expected:  false,
		},
		{
			name:      "Нет конфликтов",
			intervals: []Interval{{1, 2}, {3, 4}, {5, 6}},
			expected:  true,
		},
		{
			name:      "Одна встреча",
			intervals: []Interval{{1, 4}},
			expected:  true,
		},
		{
			name:      "Пустой список",
			intervals: []Interval{},
			expected:  true,
		},
		{
			name:      "Граничные случаи (конец = начало)",
			intervals: []Interval{{1, 2}, {2, 3}, {3, 4}},
			expected:  true, // Можно посетить, если встречи не перекрываются
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanAttendAllMeetings(tt.intervals)
			if result != tt.expected {
				t.Errorf("CanAttendAllMeetings(%v) = %v; ожидалось %v",
					tt.intervals, result, tt.expected)
			}
		})
	}
}

// TestMinMeetingRooms проверяет поиск минимального количества комнат
func TestMinMeetingRooms(t *testing.T) {
	tests := []struct {
		name      string
		intervals []Interval
		expected  int
	}{
		{
			name:      "Нужны 2 комнаты",
			intervals: []Interval{{1, 4}, {2, 5}, {7, 9}},
			expected:  2,
		},
		{
			name:      "Нужна 1 комната",
			intervals: []Interval{{1, 2}, {3, 4}, {5, 6}},
			expected:  1,
		},
		{
			name:      "Нужны 3 комнаты",
			intervals: []Interval{{1, 3}, {2, 4}, {3, 5}},
			expected:  2, // В момент времени 3 одновременно идут 2 встречи
		},
		{
			name:      "Все перекрываются",
			intervals: []Interval{{1, 3}, {2, 4}, {2, 5}},
			expected:  3,
		},
		{
			name:      "Пустой список",
			intervals: []Interval{},
			expected:  0,
		},
		{
			name:      "Одна встреча",
			intervals: []Interval{{1, 3}},
			expected:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MinMeetingRooms(tt.intervals)
			if result != tt.expected {
				t.Errorf("MinMeetingRooms(%v) = %d; ожидалось %d",
					tt.intervals, result, tt.expected)
			}
		})
	}
}

