// Package fastslowpointers демонстрирует паттерн "Быстрые и медленные указатели"
package fastslowpointers

// ListNode представляет узел связанного списка
type ListNode struct {
	Val  int       // Значение узла
	Next *ListNode // Указатель на следующий узел
}

// HasCycle определяет, есть ли цикл в связанном списке
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Алгоритм: если быстрый указатель догонит медленный, значит есть цикл
//
// Пример:
//
//	// Создаём список с циклом: 1 -> 2 -> 3 -> 4 -> 2 (цикл)
//	head := &ListNode{Val: 1}
//	head.Next = &ListNode{Val: 2}
//	head.Next.Next = &ListNode{Val: 3}
//	head.Next.Next.Next = &ListNode{Val: 4}
//	head.Next.Next.Next.Next = head.Next // Создаём цикл
//	result := HasCycle(head) // Вернёт true
func HasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow := head // Медленный указатель (черепаха)
	fast := head // Быстрый указатель (заяц)

	for fast != nil && fast.Next != nil {
		slow = slow.Next      // Медленный делает 1 шаг
		fast = fast.Next.Next // Быстрый делает 2 шага

		// Если встретились — есть цикл
		if slow == fast {
			return true
		}
	}

	// Быстрый достиг конца — цикла нет
	return false
}

// FindCycleStart находит начало цикла в связанном списке
// Возвращает nil, если цикла нет
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Алгоритм:
// 1. Находим точку встречи быстрого и медленного указателей
// 2. Сбрасываем медленный указатель на начало
// 3. Двигаем оба указателя по одному шагу — они встретятся в начале цикла
func FindCycleStart(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}

	slow := head
	fast := head

	// Шаг 1: Находим точку встречи
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			// Шаг 2: Сбрасываем slow на начало
			slow = head

			// Шаг 3: Двигаем оба по одному шагу до встречи
			for slow != fast {
				slow = slow.Next
				fast = fast.Next
			}

			return slow // Начало цикла
		}
	}

	return nil // Цикла нет
}

// FindMiddle находит средний узел связанного списка
// Если узлов чётное количество, возвращает второй из двух средних
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	// Список: 1 -> 2 -> 3 -> 4 -> 5
//	result := FindMiddle(head) // Вернёт узел со значением 3
//
//	// Список: 1 -> 2 -> 3 -> 4 -> 5 -> 6
//	result := FindMiddle(head) // Вернёт узел со значением 4
func FindMiddle(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	slow := head
	fast := head

	// Когда быстрый дойдёт до конца, медленный будет в середине
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	return slow
}

// IsHappyNumber определяет, является ли число "счастливым"
// Счастливое число: если повторно заменять число суммой квадратов его цифр,
// в конечном итоге получится 1
// Временная сложность: O(log n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	IsHappyNumber(23) // true: 2²+3²=13 → 1²+3²=10 → 1²+0²=1
//	IsHappyNumber(12) // false: войдёт в цикл
func IsHappyNumber(n int) bool {
	slow := n
	fast := n

	for {
		slow = sumOfSquares(slow)           // 1 шаг
		fast = sumOfSquares(sumOfSquares(fast)) // 2 шага

		if fast == 1 {
			return true // Достигли 1 — число счастливое
		}

		if slow == fast {
			return false // Цикл без 1 — число не счастливое
		}
	}
}

// sumOfSquares вычисляет сумму квадратов цифр числа
func sumOfSquares(n int) int {
	sum := 0
	for n > 0 {
		digit := n % 10
		sum += digit * digit
		n /= 10
	}
	return sum
}

// GetCycleLength возвращает длину цикла в связанном списке
// Возвращает 0, если цикла нет
// Временная сложность: O(n)
// Пространственная сложность: O(1)
func GetCycleLength(head *ListNode) int {
	if head == nil || head.Next == nil {
		return 0
	}

	slow := head
	fast := head

	// Находим точку встречи
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			// Считаем длину цикла
			return calculateCycleLength(slow)
		}
	}

	return 0
}

// calculateCycleLength вычисляет длину цикла, начиная с узла
func calculateCycleLength(node *ListNode) int {
	current := node
	length := 0

	// Обходим цикл, пока не вернёмся к начальному узлу
	for {
		current = current.Next
		length++
		if current == node {
			break
		}
	}

	return length
}

// CreateLinkedList создаёт связанный список из слайса
// Вспомогательная функция для тестирования
func CreateLinkedList(values []int) *ListNode {
	if len(values) == 0 {
		return nil
	}

	head := &ListNode{Val: values[0]}
	current := head

	for i := 1; i < len(values); i++ {
		current.Next = &ListNode{Val: values[i]}
		current = current.Next
	}

	return head
}

// LinkedListToSlice преобразует связанный список в слайс
// Максимум limit элементов (для предотвращения бесконечных циклов)
func LinkedListToSlice(head *ListNode, limit int) []int {
	result := []int{}
	current := head

	for current != nil && len(result) < limit {
		result = append(result, current.Val)
		current = current.Next
	}

	return result
}


