// Package reverselinkedlist демонстрирует паттерн "Разворот связанного списка"
package reverselinkedlist

// ListNode представляет узел связанного списка
type ListNode struct {
	Val  int
	Next *ListNode
}

// ReverseLinkedList итеративно разворачивает связанный список
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	// Список: 1 -> 2 -> 3 -> 4 -> nil
//	head := CreateList([]int{1, 2, 3, 4})
//	reversed := ReverseLinkedList(head) // Вернёт: 4 -> 3 -> 2 -> 1 -> nil
func ReverseLinkedList(head *ListNode) *ListNode {
	var prev *ListNode
	current := head

	for current != nil {
		// Сохраняем следующий узел
		next := current.Next
		// Разворачиваем указатель
		current.Next = prev
		// Двигаемся вперёд
		prev = current
		current = next
	}

	return prev
}

// ReverseLinkedListRecursive рекурсивно разворачивает связанный список
// Временная сложность: O(n)
// Пространственная сложность: O(n) из-за стека вызовов
func ReverseLinkedListRecursive(head *ListNode) *ListNode {
	// Базовый случай: пустой список или один элемент
	if head == nil || head.Next == nil {
		return head
	}

	// Рекурсивно разворачиваем остаток списка
	reversed := ReverseLinkedListRecursive(head.Next)

	// Разворачиваем текущий узел
	head.Next.Next = head
	head.Next = nil

	return reversed
}

// ReverseBetween разворачивает часть списка между позициями left и right
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	// Список: 1 -> 2 -> 3 -> 4 -> 5, left=2, right=4
//	head := CreateList([]int{1, 2, 3, 4, 5})
//	result := ReverseBetween(head, 2, 4) // Вернёт: 1 -> 4 -> 3 -> 2 -> 5
func ReverseBetween(head *ListNode, left int, right int) *ListNode {
	if head == nil || left == right {
		return head
	}

	// Создаём фиктивный узел для упрощения обработки случая, когда left = 1
	dummy := &ListNode{Next: head}
	prev := dummy

	// Двигаемся до позиции left - 1
	for i := 0; i < left-1; i++ {
		prev = prev.Next
	}

	// Начинаем разворот
	current := prev.Next
	for i := 0; i < right-left; i++ {
		next := current.Next
		current.Next = next.Next
		next.Next = prev.Next
		prev.Next = next
	}

	return dummy.Next
}

// ReverseKGroup разворачивает каждые k элементов связанного списка
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	// Список: 1 -> 2 -> 3 -> 4 -> 5, k=2
//	head := CreateList([]int{1, 2, 3, 4, 5})
//	result := ReverseKGroup(head, 2) // Вернёт: 2 -> 1 -> 4 -> 3 -> 5
func ReverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil || k == 1 {
		return head
	}

	// Проверяем, есть ли k элементов
	count := 0
	current := head
	for current != nil && count < k {
		current = current.Next
		count++
	}

	// Если есть k элементов, разворачиваем их
	if count == k {
		// Рекурсивно обрабатываем остаток
		current = ReverseKGroup(current, k)

		// Разворачиваем текущую группу из k элементов
		for count > 0 {
			next := head.Next
			head.Next = current
			current = head
			head = next
			count--
		}
		head = current
	}

	return head
}

// IsPalindrome проверяет, является ли связанный список палиндромом
// Временная сложность: O(n)
// Пространственная сложность: O(1)
//
// Пример:
//
//	// Список: 1 -> 2 -> 2 -> 1
//	head := CreateList([]int{1, 2, 2, 1})
//	result := IsPalindrome(head) // Вернёт true
func IsPalindrome(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true
	}

	// Находим середину списка
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// Разворачиваем вторую половину
	secondHalf := ReverseLinkedList(slow)
	firstHalf := head

	// Сравниваем первую и вторую половины
	for secondHalf != nil {
		if firstHalf.Val != secondHalf.Val {
			return false
		}
		firstHalf = firstHalf.Next
		secondHalf = secondHalf.Next
	}

	return true
}

// CreateList создаёт связанный список из слайса
func CreateList(values []int) *ListNode {
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

// ListToSlice преобразует связанный список в слайс
func ListToSlice(head *ListNode) []int {
	result := []int{}
	current := head

	for current != nil {
		result = append(result, current.Val)
		current = current.Next
	}

	return result
}


