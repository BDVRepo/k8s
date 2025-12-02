// Package mergeintervals демонстрирует паттерн "Слияние интервалов"
package mergeintervals

import "sort"

// Interval представляет интервал [Start, End]
type Interval struct {
	Start int // Начало интервала
	End   int // Конец интервала
}

// MergeIntervals объединяет перекрывающиеся интервалы
// Временная сложность: O(n log n) из-за сортировки
// Пространственная сложность: O(n) для результата
//
// Пример:
//
//	intervals := []Interval{{1, 4}, {2, 5}, {7, 9}}
//	result := MergeIntervals(intervals) // Вернёт [{1, 5}, {7, 9}]
func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) < 2 {
		return intervals
	}

	// Сортируем интервалы по началу
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	merged := []Interval{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		lastMerged := &merged[len(merged)-1]
		current := intervals[i]

		// Если текущий интервал перекрывается с последним объединённым
		// (начало текущего <= конец последнего)
		if current.Start <= lastMerged.End {
			// Объединяем: берём максимум концов
			if current.End > lastMerged.End {
				lastMerged.End = current.End
			}
		} else {
			// Не перекрывается — добавляем как новый интервал
			merged = append(merged, current)
		}
	}

	return merged
}

// InsertInterval вставляет новый интервал в отсортированный список и объединяет перекрытия
// Временная сложность: O(n)
// Пространственная сложность: O(n)
//
// Пример:
//
//	intervals := []Interval{{1, 3}, {5, 7}, {8, 12}}
//	newInterval := Interval{4, 6}
//	result := InsertInterval(intervals, newInterval) // Вернёт [{1, 3}, {4, 7}, {8, 12}]
func InsertInterval(intervals []Interval, newInterval Interval) []Interval {
	result := []Interval{}
	i := 0
	n := len(intervals)

	// Добавляем все интервалы, которые заканчиваются до начала нового
	for i < n && intervals[i].End < newInterval.Start {
		result = append(result, intervals[i])
		i++
	}

	// Объединяем все перекрывающиеся интервалы
	merged := newInterval
	for i < n && intervals[i].Start <= merged.End {
		if intervals[i].Start < merged.Start {
			merged.Start = intervals[i].Start
		}
		if intervals[i].End > merged.End {
			merged.End = intervals[i].End
		}
		i++
	}

	result = append(result, merged)

	// Добавляем оставшиеся интервалы
	for i < n {
		result = append(result, intervals[i])
		i++
	}

	return result
}

// IntervalIntersection находит пересечение двух списков интервалов
// Временная сложность: O(n + m), где n и m — длины списков
// Пространственная сложность: O(k), где k — количество пересечений
//
// Пример:
//
//	intervalsA := []Interval{{1, 3}, {5, 6}, {7, 9}}
//	intervalsB := []Interval{{2, 3}, {5, 7}}
//	result := IntervalIntersection(intervalsA, intervalsB) // Вернёт [{2, 3}, {5, 6}]
func IntervalIntersection(intervalsA, intervalsB []Interval) []Interval {
	result := []Interval{}
	i, j := 0, 0

	for i < len(intervalsA) && j < len(intervalsB) {
		a := intervalsA[i]
		b := intervalsB[j]

		// Проверяем, есть ли пересечение
		// Пересечение существует, если max(a.Start, b.Start) <= min(a.End, b.End)
		start := max(a.Start, b.Start)
		end := min(a.End, b.End)

		if start <= end {
			result = append(result, Interval{Start: start, End: end})
		}

		// Переходим к следующему интервалу в том списке, который заканчивается раньше
		if a.End < b.End {
			i++
		} else {
			j++
		}
	}

	return result
}

// CanAttendAllMeetings проверяет, можно ли посетить все встречи без конфликтов
// Временная сложность: O(n log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	intervals := []Interval{{1, 4}, {2, 5}, {7, 9}}
//	result := CanAttendAllMeetings(intervals) // Вернёт false (конфликт между первыми двумя)
func CanAttendAllMeetings(intervals []Interval) bool {
	if len(intervals) < 2 {
		return true
	}

	// Сортируем по началу
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	// Проверяем, нет ли перекрытий
	for i := 1; i < len(intervals); i++ {
		if intervals[i].Start < intervals[i-1].End {
			return false // Есть конфликт
		}
	}

	return true
}

// MinMeetingRooms находит минимальное количество переговорных комнат
// Временная сложность: O(n log n)
// Пространственная сложность: O(n)
//
// Пример:
//
//	intervals := []Interval{{1, 4}, {2, 5}, {7, 9}}
//	result := MinMeetingRooms(intervals) // Вернёт 2
func MinMeetingRooms(intervals []Interval) int {
	if len(intervals) == 0 {
		return 0
	}

	// Создаём массивы начал и концов
	starts := make([]int, len(intervals))
	ends := make([]int, len(intervals))

	for i, interval := range intervals {
		starts[i] = interval.Start
		ends[i] = interval.End
	}

	// Сортируем оба массива
	sort.Ints(starts)
	sort.Ints(ends)

	// Используем два указателя для отслеживания перекрытий
	rooms := 0
	endPointer := 0

	for i := 0; i < len(starts); i++ {
		// Если текущее начало >= конец встречи, которую отслеживает endPointer,
		// значит эта комната освободилась
		if starts[i] >= ends[endPointer] {
			endPointer++
		} else {
			// Нужна новая комната
			rooms++
		}
	}

	return rooms
}

// Вспомогательные функции min и max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}


